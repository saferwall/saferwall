// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package avlabel

import "testing"

func TestParseWindefender(t *testing.T) {
	t.Run("TestParseWindefender", func(t *testing.T) {
		testCases := []struct {
			in  string
			out Detection
		}{
			{
				in:  "ALF:HeraklezEval:TrojanSpy:AndroidOS/Anubis.B!rfn",
				out: Detection{},
			},
			{
				in:  "!#HSTR:AntiDisasmTrick",
				out: Detection{},
			},
			{
				in:  "!#FOP:Ogimant!Obfuscator!Acv!Ep",
				out: Detection{},
			},
			{
				in:  "!#Lowfi:AGGREGATOR:BanloadXorPrev6",
				out: Detection{},
			},
			{
				in: "Virus:DOS/4368",
				out: Detection{Category: "Virus", Platform: "DOS",
					Family: "4368"},
			},
			{
				in: "Virus:DOS/Screaming_Fist.894.dr",
				out: Detection{Category: "Virus", Platform: "DOS",
					Family: "Screaming_Fist", Variant: ".894.dr"},
			},
			{
				in: "Virus:Win32/Bolzano_4096.A",
				out: Detection{Category: "Virus", Platform: "Win32",
					Family: "Bolzano_4096", Variant: ".A"},
			},
			{
				in: "Worm:Win32/Buchon.F@mm",
				out: Detection{Category: "Worm", Platform: "Win32",
					Family: "Buchon", Variant: ".F@mm"},
			},
			{
				in: "Worm:Win32/Korgo.AZ.dam#2",
				out: Detection{Category: "Worm", Platform: "Win32",
					Family: "Korgo", Variant: ".dam#2"},
			},
			{
				in: "Constructor:Win32/Smoothie_2_001",
				out: Detection{Category: "Constructor", Platform: "Win32",
					Family: "Smoothie_2_001",
				},
			},
			{
				in: "Trojan:MacOS/Exploit.CVE-2007-6166!MTB",
				out: Detection{Category: "Trojan", Platform: "MacOS",
					Family: "Exploit", Variant: ".CVE-2007-6166!MTB"},
			},
			{
				in: "TrojanDownloader:MSIL/Gendwnurl.BL!bit",
				out: Detection{Category: "TrojanDownloader", Platform: "MSIL",
					Family: "Gendwnurl", Variant: ".BL!bit"},
			},
			{
				in: "Exploit:O97M/CVE-2017-11882.PDD!MTB",
				out: Detection{Category: "Exploit", Platform: "O97M",
					Family: "CVE-2017-11882", Variant: ".PDD!MTB"},
			},
			{
				in: "Spammer:Win32/MailSender.2_2",
				out: Detection{Category: "Spammer", Platform: "Win32",
					Family: "MailSender", Variant: ".2_2"},
			},
		}

		for _, tt := range testCases {
			got := ParseWindefender(tt.in)
			if got.Family != tt.out.Family {
				t.Errorf("TestParseWindefender(%s) got %v, want %v", tt.in, got, tt.out)
			}
		}
	})
}
