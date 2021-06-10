// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"runtime/debug"

	"github.com/saferwall/saferwall/pkg/grpc/multiav"
	avastclient "github.com/saferwall/saferwall/pkg/grpc/multiav/avast/client"
	avast_api "github.com/saferwall/saferwall/pkg/grpc/multiav/avast/proto"
	aviraclient "github.com/saferwall/saferwall/pkg/grpc/multiav/avira/client"
	avira_api "github.com/saferwall/saferwall/pkg/grpc/multiav/avira/proto"
	bitdefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/client"
	bitdefender_api "github.com/saferwall/saferwall/pkg/grpc/multiav/bitdefender/proto"
	clamavclient "github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/client"
	clamav_api "github.com/saferwall/saferwall/pkg/grpc/multiav/clamav/proto"
	comodoclient "github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/client"
	comodo_api "github.com/saferwall/saferwall/pkg/grpc/multiav/comodo/proto"
	drwebclient "github.com/saferwall/saferwall/pkg/grpc/multiav/drweb/client"
	drweb_api "github.com/saferwall/saferwall/pkg/grpc/multiav/drweb/proto"
	esetclient "github.com/saferwall/saferwall/pkg/grpc/multiav/eset/client"
	eset_api "github.com/saferwall/saferwall/pkg/grpc/multiav/eset/proto"
	fsecureclient "github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/client"
	fsecure_api "github.com/saferwall/saferwall/pkg/grpc/multiav/fsecure/proto"
	kasperskyclient "github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/client"
	kaspersky_api "github.com/saferwall/saferwall/pkg/grpc/multiav/kaspersky/proto"
	mcafeeclient "github.com/saferwall/saferwall/pkg/grpc/multiav/mcafee/client"
	mcafee_api "github.com/saferwall/saferwall/pkg/grpc/multiav/mcafee/proto"
	sophosclient "github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/client"
	sophos_api "github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/proto"
	symantecclient "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/client"
	symantec_api "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	trendmicroclient "github.com/saferwall/saferwall/pkg/grpc/multiav/trendmicro/client"
	trendmicro_api "github.com/saferwall/saferwall/pkg/grpc/multiav/trendmicro/proto"
	windefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/client"
	windefender_api "github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/proto"
	"github.com/saferwall/saferwall/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func avScan(engine string, filePath string, c chan multiav.ScanResult,
	ctxLogger *log.Entry, cfg *Config) {

	// Fail safe.
	defer func() {
		if r := recover(); r != nil {
			ctxLogger.Errorf("panic occured in av scan: %v", debug.Stack())
		}
	}()

	// Get the address of AV gRPC server.
	address := cfg.MultiAV.Vendors[engine].Address
	enabled := cfg.MultiAV.Vendors[engine].Enabled

	// Is this engine enabled.
	if !enabled {
		c <- multiav.ScanResult{}
		return
	}

	// Get a gRPC client scanner instance for a given engine.
	conn, err := multiav.GetClientConn(address)
	if err != nil {
		ctxLogger.Errorf("Failed to get client conn for [%s]: %v", engine, err)
		c <- multiav.ScanResult{}
		return
	}
	defer conn.Close()

	// Make a copy of the file for each AV engine.
	// This tries to fix the file locking issues which happens
	// if you try to scan a filepath in a nfs share with
	// different engines at the same time.
	filecopyPath := filePath + "-" + engine
	err = utils.CopyFile(filePath, filecopyPath)
	if err != nil {
		ctxLogger.Errorf("Failed to copy the file for engine %s.", engine)
		c <- multiav.ScanResult{}
		return
	}

	filePath = filecopyPath
	res := multiav.ScanResult{}

	switch engine {
	case "avast":
		res, err = avastclient.ScanFile(
			avast_api.NewAvastScannerClient(conn), filePath)
	case "avira":
		res, err = aviraclient.ScanFile(
			avira_api.NewAviraScannerClient(conn), filePath)
	case "bitdefender":
		res, err = bitdefenderclient.ScanFile(
			bitdefender_api.NewBitdefenderScannerClient(conn), filePath)
	case "drweb":
		res, err = drwebclient.ScanFile(
			drweb_api.NewDrWebScannerClient(conn), filePath)
	case "clamav":
		res, err = clamavclient.ScanFile(
			clamav_api.NewClamAVScannerClient(conn), filePath)
	case "comodo":
		res, err = comodoclient.ScanFile(
			comodo_api.NewComodoScannerClient(conn), filePath)
	case "eset":
		res, err = esetclient.ScanFile(
			eset_api.NewEsetScannerClient(conn), filePath)
	case "fsecure":
		res, err = fsecureclient.ScanFile(
			fsecure_api.NewFSecureScannerClient(conn), filePath)
	case "kaspersky":
		res, err = kasperskyclient.ScanFile(
			kaspersky_api.NewKasperskyScannerClient(conn), filePath)
	case "mcafee":
		res, err = mcafeeclient.ScanFile(
			mcafee_api.NewMcAfeeScannerClient(conn), filePath)
	case "symantec":
		res, err = symantecclient.ScanFile(
			symantec_api.NewSymantecScannerClient(conn), filePath)
	case "sophos":
		res, err = sophosclient.ScanFile(
			sophos_api.NewSophosScannerClient(conn), filePath)
	case "trendmicro":
		res, err = trendmicroclient.ScanFile(
			trendmicro_api.NewTrendMicroScannerClient(conn), filePath)
	case "windefender":
		res, err = windefenderclient.ScanFile(
			windefender_api.NewWinDefenderScannerClient(conn), filePath)
	}

	if err != nil {
		ctxLogger.Errorf("Failed to scan file [%s]: %v", engine, err)
	}
	c <- multiav.ScanResult{Enabled: enabled, Output: res.Output,
		Infected: res.Infected, Update: res.Update}

	if utils.Exists(filecopyPath) {
		if err = utils.DeleteFile(filecopyPath); err != nil {
			ctxLogger.Errorf("Failed to delete file path %s.", filecopyPath)
		}
	}
}

