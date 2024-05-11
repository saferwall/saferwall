// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import "github.com/saferwall/saferwall/internal/behavior"

// Capability represents any capability found in executable files.
// An example of a capability is: Exfiltration over C2 server.
type Capability struct {
	// Description describes in a few words the capability.
	Description string `json:"description"`
	// The severity of the capability: low, suspicious, high, etc.
	Severity string `json:"severity"`
	// Category of the capability: persistence, anti-analysis, etc.
	Category string `json:"category"`
	// The module that generated the capability: yara, behavior.
	Module string `json:"module"`
	// Rule ID which matched.
	RuleID string `json:"rule_id"`
	// Process identifier responsible for generating this capability.
	ProcessID string `json:"pid"`
	// Optional field indicating the malware family name.
	Family string `json:"-"`
}

func generateCapabilities(bhvRules []behavior.MatchRule, yaraRules []MatchRule) []Capability {
	capabilities := make([]Capability, 0)
	for _, rule := range bhvRules {
		// Add family to the behavior rules.
		capabilities = append(capabilities, Capability{
			Description: rule.Description,
			Severity:    rule.Severity,
			Category:    rule.Category,
			RuleID:      rule.ID,
			ProcessID:   rule.ProcessID,
			Module:      "behavior",
		})
	}

	for _, rule := range yaraRules {
		capabilities = append(capabilities, Capability{
			Description: rule.Description,
			Severity:    rule.Severity,
			Category:    rule.Category,
			RuleID:      rule.ID,
			ProcessID:   rule.ProcessID,
			Family:      rule.Family,
			Module:      "yara",
		})
	}

	return capabilities
}
