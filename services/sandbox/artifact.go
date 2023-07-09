// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"errors"
	"strings"

	pb "github.com/saferwall/saferwall/internal/agent/proto"
)

// Artifact represents dumped memory buffers (during process injection, memory
// decryption, or anything alike) and files dropped by the sample.
type Artifact struct {
	// File  name of the artifact.
	// * Memory buffers are in this format: PID.TID-VA-BuffSize-API.membuff
	//  -> 2E00.A60-0x1A46D880000-77824-Free.membuff
	//  -> 2E00.A60-0x1A46D8A0000-12824-Crypto.membuff
	//  -> 2E00.A60-0x1A46D9B0000-12824-WriteProcess.membuff
	// * Files dropped are in this format: FilePath-Size.filecreate
	//  -> C##ProgramData##Delete.vbs-9855.filecreate
	Name string `json:"name"`
	// The binary content of the artifact.
	Content []byte `json:"-"`
	// The artifact kind: membuff, filecreate, ..
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
	parts := strings.Split(name, ",")
	if len(parts) < 5 {
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

		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}
