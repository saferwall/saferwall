// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"errors"
	"strings"

	pb "github.com/saferwall/saferwall/internal/agent/proto"
)

// Artifact represents an extracted artifact during the dynamic analysis
// like a memory dumps or written files
type Artifact struct {
	// File  name of the artifact.
	// * Memory buffers are in this format: ProcessName__PID__TID__VA__BuffSize.memfree
	// *** svchost.exe__0x2E00__0xA60__0x1A46D880000__0x77824.memfree
	// * Files dropped are in this format: ProcessName_PID_TID_FilePath__FileSize.filecreate
	// *** explorer.exe__0x2E00__0xA60__C##ProgramData##Delete.vbs__0x9855.filecreate
	Name string `json:"name"`
	// The binary content of the artifact.
	Content []byte `json:"-"`
	// The artifact kind: memfree, filecreate, ..
	Kind string `json:"kind"`
	// The SHA256 hash of the artifact.
	SHA256 string `json:"sha256"`
	// Detection contains the family name of the malware if it is malicious,
	// or clean otherwise.
	Detection string `json:"detection"`
	// The file type, i.e docx, dll, etc.
	FileType string `json:"file_type"`
}

// extract the kind of the artifact from the artifact name.
func deduceKindFromName(name string) (string, error) {
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid artifact name")
	}
	return parts[len(parts)-1], nil
}

// generate artifacts metadata.
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

		artifact.SHA256 = s.hasher.Hash(artifact.Content)

		matches, err := s.yaraScanner.ScanBytes(artifact.Content)
		if err != nil {
			s.logger.Errorf("failed to scan artifact with yara: %s", artifact.Name)
			continue
		}

		for _, match := range matches {
			s.logger.Infof("yara rules matches: %v", match.Rule)
		}

		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}
