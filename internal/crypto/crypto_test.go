// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package crypto

import (
	"testing"

	"github.com/glaslos/ssdeep"
	"github.com/saferwall/saferwall/internal/utils"
)

var eicarFilePath = "../../testdata/eicar.com"

var crc32tests = []struct {
	in  string
	out string
}{
	{eicarFilePath, "0x6851cf3c"},
}

var md5tests = []struct {
	in  string
	out string
}{
	{eicarFilePath, "44d88612fea8a8f36de82e1278abb02f"},
}

var sha1tests = []struct {
	in  string
	out string
}{
	{eicarFilePath,
		"3395856ce81f2b7382dee72602f798b642f14140"},
}

var sha256tests = []struct {
	in  string
	out string
}{
	{eicarFilePath,
		"275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f"},
}

var sha512tests = []struct {
	in  string
	out string
}{
	{eicarFilePath,
		"cc805d5fab1fd71a4ab352a9c533e65fb2d5b885518f4e565e68847223b8e6b85cb48f3afad842726d99239c9e36505c64b0dc9a061d9e507d833277ada336ab"},
}

var ssdeeptests = []struct {
	in  string
	out string
}{
	{"../../testdata/putty.exe",
		"24576:wpPg/wTlg6Xklt9e/Y/iIpNh6liEmE2CebHNpVffB:XwRg6X+twii8N0oCeLNbfB"},
}

var hashTests = []struct {
	in  string
	out []string
}{
	{
		eicarFilePath,
		[]string{
			"0x6851cf3c",
			"44d88612fea8a8f36de82e1278abb02f",
			"3395856ce81f2b7382dee72602f798b642f14140",
			"275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f",
			"cc805d5fab1fd71a4ab352a9c533e65fb2d5b885518f4e565e68847223b8e6b85cb48f3afad842726d99239c9e36505c64b0dc9a061d9e507d833277ada336ab",
			"",
		},
	},
}

func TestHashBytes(t *testing.T) {

	for _, tt := range hashTests {
		b, _ := utils.ReadAll(tt.in)
		res := HashBytes(b)

		if res.CRC32 != tt.out[0] {
			t.Errorf("TestHashBytes(%s) got %v, want %v", tt.in, res.CRC32, tt.out[0])
		}
		if res.MD5 != tt.out[1] {
			t.Errorf("TestHashBytes(%s) got %v, want %v", tt.in, res.MD5, tt.out[1])
		}
		if res.SHA1 != tt.out[2] {
			t.Errorf("TestHashBytes(%s) got %v, want %v", tt.in, res.SHA1, tt.out[2])
		}
		if res.SHA256 != tt.out[3] {
			t.Errorf("TestHashBytes(%s) got %v, want %v", tt.in, res.SHA256, tt.out[3])
		}
		if res.SHA512 != tt.out[4] {
			t.Errorf("TestHashBytes(%s) got %v, want %v", tt.in, res.SHA512, tt.out[4])
		}
		if res.SSDeep != tt.out[5] {
			t.Errorf("TestHashBytes(%s) got %v, want %v", tt.in, res.SSDeep, tt.out[5])
		}
	}
}

func TestGetCRC32(t *testing.T) {
	for _, tt := range crc32tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetCRC32 failed, got: %s", err)
			}
			got := GetCRC32(b)
			if got != tt.out {
				t.Errorf("TestGetCRC32(%s) got %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

func TestGetMD5(t *testing.T) {
	for _, tt := range md5tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetMD5 failed, got: %s", err)
			}
			got := GetMD5(b)
			if got != tt.out {
				t.Errorf("TestGetMD5(%s) got %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

func TestGetSHA1(t *testing.T) {
	for _, tt := range sha1tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetSHA1 failed, got: %s", err)
			}
			got := GetSHA1(b)
			if got != tt.out {
				t.Errorf("TestGetSHA1(%s) got %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

func TestGetSHA256(t *testing.T) {
	for _, tt := range sha256tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetSHA256 failed, got: %s", err)
			}
			got := GetSHA256(b)
			if got != tt.out {
				t.Errorf("TestGetSHA256(%s) got %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

func TestGetSHA512(t *testing.T) {
	for _, tt := range sha512tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetSHA512 failed, got: %s", err)
			}
			got := GetSHA512(b)
			if got != tt.out {
				t.Errorf("TestGetSHA512(%s) got %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}

func TestGetSSDeep(t *testing.T) {
	for _, tt := range ssdeeptests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := utils.ReadAll(tt.in)
			if err != nil {
				t.Errorf("TestGetSSDeep failed, got: %s", err)
			}
			got, err := GetSSDeep(b)
			if err != nil && err != ssdeep.ErrFileTooSmall {
				t.Errorf("TestGetSSDeep failed, got: %s", err)
			}
			if got != tt.out {
				t.Errorf("TestGetSSDeep(%s) got %v, want %v", tt.in, got, tt.in)
			}
		})
	}
}
