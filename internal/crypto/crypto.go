// Copyright 2022 Saferwall. All rights reserved.
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

// Result aggregates all hashes.
type Result struct {
	CRC32  string
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
	SSDeep string
}

// GetCRC32 returns CRC32 checksum in hex format.
func GetCRC32(b []byte) string {
	checksum := crc32.ChecksumIEEE(b)
	h := fmt.Sprintf("0x%x", checksum)
	return h
}

// GetMD5 returns MD5 hash.
func GetMD5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSHA1 returns SHA1 hash.
func GetSHA1(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSHA256 returns SHA256 hash.
func GetSHA256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSHA512 returns SHA512 hash.
func GetSHA512(b []byte) string {
	h := sha512.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// GetSSDeep returns ssdeep fuzzy hash.
func GetSSDeep(b []byte) (string, error) {
	return ssdeep.FuzzyBytes(b)
}

// HashBytes run all crypto modules and return results.
func HashBytes(data []byte) Result {
	FuzzyHash, err := GetSSDeep(data)
	if err != nil && err != ssdeep.ErrFileTooSmall {
		log.Printf("GetSSDeep() failed, got %v", err)
	}
	r := Result{
		CRC32:  GetCRC32(data),
		MD5:    GetMD5(data),
		SHA1:   GetSHA1(data),
		SHA256: GetSHA256(data),
		SHA512: GetSHA512(data),
		SSDeep: FuzzyHash,
	}
	return r
}
