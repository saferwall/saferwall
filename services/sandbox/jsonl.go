// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"bytes"
	"encoding/json"
	"io"
)

func (s *Service) jsonl2json(data []byte) []byte {

	reader := bytes.NewReader(data)
	d := json.NewDecoder(reader)
	var jsonData []interface{}

	for {
		// Decode one JSON document.
		var v interface{}
		err := d.Decode(&v)
		if err != nil {
			// io.EOF is expected at end of stream.
			if err != io.EOF {
				s.logger.Error(err)
			}
			break
		}

		// Do something with the value.
		jsonData = append(jsonData, v)
	}

	return toJSON(jsonData)
}
