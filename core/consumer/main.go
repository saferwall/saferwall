// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	minio "github.com/minio/minio-go"

	avast "github.com/saferwall/saferwall/core/multiav/avast/client"
	avira "github.com/saferwall/saferwall/core/multiav/avira/client"
	bitdefender "github.com/saferwall/saferwall/core/multiav/bitdefender/client"
	clamav "github.com/saferwall/saferwall/core/multiav/clamav/client"
	comodo "github.com/saferwall/saferwall/core/multiav/comodo/client"
	eset "github.com/saferwall/saferwall/core/multiav/eset/client"
	fsecure "github.com/saferwall/saferwall/core/multiav/fsecure/client"
	kaspersky "github.com/saferwall/saferwall/core/multiav/kaspersky/client"
	symantec "github.com/saferwall/saferwall/core/multiav/symantec/client"
	windefender "github.com/saferwall/saferwall/core/multiav/windefender/client"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/trid"
	"github.com/saferwall/saferwall/pkg/utils"
	"github.com/saferwall/saferwall/pkg/utils/do"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	client *minio.Client

	backendEndpoint string
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}
type result struct {
	Crc32   string                 `json:"crc32,omitempty"`
	Md5     string                 `json:"md5,omitempty"`
	Sha1    string                 `json:"sha1,omitempty"`
	Sha256  string                 `json:"sha256,omitempty"`
	Sha512  string                 `json:"sha512,omitempty"`
	Ssdeep  string                 `json:"ssdeep,omitempty"`
	Exif    map[string]string      `json:"exif,omitempty"`
	TriD    []string               `json:"trid,omitempty"`
	Packer  []string               `json:"packer,omitempty"`
	Magic   string                 `json:"magic,omitempty"`
	Strings []stringStruct         `json:"strings,omitempty"`
	MultiAV map[string]interface{} `json:"multiav,omitempty"`
	Status  int                    `json:"status,omitempty"`
}

const (
	queued     = iota
	processing = iota
	finished   = iota
)

// NoopNSQLogger allows us to pipe NSQ logs to dev/null
// The default NSQ logger is great for debugging, but did
// not fit our normally well structured JSON logs. Luckily
// NSQ provides a simple interface for injecting your own
// logger.
type NoopNSQLogger struct{}

// Output allows us to implement the nsq.Logger interface
func (l *NoopNSQLogger) Output(calldepth int, s string) error {
	log.Info(s)
	return nil
}

// MessageHandler adheres to the nsq.Handler interface.
// This allows us to define our own custome handlers for
// our messages. Think of these handlers much like you would
// an http handler.
type MessageHandler struct{}

