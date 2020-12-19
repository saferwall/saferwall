# Byte and Entropy Statistics for Binary Files

This module implements byte and entropy extract for binary files
it is loosely based on [Deep Neural Network Based Malware Detection Using Two Dimensional Binary Program Features](https://arxiv.org/pdf/1508.03096.pdf).

This module outputs byte historgram statistics and 2D byte & entropy histogram over byte
values in a binary file.

## Usage

```golang
package main

func main() {
    bytez,err := ioutil.ReadFile("bin/sh")
    if err != nil {
        panic(err)
    }
    // Compute a histogram of byte distributions
    byteHistogram := ByteHistogram(bytez)
    // Compute a byte-entropy histogram
    byteEntropyHist := ByteEntropyHistogram(bytez)
}

```