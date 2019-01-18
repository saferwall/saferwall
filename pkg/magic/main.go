// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package magic

import (
	"github.com/rakyll/magicmime"
)

const (
	// Command to invoke exiftool scanner
	Command = "exiftool"
)

// GetMimeType returns the mime-type from a blob of data.
func GetMimeType(data []byte) (string, error) {

	err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
	if err != nil {
		return "", err
	}

	defer magicmime.Close()

	return magicmime.TypeByBuffer(data)
}
