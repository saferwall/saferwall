package bytestats

const (
	minLength = 256
	window    = 2048
	step      = 1024
)

// binCount counts the number of occurences of each byte value in a buffer.
func binCount(buf []byte, minlength int) []int {

	count := make([]int, minLength)

	for _, b := range buf {
		count[int(b)]++
	}

	return count
}

// entropyBinCount calculates the coarse entropy histogram of byte values.
func entropyBinCount(buf []byte, block int) {

}
