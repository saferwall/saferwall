// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"

	"context"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/avast/proto"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/avira/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/dummy/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/eset/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	// "github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/proto"
	"google.golang.org/grpc"
)

const (
	avEnginesCount = 1
)

// loadConfig loads our configration.
func loadConfig(cfgPath string) error {
	viper.SetConfigName("saferwall") // no need to include file extension
	viper.AddConfigPath(cfgPath)     // set the path of your config file
	err := viper.ReadInConfig()
	return err
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

func multiAvScan(filePath string) map[string]interface{} {
	// We Start as much georoutines as the AV engines we have.
	// Each goroutines makes a gRPC calls and waits for results.
	multiavScanResults := map[string]interface{}{}
	// avastChan := make(chan multiav.ScanResult)
	aviraChan := make(chan multiav.ScanResult)
	// bitdefenderChan := make(chan multiav.ScanResult)
	// clamavChan := make(chan multiav.ScanResult)
	// comodoChan := make(chan multiav.ScanResult)
	// esetChan := make(chan multiav.ScanResult)
	// fsecureChan := make(chan multiav.ScanResult)
	// kasperskyChan := make(chan multiav.ScanResult)
	// windefenderChan := make(chan multiav.ScanResult)
	// symantecChan := make(chan multiav.ScanResult)
	// dummyChan := make(chan multiav.ScanResult)

	// go func() {
	// 	address := viper.GetString("multiav.dummy_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 		dummyChan <- multiav.ScanResult{}
	// 		return
	// 	}
	// 	defer conn.Close()

	// 	client := dummy_api.NewDummyScannerClient(conn)
	// 	scanFileRequest := &dummy_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with dummy: :v", err)
	// 		dummyChan <- multiav.ScanResult{}
	// 	} else {
	// 		dummyChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()
	// go func() {
	// 	address := viper.GetString("multiav.avast_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 		avastChan <- multiav.ScanResult{}
	// 		return

	// 	}
	// 	defer conn.Close()

	// 	client := avast_api.NewAvastScannerClient(conn)
	// 	scanFileRequest := &avast_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFilePath(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with avast: :v", err)
	// 		avastChan <- multiav.ScanResult{}
	// 	} else {
	// 		avastChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	go func() {
		address := viper.GetString("multiav.avira_addr")
		conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
		if err != nil {
			log.Errorf("fail to dial: %v", err)
			aviraChan <- multiav.ScanResult{}
			return

		}
		defer conn.Close()

		client := avira_api.NewAviraScannerClient(conn)
		scanFileRequest := &avira_api.ScanFileRequest{Filepath: filePath}
		res, err := client.ScanFile(context.Background(), scanFileRequest)
		if err != nil {
			log.Errorln("Failed to scan with avira: :v", err)
			aviraChan <- multiav.ScanResult{}
		} else {
			aviraChan <- multiav.ScanResult{
				Output:   res.Output,
				Infected: res.Infected,
				Update:   res.Update}
		}
	}()

	// go func() {
	// 	address := viper.GetString("multiav.bitdefender_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := bitdefender_api.NewBitdefenderScannerClient(conn)
	// 	scanFileRequest := &bitdefender_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with bitdefender: :v", err)
	// 		bitdefenderChan <- multiav.ScanResult{}
	// 	} else {
	// 		bitdefenderChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.clamav_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := clamav_api.NewClamAVScannerClient(conn)
	// 	scanFileRequest := &clamav_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with clamav: :v", err)
	// 		clamavChan <- multiav.ScanResult{}
	// 	} else {
	// 		clamavChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.comodo_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := comodo_api.NewComodoScannerClient(conn)
	// 	scanFileRequest := &comodo_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with comodo: :v", err)
	// 		comodoChan <- multiav.ScanResult{}
	// 	} else {
	// 		comodoChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.eset_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := eset_api.NewEsetScannerClient(conn)
	// 	scanFileRequest := &eset_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with eset: :v", err)
	// 		esetChan <- multiav.ScanResult{}
	// 	} else {
	// 		esetChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.fsecure_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := fsecure_api.NewFSecureScannerClient(conn)
	// 	scanFileRequest := &fsecure_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with fsecure: :v", err)
	// 		fsecureChan <- multiav.ScanResult{}
	// 	} else {
	// 		fsecureChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.kaspersky_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := kaspersky_api.NewKasperskyScannerClient(conn)
	// 	scanFileRequest := &kaspersky_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with kaspersky: :v", err)
	// 		kasperskyChan <- multiav.ScanResult{}
	// 	} else {
	// 		kasperskyChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.symantec_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := symantec_api.NewSymantecScannerClient(conn)
	// 	scanFileRequest := &symantec_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with symantec: :v", err)
	// 		symantecChan <- multiav.ScanResult{}
	// 	} else {
	// 		symantecChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	// go func() {
	// 	address := viper.GetString("multiav.windefender_addr")
	// 	conn, err := grpc.Dial(address, []grpc.DialOption{grpc.WithInsecure()}...)
	// 	if err != nil {
	// 		log.Errorf("fail to dial: %v", err)
	// 	}
	// 	defer conn.Close()

	// 	client := windefender_api.NewWinDefenderScannerClient(conn)
	// 	scanFileRequest := &windefender_api.ScanFileRequest{Filepath: filePath}
	// 	res, err := client.ScanFile(context.Background(), scanFileRequest)
	// 	if err != nil {
	// 		log.Errorln("Failed to scan with windows defender: :v", err)
	// 		windefenderChan <- multiav.ScanResult{}
	// 	} else {
	// 		windefenderChan <- multiav.ScanResult{
	// 			Output:   res.Output,
	// 			Infected: res.Infected,
	// 			Update:   res.Update}
	// 	}
	// }()

	avCount := 0
	for {
		select {
		case aviraRes := <-aviraChan:
			multiavScanResults["avira"] = aviraRes
			avCount++
		// case avastRes := <-avastChan:
		// 	multiavScanResults["avast"] = avastRes
		// 	avCount++
		// case bitdefenderRes := <-bitdefenderChan:
		// 	multiavScanResults["bitdefender"] = bitdefenderRes
		// 	avCount++
		// case clamavRes := <-clamavChan:
		// 	multiavScanResults["clamav"] = clamavRes
		// 	avCount++
		// case comodoRes := <-comodoChan:
		// 	multiavScanResults["comodo"] = comodoRes
		// 	avCount++
		// case esetRes := <-esetChan:
		// 	multiavScanResults["eset"] = esetRes
		// 	avCount++
		// case fsecureRes := <-fsecureChan:
		// 	multiavScanResults["fsecure"] = fsecureRes
		// 	avCount++
		// case kasperskyRes := <-kasperskyChan:
		// 	multiavScanResults["kaspersky"] = kasperskyRes
		// 	avCount++
		// case symantecRes := <-symantecChan:
		// 	multiavScanResults["symanetc"] = symantecRes
		// 	avCount++
		// case windefenderRes := <-windefenderChan:
		// 	multiavScanResults["windefender"] = windefenderRes
		// 	avCount++
		// case dummyRes := <-dummyChan:
		// 	multiavScanResults["dummy"] = dummyRes
		// 	avCount++
		}

		if avCount == avEnginesCount {
			break
		}
	}

	return multiavScanResults
}
