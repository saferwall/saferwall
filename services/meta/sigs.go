// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package meta

import (
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// Compilers, Installers, Packers names as seen by DiE (Detect It Easy)
// This map a signature name substring into a tag.
var sigMap = map[string]string{
	"Nullsoft Scriptable Install": "nsis",
	"Inno Setup":                  "innosetup",
	"UPX":                         "upx",
	"FSG":                         "fsg",
	"MEW":                         "mew",
	"PKLITE":                      "pklite",
	"ASPack":                      "aspack",
	"RLPack":                      "rlpack",
	"NsPacK":                      "nspack",
	"ASProtect":                   "asprotect",
	"ACProtect":                   "acprotect",
	"AntiDote":                    "antidote",
	"PECompact":                   "pecompact",
	"PECrypt32":                   "pecrypt32",
	"PE-Armor":                    "pe-armor",
	"PESpin":                      "pespin",
	"PEBundle":                    "pebundle",
	"Petite":                      "petite",
	"PELock":                      "pelock",
	"tElock":                      "telock",
	"PKLITE32":                    "pklite",
	"EXECryptor":                  "execryptor",
	"ExeStealth":                  "exestealth",
	"RCryptor":                    "rcryptor",
	"SDProtector":                 "sdprotector",
	"Obsidium":                    "obsidium",
	"VMProtect":                   "vmprotect",
	"Themida/Winlicense":          "themida-winlicense",
	"MoleBox":                     "molebox",
	"ENIGMA":                      "enigma",
	"MPRESS":                      "mpress",
	"NeoLite":                     "neolite",
	"Armadillo":                   "armadillo",
	"Krypton":                     "krypton",
	"HidePE":                      "hidepe",
	"MSLRH":                       "mslrh",
	"gcc":                         "gcc",
	"MinGW":                       "mingw",
	"Microsoft Visual C/C++":      "vc",
	"Microsoft Visual Basic":      "vb",
	"PureBasic":                   "purebasic",
	"Borland C++":                 "borland-c++",
	"MASM":                        "masm",
	"FASM":                        "fasm",
	"Library: .NET":               "dotnet",
	"MFC":                         "mfc",
	"Yoda's Crypter":              "yodascrypter",
	"Delphi":                      "delphi",
	"AutoIt":                      "autoit",
	"StarForce":                   "starforce",
	"eXPressor":                   "expressor",
	"sfx: Microsoft Cabinet":      "sfx-cab",
	"sfx: 7-Zip":                  "sfx-7z",
	"sfx: WinZip":                 "sfx-zip",
	"sfx: WinACE":                 "sfx-ace",
	"Smart Assembly":              "smart-assembly",
	".NET Reactor":                "dotnet-reactor",
	"Babel .NET":                  "babel.net",
	"Confuser":                    "confuser",
	"Dotfuscator":                 "dotfuscator",
	"Eazfuscator":                 "eazfuscator",
	"InstallShield":               "installshield",
	"Inno Setup Module":           "innosetup",
	"Enigma Installer":            "enigma-install",
	"FDM Installer":               "fdm-install",
	"Gentee Installer":            "gentee-install",
	"Ghost Installer":             "ghost-install",
}

var typeMap = map[string]string{
	// binary executables
	"PE32":                "pe",
	"ELF":                 "elf",
	"Mach-O":              "mach-o",
	"MS Windows shortcut": "lnk",
	"MS-DOS":              "msdos",

	// documents
	"PDF document":                        "pdf",
	"Rich Text Format":                    "rtf",
	"Microsoft Word 2007+":                "ooxml",
	"Microsoft Excel 2007+":               "ooxml",
	"Microsoft PowerPoint 2007+":          "ooxml",
	"Composite Document File V2 Document": "ole2",

	// images and media
	"PC bitmap":           "bmp",
	"JPEG image data":     "jpeg",
	"PNG image data":      "png",
	"GIF image data":      "gif",
	"SVG Scalable Vector": "svg",
	"Macromedia Flash":    "swf",

	// archives
	"Zip archive data":          "zip",
	"RAR archive data":          "rar",
	"7-zip archive data":        "7-zip",
	"gzip compressed data":      "gzip",
	"bzip2 compressed data":     "bzip2",
	"tar archive":               "tar",
	"XZ compressed data":        "xz",
	"Java archive data (JAR)":   "jar",
	"Microsoft Cabinet archive": "cab",

	// misc
	"ISO 9660 CD-ROM": "iso",

	// text-based: xml, html, js, hta, swf, ....
	"ASCII text":              "txt",
	"Unicode text":            "txt",
	"ISO-8859 text":           "txt",
	"Unicode (with BOM) text": "txt",
}

func guessFileExtension(data []byte, magic string, format string, trid []string) string {

	switch format {
	// binaries
	case "lnk":
		return "lnk"

	// documents
	case "pdf":
		return "pdf"
	case "rtf":
		return "rtf"
	case "ooxml":
		// for now we assume that it is .X, later, we need to do more parsing
		// to figure out the exact extension (docm, dotm, ...)
		if strings.Contains(magic, "Word") {
			return "docx"
		} else if strings.Contains(magic, "Excel") {
			return "xlsx"
		} else if strings.Contains(magic, "PowerPoint") {
			return "pptx"
		}
	case "ole2":
		// same remark as with ooxml.
		if strings.Contains(magic, "Word") {
			return "doc"
		} else if strings.Contains(magic, "Excel") {
			return "xls"
		} else if strings.Contains(magic, "PowerPoint") {
			return "ppt"
		}
		// If file magic does not work, try trid.
		if len(trid) > 0 {
			tridOut := trid[0]
			if strings.Contains(tridOut, "Publisher") {
				return "pub"
			}
		}
	// images and media
	case "bmp":
		return "bmp"
	case "jpeg":
		return "jpeg"
	case "png":
		return "png"
	case "gif":
		return "gif"
	case "svg":
		return "svg"
	case "swf":
		return "swf"

	// archives
	case "zip":
		return "zip"
	case "rar":
		return "rar"
	case "7-zip":
		return "7z"
	case "gzip":
		return "gz"
	case "bzip2":
		return "bz2"
	case "tar":
		return "tar"
	case "xz":
		return "xz"
	case "jar":
		return "jar"
	case "cab":
		return "cab"

	// misc
	case "iso":
		return "iso"

	// txt based files: powershell, batch, html, javascript,
	// vbscript, jscript, wsf, hta

	case "txt":
		// Order matters in this logic.
		if IsWsf(data) {
			return "wsf"
		}

		// HTML or HTA
		if strings.Contains(magic, "HTML document") {
			if IsHTMLApp(data) {
				return "hta"
			}

			return "html"
		}

		lang, _ := enry.GetLanguageByClassifier(data, []string{
			"powershell", "batch", "vbscript", "javascript"})
		switch lang {
		case "PowerShell":
			return "ps1"
		case "JavaScript":
			return "js"
		case "VBScript":
			return "vbs"
		case "Batchfile":
			return "bat"
		}

	}
	return "unknown"
}

func IsHTMLApp(data []byte) bool {

	content := strings.ToLower(string(data))
	content = strings.Join(strings.Fields(content), " ")
	if strings.Contains(content, "<hta:application") ||
		strings.Contains(content, "<script language=") ||
		strings.Contains(content, "activexobject") {
		return true
	}

	return false
}

func IsWsf(data []byte) bool {

	content := strings.ToLower(string(data))
	content = strings.Join(strings.Fields(content), " ")

	if strings.Contains(content, "<job id=") ||
		strings.Contains(content, "<script language=") ||
		strings.Contains(content, "<package>") {
		return true
	}

	return false
}