// HandleMessage is the only requirement needed to fulfill the
// nsq.Handler interface. This where you'll write your message
// handling logic.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// returning an error results in the message being re-enqueued
		// a REQ is sent to nsqd
		return errors.New("body is blank re-enqueue message")
	}

	sha256 := string(m.Body)
	log.Infof("Processing %s", sha256)

	// set the file status to `processing`
	res := result{}
	res.Status = processing

	// Marshell results
	buff, err := json.Marshal(res)
	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}

	// Update document
	updateDocument(sha256, buff)

	// 	Where to save the sample, in k8s, our nfs share
	filePath := path.Join("/samples", sha256)

	// Download the sample
	bucketName := viper.GetString("do.spacename")
	err = client.FGetObject(bucketName, sha256, filePath, minio.GetObjectOptions{})

	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}
	log.Infof("File fetched success %s", sha256)

	b, err := utils.ReadAll(filePath)
	if err != nil {
		log.Error("Failed to read file, err: ", err)
		return err
	}

	// Run crypto pkg
	r := crypto.HashBytes(b)
	res.Crc32 = r.Crc32
	res.Md5 = r.Md5
	res.Sha1 = r.Sha1
	res.Sha256 = r.Sha256
	res.Sha512 = r.Sha512
	res.Ssdeep = r.Ssdeep
	log.Infof("HashBytes success %s", sha256)

	// Run exiftool pkg
	res.Exif, err = exiftool.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with exiftool, err: ", err)
	}
	log.Infof("exiftool success %s", sha256)

	// Run TRiD pkg
	res.TriD, err = trid.Scan(filePath)
	if err != nil {
		log.Error("Faileds to scan file with trid, err: ", err)
	}
	log.Infof("trid success %s", sha256)

	// Run Magic Pkg
	res.Magic, err = magic.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with magic, err: ", err)
	}
	log.Infof("magic extraction success %s", sha256)

	// Run Die Pkg
	res.Packer, err = packer.Scan(filePath)
	if err != nil {
		log.Error("Failed to scan file with packer, err: ", err)
	}
	log.Infof("packer extraction success %s", sha256)

	// Run strings pkg
	n := 10
	asciiStrings := s.GetASCIIStrings(b, n)
	wideStrings := s.GetUnicodeStrings(b, n)
	asmStrings := s.GetAsmStrings(b)

	// Remove duplicates
	uniqueASCII := utils.UniqueSlice(asciiStrings)
	uniqueWide := utils.UniqueSlice(wideStrings)
	uniqueAsm := utils.UniqueSlice(asmStrings)

	var strResults []stringStruct
	for _, str := range uniqueASCII {
		strResults = append(strResults, stringStruct{"ascii", str})
	}

	for _, str := range uniqueWide {
		strResults = append(strResults, stringStruct{"wide", str})
	}

	for _, str := range uniqueAsm {
		strResults = append(strResults, stringStruct{"asm", str})
	}
	res.Strings = strResults
	log.Infof("strings success %s", sha256)

	multiavScanResults := map[string]interface{}{}

	// Scan with Kaspersky
	kasperskyClient, err := kaspersky.Init()
	if err != nil {
		log.Errorf("kaspersky init failed %s", err)
	} else {
		kasperskyRes, err := kaspersky.ScanFile(kasperskyClient, filePath)
		if err != nil {
			log.Errorf("kaspersky scanfile failed %s", err)
		} else {
			multiavScanResults["kaspersky"] = kasperskyRes
			log.Infof("kaspersky success %s", sha256)
		}
	}

	// Scan with ClamAV
	clamclient, err := clamav.Init()
	if err != nil {
		log.Errorf("clamav init failed %s", err)
	} else {
		clamres, err := clamav.ScanFile(clamclient, filePath)
		if err != nil {
			log.Errorf("clamav scanfile failed %s", err)
		} else {
			multiavScanResults["clamav"] = clamres
			log.Infof("clamav success %s", sha256)
		}
	}

	// Scan with Avast
	avastClient, err := avast.Init()
	if err != nil {
		log.Errorf("avast init failed %s", err)
	} else {
		avastres, err := avast.ScanFile(avastClient, filePath)
		if err != nil {
			log.Errorf("avast scanfile failed %s", err)
		} else {
			multiavScanResults["avast"] = avastres
			log.Infof("avast success %s", sha256)
		}
	}

	// Scan with Avira
	aviraClient, err := avira.Init()
	if err != nil {
		log.Errorf("avira init failed %s", err)
	} else {
		avirares, err := avira.ScanFile(aviraClient, filePath)
		if err != nil {
			log.Errorf("avira scanfile failed %s", err)
		} else {
			multiavScanResults["avira"] = avirares
			log.Infof("avira success %s", sha256)
		}
	}

	// Scan with Bitdefender
	bitdefenderClient, err := bitdefender.Init()
	if err != nil {
		log.Errorf("bitdefender init failed %s", err)
	} else {
		bitdefenderres, err := bitdefender.ScanFile(bitdefenderClient, filePath)
		if err != nil {
			log.Errorf("bitdefender scanfile failed %s", err)
		} else {
			multiavScanResults["bitdefender"] = bitdefenderres
			log.Infof("bitdefender success %s", sha256)
		}
	}

	// Scan with Comodo
	comodoClient, err := comodo.Init()
	if err != nil {
		log.Errorf("comodo init failed %s", err)
	} else {
		comodores, err := comodo.ScanFile(comodoClient, filePath)
		if err != nil {
			log.Errorf("comodo scanfile failed %s", err)
		} else {
			multiavScanResults["comodo"] = comodores
			log.Infof("comodo success %s", sha256)
		}
	}

	// Scan with Windows Defender
	windefenderClient, err := windefender.Init()
	if err != nil {
		log.Errorf("windefender init failed %s", err)
	} else {
		windefenderRes, err := windefender.ScanFile(windefenderClient, filePath)
		if err != nil {
			log.Errorf("windefender scanfile failed %s", err)
		} else {
			multiavScanResults["windefender"] = windefenderRes
			log.Infof("windefender success %s", sha256)
		}
	}

	// Scan with FSecure
	fsecureClient, err := fsecure.Init()
	if err != nil {
		log.Errorf("fsecure init failed %s", err)
	} else {
		fsecureRes, err := fsecure.ScanFile(fsecureClient, filePath)
		if err != nil {
			log.Errorf("fsecure scanfile failed %s", err)
		} else {
			multiavScanResults["fsecure"] = fsecureRes
			log.Infof("fsecure success %s", sha256)
		}
	}

	// Scan with Eset
	esetClient, err := eset.Init()
	if err != nil {
		log.Errorf("eset init failed %s", err)
	} else {
		esetRes, err := eset.ScanFile(esetClient, filePath)
		if err != nil {
			log.Errorf("eset scanfile failed %s", err)
		} else {
			multiavScanResults["eset"] = esetRes
			log.Infof("eset success %s", sha256)
		}
	}

	// Scan with Symantec
	symantecClient, err := symantec.Init()
	if err != nil {
		log.Errorf("symantec init failed %s", err)
	} else {
		symantecRes, err := symantec.ScanFile(symantecClient, filePath)
		if err != nil {
			log.Errorf("symantec scanfile failed %s", err)
		} else {
			multiavScanResults["symantec"] = symantecRes
			log.Infof("symantec success %s", sha256)
		}
	}

	res.MultiAV = multiavScanResults

	// analysis finished
	res.Status = finished

	// Marshell results
	buff, err = json.Marshal(res)
	if err != nil {
		log.Error("Failed to get object, err: ", err)
		return err
	}

	// Update document
	updateDocument(sha256, buff)

	// Returning nil signals to the consumer that the message has
	// been handled with success. A FIN is sent to nsqd
	return nil
}

