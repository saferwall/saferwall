package token

import "testing"

func TestKnowsAllModifiers(t *testing.T) {
	testCases := []struct {
		input    string
		expected Kind
	}{
		{"size", FileSize},
		{"type", FileType},
		{"extension", FileExtension},
		{"name", FileName},
		{"positives", Positives},
		{"trid", Trid},
		{"packer", Packer},
		{"magic", FileMagic},
		{"tag", Tag},
		{"fs", FirstSeen},
		{"ls", LastScanned},
		{"crc32", CRC32},
		{"avast", Avast},
		{"avira", Avira},
		{"bitdefender", Bitdefender},
		{"clamav", Clamav},
		{"comodo", Comodo},
		{"drweb", DrWeb},
		{"eset", Eset},
		{"fsecure", FSecure},
		{"kaspersky", Kaspersky},
		{"mcafee", McAfee},
		{"sophos", Sophos},
		{"symantec", Symantec},
		{"trendmicro", TrendMicro},
		{"windefender", Windefender},
		{"md5", MD5},
		{"sha1", SHA1},
		{"sha256", SHA256},
		{"sha512", SHA512},
		{"ssdeep", SSDeep},
	}

	for _, tc := range testCases {
		modifier, ok := GetModifier(tc.input)
		if !ok || modifier != tc.expected {
			t.Fatal("expected modifier to exist for :", tc.input)
		}
	}
}