func (f *File) multiAvScan(filePath string, cfg *Config,
	ctxLogger *log.Entry) map[string]interface{} {

	// Create channels to receive scan results.
	aviraChan := make(chan multiav.ScanResult)
	avastChan := make(chan multiav.ScanResult)
	bitdefenderChan := make(chan multiav.ScanResult)
	drwebChan := make(chan multiav.ScanResult)
	clamavChan := make(chan multiav.ScanResult)
	comodoChan := make(chan multiav.ScanResult)
	esetChan := make(chan multiav.ScanResult)
	fsecureChan := make(chan multiav.ScanResult)
	kasperskyChan := make(chan multiav.ScanResult)
	mcafeeChan := make(chan multiav.ScanResult)
	symantecChan := make(chan multiav.ScanResult)
	sophosChan := make(chan multiav.ScanResult)
	trendmicroChan := make(chan multiav.ScanResult)
	windefenderChan := make(chan multiav.ScanResult)

	// We Start as much go routines as the AV engines we have.
	// Each go-routines makes a gRPC calls and waits for results.
	// Avast
	go avScan("eset", filePath, esetChan, ctxLogger, cfg)
	go avScan("fsecure", filePath, fsecureChan, ctxLogger, cfg)
	go avScan("avira", filePath, aviraChan, ctxLogger, cfg)
	go avScan("bitdefender", filePath, bitdefenderChan, ctxLogger, cfg)
	go avScan("kaspersky", filePath, kasperskyChan, ctxLogger, cfg)
	go avScan("symantec", filePath, symantecChan, ctxLogger, cfg)
	go avScan("sophos", filePath, sophosChan, ctxLogger, cfg)
	go avScan("windefender", filePath, windefenderChan, ctxLogger, cfg)
	go avScan("clamav", filePath, clamavChan, ctxLogger, cfg)
	go avScan("comodo", filePath, comodoChan, ctxLogger, cfg)
	go avScan("avast", filePath, avastChan, ctxLogger, cfg)
	go avScan("mcafee", filePath, mcafeeChan, ctxLogger, cfg)
	go avScan("drweb", filePath, drwebChan, ctxLogger, cfg)
	go avScan("trendmicro", filePath, trendmicroChan, ctxLogger, cfg)

	multiavScanResults := map[string]interface{}{}
	avEnginesCount := 14
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
		case drwebRes := <-drwebChan:
			multiavScanResults["drweb"] = drwebRes
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
			multiavScanResults["symantec"] = symantecRes
			avCount++
		case sophosRes := <-sophosChan:
			multiavScanResults["sophos"] = sophosRes
			avCount++
		case trendmicroRes := <-trendmicroChan:
			multiavScanResults["trendmicro"] = trendmicroRes
			avCount++
		case windefenderRes := <-windefenderChan:
			multiavScanResults["windefender"] = windefenderRes
			avCount++
		}

		if avCount == avEnginesCount {
			break
		}
	}

	ctxLogger.Debug("multiav scan success")
	return multiavScanResults
}
