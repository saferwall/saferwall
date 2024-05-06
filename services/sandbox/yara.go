// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"errors"
	"strings"

	"github.com/hillu/go-yara/v4"
	pb "github.com/saferwall/saferwall/internal/agent/proto"
)

// MatchRule represents a matched yara rule.
type MatchRule struct {
	// Description describes the purpose of the rule.
	Description string `json:"description"`
	// ID uniquely identify the rule.
	ID string `json:"id"`
	// Category indicates the category of the yara rule.
	// examples include: anti-analysis, ransomware, ..
	Category string `json:"category"`
	// Severity indicates how confident the rule is to classify
	// the threat as malicious.
	Severity string `json:"severity"`
	// Process identifier responsible for matching the rule.
	// This field is not always available as some behavior rules matches over
	// multiple processes.
	ProcessID string `json:"proc_id"`
	// Family describes the family name of the malware.
	Family string `json:"family"`
}

// Extract the process ID responsible for geneting this artifact from its name.
func getPIDFromName(name string) (string, error) {
	parts := strings.Split(name, "__")
	if len(parts) < 2 {
		return "", errors.New("invalid artifact name")
	}
	return parts[1], nil
}

func categoryFromMeta(match yara.MatchRule) string {
	for _, meta := range match.Metas {
		if meta.Identifier == "category" {
			return strings.ToLower(meta.Value.(string))
		}
	}

	return ""
}

func (s *Service) scanArtifactsWithYara(artifacts []*pb.AnalyzeFileReply_Artifact) ([]MatchRule, error) {

	yaraMatches := []MatchRule{}

	for _, artifact := range artifacts {

		matches, err := s.yaraScanner.ScanBytes(artifact.GetContent())
		if err != nil {
			s.logger.Errorf("failed to scan artifact %s with behavior: %v",
				artifact.GetName(), err)
			continue
		}

		if len(matches) == 0 {
			continue
		}

		processID, err := getPIDFromName(artifact.GetName())
		if err != nil {
			s.logger.Errorf("failed to extract artifact %s process ID: %v",
				artifact.GetName(), err)
			continue
		}

		for _, match := range matches {

			yaraMatch := MatchRule{ProcessID: processID}

			for _, meta := range match.Metas {
				if meta.Identifier == "severity" {
					yaraMatch.Severity = strings.ToLower(meta.Value.(string))
				} else if meta.Identifier == "category" {
					yaraMatch.Category = strings.ToLower(meta.Value.(string))
				} else if meta.Identifier == "description" {
					yaraMatch.Description = meta.Value.(string)
				} else if meta.Identifier == "id" {
					yaraMatch.ID = meta.Value.(string)
				} else if meta.Identifier == "malware" {
					yaraMatch.Family = strings.ToLower(meta.Value.(string))
				}
			}

			// Skip duplicated matches.
			found := false
			for _, ym := range yaraMatches {
				if ym == yaraMatch {
					found = true
					break
				}
			}
			if !found {
				yaraMatches = append(yaraMatches, yaraMatch)
			}
		}
	}

	return yaraMatches, nil
}
