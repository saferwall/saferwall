// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"encoding/json"
	"path"
	"runtime/debug"
	"strings"
	"time"

	s "github.com/saferwall/saferwall/pkg/strings"
	"github.com/saferwall/saferwall/pkg/utils"

	peparser "github.com/saferwall/pe"
	bs "github.com/saferwall/saferwall/pkg/bytestats"
	"github.com/saferwall/saferwall/pkg/crypto"
	"github.com/saferwall/saferwall/pkg/exiftool"
	"github.com/saferwall/saferwall/pkg/magic"
	"github.com/saferwall/saferwall/pkg/ml"
	"github.com/saferwall/saferwall/pkg/packer"
	"github.com/saferwall/saferwall/pkg/trid"
	log "github.com/sirupsen/logrus"
)

const (
	// defaultStringsMinLength represents the minimum length of extracted strings
	defaultStringsMinLength = 5
)

type stringStruct struct {
	Encoding string `json:"encoding"`
	Value    string `json:"value"`
}

// File represents a file object.
type File struct {
	MD5         string                 `json:"md5,omitempty"`
	SHA1        string                 `json:"sha1,omitempty"`
	SHA256      string                 `json:"sha256,omitempty"`
	SHA512      string                 `json:"sha512,omitempty"`
	SSDeep      string                 `json:"ssdeep,omitempty"`
	CRC32       string                 `json:"crc32,omitempty"`
	Magic       string                 `json:"magic,omitempty"`
	Size        int64                  `json:"size,omitempty"`
	Exif        map[string]string      `json:"exif,omitempty"`
	TriD        []string               `json:"trid,omitempty"`
	Tags        map[string]interface{} `json:"tags,omitempty"`
	Packer      []string               `json:"packer,omitempty"`
	LastScanned time.Time              `json:"last_scanned,omitempty"`
	Strings     []stringStruct         `json:"strings,omitempty"`
	MultiAV     map[string]interface{} `json:"multiav,omitempty"`
	Status      int                    `json:"status,omitempty"`
	PE          *peparser.File         `json:"pe,omitempty"`
	Histogram   []int                  `json:"histogram,omitempty"`
	ByteEntropy []int                  `json:"byte_entropy,omitempty"`
	Ml          map[string]interface{} `json:"ml,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

func determineType(magic string) string {

	var fileType string

	typeMap := map[string]string{
		"PE32":                    "pe",
		"MS-DOS":                  "msdos",
		"XML":                     "xml",
		"HTML":                    "html",
		"ELF":                     "elf",
		"PDF":                     "pdf",
		"Macromedia Flash":        "swf",
		"Zip archive data":        "zip",
		"Java archive data (JAR)": "jar",
		"PEG image data":          "jpeg",
		"PNG image data":          "png",
		"SVG Scalable Vector":     "svg",
	}

	for k, v := range typeMap {
		if strings.HasPrefix(magic, k) {
			fileType = v
			break
		}
	}

	return fileType
}

// Scan runs all scanners on the queued file.
func (f *File) Scan(sha256, filePath string, b []byte,
	ctxLogger *log.Entry, cfg *Config) error {

	var err error

	// Calculate the file size.
	f.Size = int64(len(b))

	// Calculates hashes.
	r := crypto.HashBytes(b)
	f.CRC32 = r.CRC32
	f.MD5 = r.MD5
	f.SHA1 = r.SHA1
	f.SHA256 = r.SHA256
	f.SHA512 = r.SHA512
	f.SSDeep = r.SSDeep

	// Extract file metadata using packer/trid/magic and exif.
	f.extractMetadata(sha256, filePath, ctxLogger, cfg)
	ctxLogger.Debug("metadata scan success")

	f.Strings = extractStrings(b, defaultStringsMinLength)
	ctxLogger.Debug("strings scan success")

	// Determine the file type.
	f.Type = determineType(f.Magic)

	// Parse the file.
	switch f.Type {
	case "pe":
		if f.PE, err = parsePE(filePath); err != nil {
			ctxLogger.Errorf("pe parser failed: %v", err)
		} else {
			ctxLogger.Debug("pe scan success")
		}

		// Extract Byte Histogram and byte entropy.
		f.Histogram, f.ByteEntropy = extractHistogramData(b)
		ctxLogger.Debug("bytestats scan success")
	}

	// Run the multi-av scanning.
	multiavScanRes := f.multiAvScan(filePath, cfg, ctxLogger)
	f.MultiAV = map[string]interface{}{}
	f.MultiAV["last_scan"] = multiavScanRes

	// Extract tags.
	f.extractTags()

	// Marshal results.
	var buff []byte
	if buff, err = json.Marshal(f); err != nil {
		ctxLogger.Errorf("failed to json marshal object: %v", err)
		return err
	}

	// Get ML classification results.
	f.Ml = map[string]interface{}{}
	switch f.Type {
	case "pe":
		if mlPredictionResults, err :=
			ml.PEClassPrediction(cfg.Ml.Address, buff); err != nil {
			ctxLogger.Errorf(
				"failed to get ml pe classifier prediction results: %v", err)
		} else {
			mlPredictionResults.SHA256 = ""
			f.Ml["pe"] = mlPredictionResults
		}
	}

	// Get ranked strings results.
	if mlStrRankerResults, err :=
		ml.RankStrings(cfg.Ml.Address, buff); err != nil {
		ctxLogger.Errorf(
			"failed to get ml string ranker prediction results: %v", err)
	} else {
		mlStrRankerResults.SHA256 = ""
		f.Ml["strings"] = mlStrRankerResults
	}

	// Finished scanning the file.
	f.LastScanned = time.Now().UTC()

	return nil
}

func parsePE(filePath string) (*peparser.File, error) {

	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	// Open the file and prepare it to be parsed.
	opts := peparser.Options{SectionEntropy: true}
	pe, err := peparser.New(filePath, &opts)
	if err != nil {
		return nil, err
	}
	defer pe.Close()

	// Parse the PE.
	err = pe.Parse()
	return pe, err
}

func scanFile(sha256 string, ctxLogger *log.Entry, h *MessageHandler) error {

	// Handle unexpected panics.
	defer func() {
		if r := recover(); r != nil {
			ctxLogger.Errorf("panic occured in file scan: %v", debug.Stack())
		}
	}()

	// Create a new file instance.
	f := File{SHA256: sha256}

	// Set the file status to `processing`.
	f.Status = processing
	err := h.updateMsgProgress(&f)
	if err != nil {
		ctxLogger.Errorf("failed to update message status: %v", err)
		return err
	}

	// Download the sample.
	filePath := path.Join(h.cfg.Consumer.DownloadDir, f.SHA256)
	b, err := h.downloadSample(filePath, &f)
	if err != nil {
		ctxLogger.Errorf("failed to download sample from s3: %v", err)
		return err
	}

	// Scan the file.
	err = f.Scan(sha256, filePath, b, ctxLogger, h.cfg)
	if err != nil {
		ctxLogger.Errorf("failed to scan the file: %v", err)
		return err
	}

	// Set the file status to `finished`.
	f.Status = finished
	err = h.updateMsgProgress(&f)
	if err != nil {
		ctxLogger.Errorf("failed to update message status: %v", err)
		return err
	}

	// Delete the file from the network share.
	if utils.Exists(filePath) {
		if err = utils.DeleteFile(filePath); err != nil {
			log.Errorf("failed to delete file path %s", filePath)
		}
	}

	return nil
}

func extractHistogramData(b []byte) (hist []int, byteentropy []int) {
	hist = bs.ByteHistogram(b)
	byteentropy = bs.ByteEntropyHistogram(b)
	return
}

func extractStrings(b []byte, minLength int) (strResults []stringStruct) {
	// Extract strings.
	n := minLength
	asciiStrings := s.GetASCIIStrings(b, n)
	wideStrings := s.GetUnicodeStrings(b, n)

	// Remove duplicates
	uniqueASCII := utils.UniqueSlice(asciiStrings)
	uniqueWide := utils.UniqueSlice(wideStrings)

	for _, str := range uniqueASCII {
		strResults = append(strResults, stringStruct{"ascii", str})
	}
	for _, str := range uniqueWide {
		strResults = append(strResults, stringStruct{"wide", str})
	}
	return
}

func (f *File) extractMetadata(sha256, filePath string, ctxLogger *log.Entry, cfg *Config) {
	var err error
	// Get exif metadata.
	if f.Exif, err = exiftool.Scan(filePath); err != nil {
		ctxLogger.Errorf("exiftool scan failed with: %v", err)
	} else {
		ctxLogger.Debug("exiftool scan success")
	}

	// Get TriD file identifier results.
	if f.TriD, err = trid.Scan(filePath); err != nil {
		ctxLogger.Errorf("trid scan failed with: %v", err)
	} else {
		ctxLogger.Debug("trid scan success")
	}

	// Get lib magic scan results.
	if f.Magic, err = magic.Scan(filePath); err != nil {
		ctxLogger.Errorf("magic scan failed with: %v", err)
	} else {
		ctxLogger.Debug("magic scan success")
	}

	// Retrieve packer/crypter scan results.
	if f.Packer, err = packer.Scan(filePath); err != nil {
		ctxLogger.Errorf("packer scan failed with: %v", err)
	} else {
		ctxLogger.Debug("packer scan success")
	}

}
