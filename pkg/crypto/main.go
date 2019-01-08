// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"log"

	"github.com/glaslos/ssdeep"
)

// Result contains result for all hashes
type Result struct {
	Crc32  string `json:"crc32"`
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
	Sha512 string `json:"sha512"`
	Ssdeep string `json:"ssdeep"`
}

// GetCrc32 returns CRC32 checksum in hex format
func GetCrc32(b []byte) string {
	checksum := crc32.ChecksumIEEE(b)
	h := fmt.Sprintf("0x%x", checksum)
	return h
}

// GetMd5 returns MD5 hash
func GetMd5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSha1 returns SHA1 hash
func GetSha1(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSha256 returns SHA256 hash
func GetSha256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSha512 returns SHA512 hash
func GetSha512(b []byte) string {
	h := sha512.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSsdeep returns ssdeep fuzzy hash
func GetSsdeep(b []byte) (FuzzyHash string, err error) {
	return ssdeep.FuzzyBytes(b)
}

// HashBytes run all crypto modules and return results
func HashBytes(data []byte) Result {
	FuzzyHash, err := GetSsdeep(data)
	if err != nil && err != ssdeep.ErrFileTooSmall {
		log.Printf("GetSsdeep() failed, got %s", err)
	}
	r := Result{
		Crc32:  GetCrc32(data),
		Md5:    GetMd5(data),
		Sha1:   GetSha1(data),
		Sha256: GetSha256(data),
		Sha512: GetSha512(data),
		Ssdeep: FuzzyHash,
	}
	return r
}
