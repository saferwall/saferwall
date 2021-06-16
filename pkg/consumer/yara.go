// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"path"

	"github.com/hillu/go-yara/v4"
	goyara "github.com/saferwall/saferwall/pkg/yara"
)

const (
	// YaraRulesPath is the OS level (inside docker) path to Yara rules.
	yaraRulesPath = "/opt/yararules"
)

// LoadYaraRules will read yara rules in-memory.
func LoadYaraRules() (*yara.Rules, error) {
	yaraRules := []goyara.Rule{
		{
			Namespace: "antidebug_antivm",
			Filename:  path.Join(yaraRulesPath, "antidebug_antivm_index.yar"),
		},
		{
			Namespace: "capabilities",
			Filename:  path.Join(yaraRulesPath, "capabilities_index.yar"),
		}, {
			Namespace: "crypto",
			Filename:  path.Join(yaraRulesPath, "crypto_index.yar"),
		}, {
			Namespace: "packers",
			Filename:  path.Join(yaraRulesPath, "packers_index.yar"),
		},
	}
	rules, err := goyara.Load(yaraRules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}
