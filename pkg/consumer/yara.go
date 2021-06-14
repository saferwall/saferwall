package consumer

import (
	"github.com/hillu/go-yara/v4"
	goyara "github.com/saferwall/saferwall/pkg/yara"
)

const (
	// YaraRulesPath is the OS level (inside docker) path to Yara rules.
	YaraRulesPath = "/opt/yararules/rules"
)

// LoadYaraRules will read yara rules in-memory.
func LoadYaraRules() (*yara.Rules, error) {
	yaraRules := []goyara.Rule{
		{
			Namespace: "capabilities",
			Filename:  YaraRulesPath + "/Capabilities/capabilities.yar",
		},
	}

	rules, err := goyara.Load(yaraRules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}
