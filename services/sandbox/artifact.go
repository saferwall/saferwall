// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

// Artifact represents dumped memory buffers (during process injection, memory
// decryption, or anything alike) and files dropped by the sample.
type Artifact struct {
	// File  name of the artifact.
	// * Memory buffers are in this format: PID.TID-VA-BuffSize-API.membuff
	//  -> 2E00.A60-0x1A46D880000-77824-Free.membuff
	//  -> 2E00.A60-0x1A46D8A0000-12824-Crypto.membuff
	//  -> 2E00.A60-0x1A46D9B0000-12824-WriteProcess.membuff
	// * Files dropped are in this format: PID.TID-FilePath-Size.filecreate
	//  -> 2E00.A60-C##ProgramData##Delete.vbs-9855.filecreate
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
