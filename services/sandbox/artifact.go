// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"errors"
	"strings"

	pb "github.com/saferwall/saferwall/internal/agent/proto"
	"github.com/saferwall/saferwall/internal/magic"
)

// ArtifactType represents the type of the artifact.
type ArtifactType string

const (
	// Freed memory regions: ProcessName__PID__TID__VA__ID.memfree
	// example: svchost.exe__0x2E00__0xA60__0x1A46D880000__5a0a1add.memfree
	MemFree ArtifactType = "memfree"

	// Files created: ProcessName__PID__FilePath.filecreate
	// example: explorer.exe__0x2E00__C##ProgramData##Delete.vbs.filecreate
	FileCreate ArtifactType = "filecreate"

	// Code injection: ProcessName__PID__RemoteProcessName__RemotePID.codeinject
	// example: emotet.exe__0x2E00__svchost.exe__0x3004.codeinject
	CodeInjection ArtifactType = "codeinject"

	// Memory dumps: ProcessName__PID__TID__VA.memdmp
	// example: explorer.exe__0x2E00__0x400000__.memdmp
	MemDmp ArtifactType = "memdmp"
)

const (
	maxArtifactCount = 10
)

// Artifact represents an extracted artifact during the dynamic analysis.
type Artifact struct {
	// File  name of the artifact.
	Name string `json:"name"`
	// The binary content of the artifact.
	Content []byte `json:"-"`
	// The artifact kind: memfree, filecreate, ..
	Kind ArtifactType `json:"kind"`
	// The SHA256 hash of the artifact.
	SHA256 string `json:"sha256"`
	// Detection contains the family name of the malware if it is malicious,
	// or clean otherwise.
	Detection string `json:"detection"`
	// List of all matched rules.
	MatchedRules []string `json:"matched_rules"`
	// The file type, i.e docx, dll, etc.
	FileType string `json:"file_type"`
}

// Extract the kind of the artifact from the artifact name.
func deduceKindFromName(name string) (ArtifactType, error) {
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid artifact name")
	}
	kind := parts[len(parts)-1]
	return ArtifactType(kind), nil
}

// Generate artifacts metadata.
func (s *Service) generateArtifacts(resArtifacts []*pb.AnalyzeFileReply_Artifact) ([]Artifact, error) {
	artifacts := []Artifact{}
	var err error
	for _, art := range resArtifacts {

		artifact := Artifact{
			Name:    art.GetName(),
			Content: art.GetContent(),
		}

		artifact.Kind, err = deduceKindFromName(artifact.Name)
		if err != nil {
			s.logger.Errorf("failed to deduce artifact kind from %s", artifact.Name)
			continue
		}

		// SHA256.
		artifact.SHA256 = s.hasher.Hash(artifact.Content)

		// Yara scan.
		matches, err := s.yaraScanner.ScanBytes(artifact.Content)
		if err != nil {
			s.logger.Errorf("failed to scan artifact with yara: %s", artifact.Name)
		}
		if len(matches) > 0 {
			artifact.MatchedRules = s.yaraScanner.StringifyMatches(matches)
			s.logger.Infof("yara rules matches: %v", artifact.MatchedRules)
		} else {
			artifact.MatchedRules = make([]string, 0)
		}

		// File type.
		artifact.FileType, err = magic.ScanBytes(artifact.Content)
		if err != nil {
			s.logger.Errorf("failed to detect file type from %s", artifact.Name)
		}

		// Pick the most representative detection name.
		for _, match := range matches {
			artifact.Detection = match.Rule
			break
		}

		artifacts = append(artifacts, artifact)
	}

	s.logger.Infof("artifact count is %d, max: %d",
		len(artifacts), maxArtifactCount)
	var curatedArtifacts []Artifact
	if len(artifacts) >= maxArtifactCount {
		for c, artifact := range artifacts {
			if c >= maxArtifactCount && len(artifact.MatchedRules) == 0 {
				continue
			}
			curatedArtifacts = append(curatedArtifacts, artifact)
		}
	}

	return curatedArtifacts, nil
}
