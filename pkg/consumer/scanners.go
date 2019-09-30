package main

import (
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/packer"
	"github.com/saferwall/saferwall/pkg/utils"
	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/trid"
	log "github.com/sirupsen/logrus"

)

func staticScan(sha256, filePath string, b []byte) (result) {
	res := result {}
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

	return res
}

