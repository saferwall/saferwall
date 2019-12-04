// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"bytes"
	"errors"
	"fmt"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	"github.com/saferwall/saferwall/pkg/utils"
	avastclient "github.com/saferwall/saferwall/pkg/grpc/multiav/avast/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/avast/proto"
	aviraclient "github.com/saferwall/saferwall/pkg/grpc/multiav/avira/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/avira/proto"
	bitdefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/proto"
	clamavclient "github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/proto"
	comodoclient "github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/proto"
	esetclient "github.com/saferwall/saferwall/pkg/grpc/multiav/eset/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/eset/proto"
	fsecureclient "github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/proto"
	kasperskyclient "github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/proto"
	mcafeeclient "github.com/saferwall/saferwall/pkg/grpc/multiav/mcafee/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/mcafee/proto"
	symantecclient "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	sophosclient "github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/proto"
	windefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"time"
)

// LoadConfig loads our configration.
func LoadConfig() {
	viper.AddConfigPath("../../configs") // set the path of your config file

	// Load the config type depending on env variable.
	var name string
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "dev":
		name = "saferwall.dev"
	case "prod":
		name = "saferwall.prod"
	case "test":
		name = "saferwall.test"
	default:
		log.Fatal("ENVIRONMENT is not set")
	}

	viper.SetConfigName(name)    // no need to include file extension
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Config %s was loaded", name)
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

// Pretty print error.
func check(engine string, err error) {
	if err != nil {
		log.Errorf("[%s]: %v", engine, err)
	}
}

// Return a gRPC client connextion for a given engine.
func avgRPCConn(engine string) (*grpc.ClientConn, error) {

	// Get the address of AV gRPC server
	key := "multiav." + engine + "_addr"
	address := viper.GetString(key)

	// Dial creates a client connection to the given target.
	conn, err := grpc.Dial(
		address, []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		msg := fmt.Sprintf("[%s]: %v", engine, err)
		return nil, errors.New(msg)
	}
	return conn, nil
}

func avScan(engine string, filePath string, c chan multiav.ScanResult) {

	// Get a gRPC client scanner instance for a given engine.
	conn, err := avgRPCConn(engine)
	if err != nil {
		c <- multiav.ScanResult{}
		return
	}
	defer conn.Close()

	// Make a copy of the file for each AV engine.
	// This tries to fix the file locking issues which happens
	// if you try to scan a filepath in a nfs share with
	// different engines at the same time.
	filecopyPath := filePath+"-"+engine
	err = utils.CopyFile(filePath, filecopyPath)
	if err != nil {
		log.Errorf("Failed to copy the file for engine %s.", engine)
		c <- multiav.ScanResult{}
		return
	}

	filePath = filecopyPath

	switch engine {
	case "avast":
		res, err := avastclient.ScanFile(avast_api.NewAvastScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "avira":
		res, err := aviraclient.ScanFile(avira_api.NewAviraScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "bitdefender":
		res, err := bitdefenderclient.ScanFile(bitdefender_api.NewBitdefenderScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "clamav":
		res, err := clamavclient.ScanFile(clamav_api.NewClamAVScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "comodo":
		res, err := comodoclient.ScanFile(comodo_api.NewComodoScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "eset":
		res, err := esetclient.ScanFile(eset_api.NewEsetScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "fsecure":
		res, err := fsecureclient.ScanFile(fsecure_api.NewFSecureScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "kaspersky":
		res, err := kasperskyclient.ScanFile(kaspersky_api.NewKasperskyScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "mcafee":
		res, err := mcafeeclient.ScanFile(mcafee_api.NewMcAfeeScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "symantec":
		res, err := symantecclient.ScanFile(symantec_api.NewSymantecScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "sophos":
		res, err := sophosclient.ScanFile(sophos_api.NewSophosScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	case "windefender":
		res, err := windefenderclient.ScanFile(windefender_api.NewWinDefenderScannerClient(conn), filePath)
		check(engine, err)
		c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}
	}
}

func multiAvScan(filePath string) map[string]interface{} {

	// Create channels to receive scan results.
	aviraChan := make(chan multiav.ScanResult)
	avastChan := make(chan multiav.ScanResult)
	bitdefenderChan := make(chan multiav.ScanResult)
	clamavChan := make(chan multiav.ScanResult)
	comodoChan := make(chan multiav.ScanResult)
	esetChan := make(chan multiav.ScanResult)
	fsecureChan := make(chan multiav.ScanResult)
	kasperskyChan := make(chan multiav.ScanResult)
	mcafeeChan := make(chan multiav.ScanResult)
	symantecChan := make(chan multiav.ScanResult)
	sophosChan := make(chan multiav.ScanResult)
	windefenderChan := make(chan multiav.ScanResult)

	// We Start as much go routines as the AV engines we have.
	// Each go-routines makes a gRPC calls and waits for results.
	// Avast
	go avScan("eset", filePath, esetChan)
	go avScan("fsecure", filePath, fsecureChan)
	go avScan("avira", filePath, aviraChan)
	go avScan("bitdefender", filePath, bitdefenderChan)
	go avScan("kaspersky", filePath, kasperskyChan)
	go avScan("symantec", filePath, symantecChan)
	go avScan("sophos", filePath, sophosChan)
	go avScan("windefender", filePath, windefenderChan)
	go avScan("clamav", filePath, clamavChan)
	go avScan("comodo", filePath, comodoChan)
	go avScan("avast", filePath, avastChan)
	go avScan("mcafee", filePath, mcafeeChan)

	multiavScanResults := map[string]interface{}{}
	avEnginesCount := 11
	avCount := 0
	for {
		select {
		case aviraRes := <-aviraChan:
			multiavScanResults["avira"] = aviraRes
			avCount++
		case avastRes := <-avastChan:
			multiavScanResults["avast"] = avastRes
			avCount++
		case bitdefenderRes := <-bitdefenderChan:
			multiavScanResults["bitdefender"] = bitdefenderRes
			avCount++
		case clamavRes := <-clamavChan:
			multiavScanResults["clamav"] = clamavRes
			avCount++
		case comodoRes := <-comodoChan:
			multiavScanResults["comodo"] = comodoRes
			avCount++
		case esetRes := <-esetChan:
			multiavScanResults["eset"] = esetRes
			avCount++
		case fsecureRes := <-fsecureChan:
			multiavScanResults["fsecure"] = fsecureRes
			avCount++
		case kasperskyRes := <-kasperskyChan:
			multiavScanResults["kaspersky"] = kasperskyRes
			avCount++
		case mcafeeRes := <-mcafeeChan:
			multiavScanResults["mcafee"] = mcafeeRes
			avCount++
		case symantecRes := <-symantecChan:
			multiavScanResults["symanetc"] = symantecRes
			avCount++
		case sophosRes := <-sophosChan:
			multiavScanResults["sophos"] = sophosRes
			avCount++
		case windefenderRes := <-windefenderChan:
			multiavScanResults["windefender"] = windefenderRes
			avCount++
		}

		if avCount == avEnginesCount {
			break
		}
	}

	log.Infoln("multiav scan finished")
	return multiavScanResults
}
