// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avlabel

import (
	"regexp"
)

func getParams(regEx, url string) (paramsMap map[string]string) {

    var compRegEx = regexp.MustCompile(regEx)
    match := compRegEx.FindStringSubmatch(url)

    paramsMap = make(map[string]string)
    for i, name := range compRegEx.SubexpNames() {
        if i > 0 && i <= len(match) {
            paramsMap[name] = match[i]
        }
    }
    return
}

// ParseWindefender parse.
func ParseWindefender (detection string) map[string]string {
	// Backdoor:Win32/Beastdoor.DQ
	// Exploit:O97M/CVE-2017-11882.M
	params := getParams(`^(?P<Category>[a-zA-Z]{1,20})\:(?P<Platform>[a-zA-Z0-9]{1,20})\/(?P<Family>[a-zA-Z0-9-]{1,20})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	return params
}

// ParseEset parse.
func ParseEset (detection string) map[string]string {
	// Win32/Yurist, Win32/Agobot, IRC/SdBot
	params := getParams(`^(?P<Platform>[a-zA-Z0-9]{3,10})\/(?P<Family>[a-zA-Z0-9]{1,20})$`, detection)
	if len(params) != 0 {
		return params
	}

	// Win32/Agent.ODC, Win32/Injector.DXDY, Win32/Dridex.BC
	params = getParams(`^(?P<Platform>[a-zA-Z0-9]{3,10})\/(?P<Family>[a-zA-Z0-9]{1,20})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	if len(params) != 0 {
		return params
	}

	// Android/TrojanDropper.Agent.BII, Win32/PSW.OnLineGames.NMY, Win32/Adware.FakeAV.R, Win32/Filecoder.CryptoWall.D
	params = getParams(`^(?P<Platform>[a-zA-Z0-9]{3,10})\/(?P<Category>[a-zA-Z0-9]{1,20})\.(?P<Family>[a-zA-Z0-9]{1,10})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	return params
}

// ParseAvira parses.
func ParseAvira (detection string) map[string]string {
	// TR/PSW.Tepfer.ockxa, TR/Patched.Ren.Gen, TR/Crypt.XPACK.Gen
	// TR/AD.Kreen.blxny
	params := getParams(`^(?P<Category>[a-zA-Z0-9]{2,10})\/(?P<Type>[a-zA-Z0-9]{1,20})\.(?P<Family>[a-zA-Z0-9]{1,10})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	if len(params) != 0 {
		return params
	}

	// HEUR/AGEN.1012588, Linux/Mirai.bonb, KIT/Exploit-M-022.B
	// EXP/CVE-2017-11882.Gen
	params = getParams(`^(?P<Category>[a-zA-Z0-9]{2,10})\/(?P<Family>[a-zA-Z0-9-]{1,20})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	return params
}