package bytestats

import "math"

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

// rollingWindow returns a rolling window of values accross 1D slice.
func rollingWindow(buf []byte, window int) [][]byte {
	r := make([][]byte, 0, len(buf))
	for i := window - 1; i < len(buf); i++ {
		r = append(r, buf[(i-window+1):i+1])
	}
	return r
}

// shiftBytes applies right shift operator to all buffer values.
func shiftBytes(buf []byte, s int) []byte {
	b := make([]byte, s)
	for _, v := range buf {
		b = append(b, v>>s)
	}
	return b
}

// apply applies fun to all entries in the list
func apply(f func(b float32) float32, list []float32) []float32 {

	r := make([]float32, len(list))

	for idx, v := range list {
		r[idx] = f(v)
	}

	return r
}

// asFloat32 casts a slice of integers to float32.
func asFloat32(ar []int) []float32 {
	newar := make([]float32, len(ar))
	for idx, v := range ar {
		newar[idx] = float32(v)
	}
	return newar
}

// nonZeroEntries returns a slice of non-zero indexes in a given slice.
func nonZeroEntries(a []int) []int {
	b := make([]int, 0, len(a))

	for idx, v := range a {
		if v != 0 {
			b = append(b, idx)
		}
	}
	return b
}

// sum computes sum of a slice elements.
func sum(s []float32) float32 {
	var r float32

	for _, v := range s {
		r += v
	}
	return r
}

// entropyBinCount calculates the coarse entropy histogram of byte values.
func entropyBinCount(block []byte, window int) (int, []int) {

	var H float64
	shiftedBlock := shiftBytes(block, 4)
	c := binCount(shiftedBlock, 16)
	p := apply(func(b float32) float32 {
		return b / float32(window)
	}, asFloat32(c))
	wh := nonZeroEntries(c)
	var tmp float64
	for _, entry := range wh {
		a := -float64(p[entry]) * math.Log2(float64(p[entry]))
		tmp += a
	}
	H = tmp * 2
	Hbin := int(H * 2)
	if Hbin == 16. {
		Hbin = 15.
	}
	return Hbin, c
}

// asFeatureVector encodes the histogram entropy values to a singular vector
func asFeatureVector(buf []byte, step, window int) []float32 {
	output := make([][]float32, 16)

	for idx := range output {
		output[idx] = make([]float32, 16)
	}

	if len(buf) < window {
		Hbin, c := entropyBinCount(buf, window)
		output[Hbin] = vectorizeSum(output[Hbin], asFloat32(c))
	} else {
		blocks := rollingWindow(buf, window)
		for _, block := range blocks {
			Hbin, c := entropyBinCount(block, window)
			output[Hbin] = vectorizeSum(output[Hbin], asFloat32(c))
		}
	}

	return flatten2D(output)
}

// flatten2D flattens a 2D slice to a 1D slice
func flatten2D(a [][]float32) []float32 {
	r := make([]float32, 0)

	for _, row := range a {
		for _, entry := range row {
			r = append(r, entry)
		}
	}
	return r
}

// vectorizeSum computes the sum of two float32 vectors
func vectorizeSum(a []float32, b []float32) []float32 {
	if len(a) != len(b) {
		panic("Operands could not be broadcasted together with incompatible dimensions")
	}
	c := make([]float32, len(b))

	for idx := range b {
		c[idx] = a[idx] + b[idx]
	}
	return c
}
