package goyara

import (
	"path"
	"testing"
)

const (
	yaraRulesPath = "../../build/data/rules/"
)

func TestYara(t *testing.T) {
	t.Run("TestYaraLoadRules", func(t *testing.T) {
		rules := []Rule{
			{
				Namespace: "capabilities",
				Filename:  path.Join(yaraRulesPath, "Capabilities/capabilities.yar"),
			}, {
				Namespace: "crypto",
				Filename:  path.Join(yaraRulesPath, "Crypto/crypto_signatures.yar"),
			},
		}
		_, err := Load(rules)
		if err != nil {
			t.Fatal("failed to load yara rules with error :", err)
		}

	})
	t.Run("TestYaraScanFile", func(t *testing.T) {
		rules := []Rule{
			{
				Namespace: "capabilities",
				Filename:  path.Join(yaraRulesPath, "Capabilities/capabilities.yar"),
			}, {
				Namespace: "crypto",
				Filename:  path.Join(yaraRulesPath, "Crypto/crypto_signatures.yar"),
			},
		}
		r, err := Load(rules)
		if err != nil {
			t.Fatal("failed to load yara rules with error :", err)
		}
		_, err = ScanFile(r, "../../testdata/putty.exe")
		if err != nil {
			t.Fatal("failed to scan file with error :", err)
		}
	})
}
