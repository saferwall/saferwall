// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/grpc/multiav"
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
	sophosclient "github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/sophos/proto"
	symantecclient "github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/symantec/proto"
	windefenderclient "github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/client"
	"github.com/saferwall/saferwall/pkg/grpc/multiav/windefender/proto"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	peparser "github.com/saferwall/saferwall/pkg/peparser"
	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/trid"
	"github.com/saferwall/saferwall/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func staticScan(sha256, filePath string, b []byte) result {
	res := result{}
	var err error

	// Crypto Pkg
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

	// Parse PE
	pe, err := parsePE(filePath)
	if err != nil {
		log.Infof("PE parser failed %v", err)
	} else {
		res.PE = pe
		res.Tags = append(res.Tags, "pe")
		log.Infof("PE parser success %s", sha256)
	}

	// Extract tags
	res.GetTags()

	return res
}


func parsePE(filePath string) (peparser.File, error) {

	pe, err := peparser.Open(filePath)
	if err != nil {
		return peparser.File{}, err
	}
	defer pe.Close()

    defer func() { 
        if err := recover(); err != nil {
			log.Printf("PE parser raised an unexpected exception: %v\n", err)
        }
    }()
	
	err = pe.Parse()
	return pe, err
}

func avScan(engine string, filePath string, c chan multiav.ScanResult) {

	// Get the address of AV gRPC server
	key := "multiav." + engine + "_addr"
	address := viper.GetString(key)

	// Get a gRPC client scanner instance for a given engine.
	conn, err := multiav.GetClientConn(address)
	if err != nil {
		log.Printf("Failed to get client conn for [%s]: %v", engine, err)
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
		log.Errorf("Failed to copy the file for engine %s.", engine)
		c <- multiav.ScanResult{}
		return
	}

	filePath = filecopyPath
	res := multiav.ScanResult{}

	switch engine {
	case "avast":
		res, err = avastclient.ScanFile(avast_api.NewAvastScannerClient(conn), filePath)
	case "avira":
		res, err = aviraclient.ScanFile(avira_api.NewAviraScannerClient(conn), filePath)
	case "bitdefender":
		res, err = bitdefenderclient.ScanFile(bitdefender_api.NewBitdefenderScannerClient(conn), filePath)
	case "clamav":
		res, err = clamavclient.ScanFile(clamav_api.NewClamAVScannerClient(conn), filePath)
	case "comodo":
		res, err = comodoclient.ScanFile(comodo_api.NewComodoScannerClient(conn), filePath)
	case "eset":
		res, err = esetclient.ScanFile(eset_api.NewEsetScannerClient(conn), filePath)
	case "fsecure":
		res, err = fsecureclient.ScanFile(fsecure_api.NewFSecureScannerClient(conn), filePath)
	case "kaspersky":
		res, err = kasperskyclient.ScanFile(kaspersky_api.NewKasperskyScannerClient(conn), filePath)
	case "mcafee":
		res, err = mcafeeclient.ScanFile(mcafee_api.NewMcAfeeScannerClient(conn), filePath)
	case "symantec":
		res, err = symantecclient.ScanFile(symantec_api.NewSymantecScannerClient(conn), filePath)
	case "sophos":
		res, err = sophosclient.ScanFile(sophos_api.NewSophosScannerClient(conn), filePath)
	case "windefender":
		res, err = windefenderclient.ScanFile(windefender_api.NewWinDefenderScannerClient(conn), filePath)
	}

	if err != nil {
		log.Errorf("Failed to scan file [%s]: %v", engine, err)
	}
	c <- multiav.ScanResult{Output: res.Output, Infected: res.Infected, Update: res.Update}

	if err = utils.DeleteFile(filecopyPath) ; err != nil {
		log.Errorf("Failed to delete file path %s.", filecopyPath)
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
	avEnginesCount := 12
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
