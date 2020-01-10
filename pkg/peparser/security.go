package pe

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"encoding/hex"
	"go.mozilla.org/pkcs7"
)

// The options for the WIN_CERTIFICATE Revision member include
// (but are not limited to) the following.
const (
	// WinCertRevision1_0 represents the WIN_CERT_REVISION_1_0 Version 1,
	// legacy version of the Win_Certificate structure.
	// It is supported only for purposes of verifying legacy Authenticode signatures
	WinCertRevision1_0 = 0x0100

	// WinCertRevision2_0 represents the WIN_CERT_REVISION_2_0. Version 2
	// is the current version of the Win_Certificate structure.
	WinCertRevision2_0 = 0x0200
)

// The options for the WIN_CERTIFICATE CertificateType member include
// (but are not limited to) the items in the following table. Note that some
// values are not currently supported.
const (
	WinCertTypeX509           = 0x0001 // Certificate contains an X.509 Certificate (Not Supported)
	WinCertTypePKCSSignedData = 0x0002 //Certificate contains a PKCS#7 SignedData structure
	WinCertTypeReserved1      = 0x0003 // Reserved
	WinCertTypeTsStackSigned  = 0x0004 // Terminal Server Protocol Stack Certificate signing (Not Supported)
)

// Certificate directory.
type Certificate struct {
	Header WinCertificate
	Content   *pkcs7.PKCS7
}

// WinCertificate encapsulates a signature used in verifying executable files.
type WinCertificate struct {
	Length          uint32 // Specifies the length, in bytes, of the signature.
	Revision        uint16 // Specifies the certificate revision.
	CertificateType uint16 // Specifies the type of certificate.
}


func (pe *File) parseSecurityDirectory(rva, size uint32) (Certificate, error) {

	certHeader := WinCertificate{}
	certSize := uint32(binary.Size(certHeader))
	var pkcs *pkcs7.PKCS7

	// The virtual address value from the Certificate Table entry in the
	// Optional Header Data Directory is a file offset to the first attribute
	// certificate entry.
	fileOffset := rva

	for {
		// Boundary check
		if fileOffset+certSize > pe.size {
			log.Print("Overflow")
			return Certificate{}, errors.New("Boundary checks failed in Security Data Dir")
		}

		buf := bytes.NewReader(pe.data[fileOffset : fileOffset+certSize])
		err := binary.Read(buf, binary.LittleEndian, &certHeader)
		if err != nil {
			return Certificate{}, err
		}

		certContent := pe.data[fileOffset+certSize : fileOffset+certHeader.Length]
		pkcs, err = pkcs7.Parse(certContent)
		if err != nil {
			return Certificate{Header: certHeader}, err
		}

		// Verify the signature
		pkcs.Verify()

		log.Printf("%s", hex.Dump(pkcs.Content))

		// Subsequent entries are accessed by advancing that entry's dwLength bytes,
		// rounded up to an 8-byte multiple, from the start of the current attribute
		// certificate entry.
		nextOffset := certHeader.Length + fileOffset
		nextOffset = ((nextOffset + 8 - 1) / 8) * 8

		// Check if we walked the entire table.
		if nextOffset == fileOffset+size {
			break
		}

		fileOffset = nextOffset
	}

	return Certificate{Header: certHeader, Content: pkcs}, nil
}
