package ssdeep

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

var (
	// ErrEmptyHash is returned when no hash string is provided for scoring.
	ErrEmptyHash = errors.New("empty string")

	// ErrInvalidFormat is returned when a hash string is malformed.
	ErrInvalidFormat = errors.New("invalid ssdeep format")
)

// Distance computes the match score between two fuzzy hash signatures.
// Returns a value from zero to 100 indicating the match score of the two signatures.
// A match score of zero indicates the signatures did not match.
// Returns an error when one of the inputs are not valid signatures.
func Distance(hash1, hash2 string) (score int, err error) {
	hash1BlockSize, hash1String1, hash1String2, err := splitSsdeep(hash1)
	if err != nil {
		return
	}
	hash2BlockSize, hash2String1, hash2String2, err := splitSsdeep(hash2)
	if err != nil {
		return
	}

	if hash1BlockSize == hash2BlockSize && hash1String1 == hash2String1 {
		return 100, nil
	}

	// We can only compare equal or *2 block sizes
	if hash1BlockSize != hash2BlockSize && hash1BlockSize != hash2BlockSize*2 && hash2BlockSize != hash1BlockSize*2 {
		return
	}

	if hash1BlockSize == hash2BlockSize {
		d1 := scoreDistance(hash1String1, hash2String1, hash1BlockSize)
		d2 := scoreDistance(hash1String2, hash2String2, hash1BlockSize*2)
		score = int(math.Max(float64(d1), float64(d2)))
	} else if hash1BlockSize == hash2BlockSize*2 {
		score = scoreDistance(hash1String1, hash2String2, hash1BlockSize)
	} else {
		score = scoreDistance(hash1String2, hash2String1, hash2BlockSize)
	}
	return
}

func splitSsdeep(hash string) (blockSize int, hashString1, hashString2 string, err error) {
	if hash == "" {
		err = ErrEmptyHash
		return
	}

	parts := strings.Split(hash, ":")
	if len(parts) != 3 {
		err = ErrInvalidFormat
		return
	}

	blockSize, err = strconv.Atoi(parts[0])
	if err != nil {
		err = ErrInvalidFormat
		return
	}

	hashString1 = parts[1]
	hashString2 = parts[2]
	return
}

func scoreDistance(h1, h2 string, blockSize int) int {
	d := distance(h1, h2)
	d = (d * spamSumLength) / (len(h1) + len(h2))
	d = (100 * d) / spamSumLength
	d = 100 - d
	/* TODO: Figure out this black magic...
	matchSize := float64(blockSize) / float64(blockMin) * math.Min(float64(len(h1)), float64(len(h2)))
	if d > int(matchSize) {
		d = int(matchSize)
	}
	*/
	return d
}
