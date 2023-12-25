// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package goyara

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func getAbsoluteFilePath(testfile string) string {
	_, p, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(p), "..", "..", testfile)
}

var (
	yaraRulesPath = getAbsoluteFilePath("testdata/yara-rules/")
)

func TestYara(t *testing.T) {
	t.Run("TestYaraLoadRules", func(t *testing.T) {
		rules := []Rule{
			{
				Namespace: "default",
				Filename:  path.Join(yaraRulesPath, "index.yara"),
			},
		}
		_, err := Load(rules)
		if err != nil {
			t.Fatal("failed to load yara rules with error :", err)
		}

	})
	t.Run("TestYaraScanFile", func(t *testing.T) {
		s, err := New(yaraRulesPath)
		if err != nil {
			t.Fatal("failed to create a new scanner with error :", err)
		}
		_, err = s.ScanFile("../../testdata/putty.exe")
		if err != nil {
			t.Fatal("failed to scan file with error :", err)
		}
	})
}
