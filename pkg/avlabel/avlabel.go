// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avlabel

import (
	"regexp"
)

type Detection struct {
	Family   string
	Category string
	Platform string
	Variant  string
}

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
func ParseWindefender(detection string) Detection {
	// Backdoor:Win32/Beastdoor.DQ
	// Exploit:O97M/CVE-2017-11882.M
	params := getParams(`^(?P<Category>[a-zA-Z]{1,20})\:(?P<Platform>[a-zA-Z0-9]{1,20})\/(?P<Family>[\w-]{1,20})(?P<Variant>.*)$`, detection)
	return Detection{
		Category: params["Category"],
		Platform: params["Platform"],
		Family:   params["Family"],
		Variant:  params["Variant"],
	}
}

// ParseEset parse.
func ParseEset(detection string) Detection {
	// Win32/Yurist, Win32/Agobot, IRC/SdBot
	params := getParams(`^(?P<Platform>[a-zA-Z0-9]{3,10})\/(?P<Family>[a-zA-Z0-9]{1,20})$`, detection)

	if len(params) != 0 {
		return Detection{
			Platform: params["Platform"],
			Family:   params["Family"],
			Variant:  params["Variant"],
		}
	}

	// Win32/Agent.ODC, Win32/Injector.DXDY, Win32/Dridex.BC
	params = getParams(`^(?P<Platform>[a-zA-Z0-9]{3,10})\/(?P<Family>[a-zA-Z0-9]{1,20})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	if len(params) != 0 {
		return Detection{
			Platform: params["Platform"],
			Family:   params["Family"],
			Variant:  params["Variant"],
		}
	}

	// Android/TrojanDropper.Agent.BII, Win32/PSW.OnLineGames.NMY, Win32/Adware.FakeAV.R, Win32/Filecoder.CryptoWall.D
	params = getParams(`^(?P<Platform>[a-zA-Z0-9]{3,10})\/(?P<Category>[a-zA-Z0-9]{1,20})\.(?P<Family>[a-zA-Z0-9]{1,10})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	return Detection{
		Platform: params["Platform"],
		Category: params["Category"],
		Family:   params["Family"],
		Variant:  params["Variant"],
	}
}

// ParseAvira parses.
func ParseAvira(detection string) Detection {
	// TR/PSW.Tepfer.ockxa, TR/Patched.Ren.Gen, TR/Crypt.XPACK.Gen
	// TR/AD.Kreen.blxny
	params := getParams(`^(?P<Category>[a-zA-Z0-9]{2,10})\/(?P<Type>[a-zA-Z0-9]{1,20})\.(?P<Family>[a-zA-Z0-9]{1,10})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	if len(params) != 0 {
		return Detection{
			Category: params["Category"],
			Platform: params["Type"],
			Family:   params["Family"],
			Variant:  params["Variant"],
		}
	}

	// HEUR/AGEN.1012588, Linux/Mirai.bonb, KIT/Exploit-M-022.B
	// EXP/CVE-2017-11882.Gen
	params = getParams(`^(?P<Category>[a-zA-Z0-9]{2,10})\/(?P<Family>[a-zA-Z0-9-]{1,20})\.(?P<Variant>[a-zA-Z0-9]{1,10})$`, detection)
	return Detection{
		Category: params["Category"],
		Family:   params["Family"],
		Variant:  params["Variant"],
	}
}

// Parse parses an Antivirus detection name.
func Parse(av, detection string) Detection {
	switch av {
	case "windefender":
		return ParseWindefender(detection)
	case "eset":
		return ParseEset(detection)
	case "avira":
		return ParseAvira(detection)
	}

	return Detection{}
}
