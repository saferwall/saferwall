// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"time"

	bs "github.com/saferwall/saferwall/pkg/bytestats"

	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
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
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	peparser "github.com/saferwall/saferwall/pkg/peparser"
	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/trid"
	"github.com/saferwall/saferwall/pkg/utils"
	"github.com/spf13/viper"
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}
type result struct {
	Md5         string                 `json:"md5,omitempty"`
	Sha1        string                 `json:"sha1,omitempty"`
	Sha256      string                 `json:"sha256,omitempty"`
	Sha512      string                 `json:"sha512,omitempty"`
	Ssdeep      string                 `json:"ssdeep,omitempty"`
	Crc32       string                 `json:"crc32,omitempty"`
	Magic       string                 `json:"magic,omitempty"`
	Size        int64                  `json:"size,omitempty"`
	Exif        map[string]string      `json:"exif,omitempty"`
	TriD        []string               `json:"trid,omitempty"`
	Tags        map[string]interface{} `json:"tags,omitempty"`
	Packer      []string               `json:"packer,omitempty"`
	LastScanned *time.Time             `json:"last_scanned,omitempty"`
	Strings     []stringStruct         `json:"strings,omitempty"`
	MultiAV     map[string]interface{} `json:"multiav,omitempty"`
	Status      int                    `json:"status,omitempty"`
	PE          *peparser.File          `json:"pe,omitempty"`
	Histogram   []int                  `json:"histogram,omitempty"`
	ByteEntropy []int                  `json:"byte_entropy,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

func (res *result) parseFile(b []byte, filePath string) {
	// Get the file type using linux magic utility.
	magic := res.Magic
	if strings.HasPrefix(magic, "PE32") {
		res.Type = "pe"
	} else if strings.HasPrefix(magic, "XML") {
		res.Type = "xml"
	} else if strings.HasPrefix(magic, "HTML") {
		res.Type = "html"
	} else if strings.HasPrefix(magic, "ELF") {
		res.Type = "elf"
	} else if strings.HasPrefix(magic, "Macromedia Flash") {
		res.Type = "swf"
	}

	// Parse it accrording to its type.
	var err error
	switch res.Type {
	case "pe":
		res.PE, err = parsePE(filePath)
		if err != nil {
			contextLogger.Errorf("pe parser failed: %v", err)
		}

		// Extract Byte Histogram and byte entropy.
		res.Histogram = bs.ByteHistogram(b)
		res.ByteEntropy = bs.ByteEntropyHistogram(b)
		contextLogger.Debug("bytestats pkg success")
	}
}

func staticScan(sha256, filePath string, b []byte) result {
	res := result{}
	var err error

	// Size
	res.Size = int64(len(b))

	// Calculates hashes.
	r := crypto.HashBytes(b)
	res.Crc32 = r.Crc32
	res.Md5 = r.Md5
	res.Sha1 = r.Sha1
	res.Sha256 = r.Sha256
	res.Sha512 = r.Sha512
	res.Ssdeep = r.Ssdeep
	contextLogger.Debug("crypto pkg success")

	// Get exif metadata.
	res.Exif, err = exiftool.Scan(filePath)
	if err != nil {
		contextLogger.Errorf("exiftool pkg failed with: %v", err)
	}
	contextLogger.Debug("exiftool pkg success")

	// Get TriD.
	res.TriD, err = trid.Scan(filePath)
	if err != nil {
		contextLogger.Errorf("trid pkg failed with: %v", err)
	}
	contextLogger.Debug("trid pkg success")

	// Get magic.
	res.Magic, err = magic.Scan(filePath)
	if err != nil {
		contextLogger.Errorf("magic pkg failed with: %v", err)
	}
	contextLogger.Debug("magic pkg success")

	// Get DiE
	res.Packer, err = packer.Scan(filePath)
	if err != nil {
		contextLogger.Errorf("die pkg failed with: %v", err)
	}
	contextLogger.Debug("die pkg success")

	// Extract strings.
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
	contextLogger.Debug("strings pkg success")

	// Run the parsers
	res.parseFile(b, filePath)

	return res
}

func parsePE(filePath string) (*peparser.File, error) {

	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	// Open the file and prepare it to be parsed.
	opts := peparser.Options{SectionEntropy: true}
	pe, err := peparser.New(filePath, &opts)
	if err != nil {
		return nil, err
	}
	defer pe.Close()

	// Do the actual parsing
	err = pe.Parse()

	// Dirty hack to fix json marshalling
	for i := range pe.Certificates.Content.Certificates {
		pe.Certificates.Content.Certificates[i].PublicKey = nil
	}

	contextLogger.Info("pe pkg success")
	return pe, err
}

func avScan(engine string, filePath string, c chan multiav.ScanResult) {

	// Get the address of AV gRPC server
	multiavCfg := viper.GetStringMap("multiav")
	engineCfg := multiavCfg[engine]
	address := engineCfg.(map[string]interface{})["addr"].(string)
	enabled := engineCfg.(map[string]interface{})["enabled"].(bool)

	// Is this engine enabled
	if !enabled {
		c <- multiav.ScanResult{}
		return
	}

	// Get a gRPC client scanner instance for a given engine.
	conn, err := multiav.GetClientConn(address)
	if err != nil {
		contextLogger.Errorf("Failed to get client conn for [%s]: %v", engine, err)
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
		contextLogger.Errorf("Failed to copy the file for engine %s.", engine)
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
	case "drweb":
		res, err = drwebclient.ScanFile(drweb_api.NewDrWebScannerClient(conn), filePath)
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
	case "trendmicro":
		res, err = trendmicroclient.ScanFile(trendmicro_api.NewTrendMicroScannerClient(conn), filePath)
	case "windefender":
		res, err = windefenderclient.ScanFile(windefender_api.NewWinDefenderScannerClient(conn), filePath)
	}

	if err != nil {
		contextLogger.Errorf("Failed to scan file [%s]: %v", engine, err)
	}
	c <- multiav.ScanResult{Enabled: enabled, Output: res.Output, Infected: res.Infected, Update: res.Update}

	if err = utils.DeleteFile(filecopyPath); err != nil {
		contextLogger.Errorf("Failed to delete file path %s.", filecopyPath)
	}
}

func multiAvScan(filePath string) map[string]interface{} {

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
	go avScan("drweb", filePath, drwebChan)
	go avScan("trendmicro", filePath, trendmicroChan)

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

	contextLogger.Debug("multiav pkg success")
	return multiavScanResults
}
