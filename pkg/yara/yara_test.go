package goyara

import (
	"testing"
)

func TestYara(t *testing.T) {
	t.Run("TestYaraLoadRules", func(t *testing.T) {
		malIndexRule := Rule{
			Namespace: "malware index",
			Filename:  "../../build/data/rules/Capabilities/capabilities.yar",
		}
		_, err := Load([]Rule{malIndexRule})
		if err != nil {
			t.Fatal("failed to load yara rules with error :", err)
		}

	})
	t.Run("TestYaraScanFile", func(t *testing.T) {
		malIndexRule := Rule{
			Namespace: "malware index",
			Filename:  "../../build/data/rules/Capabilities/capabilities.yar",
		}
		r, err := Load([]Rule{malIndexRule})
		if err != nil {
			t.Fatal("failed to load yara rules with error :", err)
		}
		m, err := ScanFile(r, "../../testdata/putty.exe")
		if err != nil {
			t.Fatal("failed to scan file with error :", err)
		}
		for _, match := range m {
			t.Log(match.Rule)
		}
	})
}