// loadConfig loads our configration.
func loadConfig() {
	viper.SetConfigName("saferwall") // no need to include file extension
	viper.AddConfigPath(".")         // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func updateDocument(sha256 string, buff []byte) {
	// Update results to DB
	client := &http.Client{}
	client.Timeout = time.Second * 15
	url := backendEndpoint + sha256
	log.Infoln("Sending results to ", url)

	body := bytes.NewBuffer(buff)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Errorf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}

	log.Infof("Response status code: %d, text: %s", resp.StatusCode, string(d))
}

func main() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	client = do.GetClient()

	// Load consumer config
	loadConfig()

	// Set backend API address
	backendEndpoint = viper.GetString("backend.address") + "/v1/files/"

	// The default config settings provide a pretty good starting point for
	// our new consumer.
	config := nsq.NewConfig()

	// Create a NewConsumer with the name of our topic, the channel, and our config
	consumer, err := nsq.NewConsumer("scan", "file", config)
	if err != nil {
		log.Errorln("Could not create consumer")
	}

	// Set the number of messages that can be in flight at any given time
	// you'll want to set this number as the default is only 1. This can be
	// a major concurrency knob that can change the performance of your application.
	consumer.ChangeMaxInFlight(200)

	// Here we set the logger to our NoopNSQLogger to quiet down the default logs.
	// At Reverb we use a custom structured logging format so we'll take the logging
	// from here.
	consumer.SetLogger(
		&NoopNSQLogger{},
		nsq.LogLevelError,
	)

	// Injects our handler into the consumer. You'll define one handler
	// per consumer, but you can have as many concurrently running handlers
	// as specified by the second argument. If your MaxInFlight is less
	// than your number of concurrent handlers you'll  starve your workers
	// as there will never be enough in flight messages for your worker pool
	consumer.AddConcurrentHandlers(
		&MessageHandler{},
		1,
	)

	// Our consumer will discover where topics are located by our three
	// nsqlookupd instances The application will periodically poll
	// these nqslookupd instances to discover new nodes or drop unhealthy
	// producers.
	nsqlds := viper.GetStringSlice("nsq.lookupd")
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		log.Fatal(err)
	}

	log.Infoln("Connected to nsqlookupd")

	// Let's allow our queues to drain properly during shutdown.
	// We'll create a channel to listen for SIGINT (Ctrl+C) to signal
	// to our application to gracefully shutdown.
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	// This is our main loop. It will continue to read off of our nsq
	// channel until either the consumer dies or our application is signaled
	// to stop.
	for {
		select {
		case <-consumer.StopChan:
			return // uh oh consumer disconnected. Time to quit.
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			consumer.Stop()
		}
	}
}
