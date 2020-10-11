// Copyright 2020 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package pe

// Architecture-specific data. This data directory is not used 
// (set to all zeros) for I386, IA64, or AMD64 architecture.
func (pe *File) parseArchitectureDirectory(rva, size uint32) error {
	return nil
}
