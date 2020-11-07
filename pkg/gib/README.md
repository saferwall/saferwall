# gib : A String Gibberish Detector in Golang

*gib* is a package to detect gibberish strings in Golang. This utility is useful when analyzing textual artefacts in malwares. To cite how we're using this package:
  - Macro documents/scripts usually obfuscate their code by randomizing variables and function names, applying `gib` over extracted tokens is a simple heuritic/ ML feature to raise when such a situation happens.
  - Be able to tell when malwares drop random files names, or uses random PE section names. Well you get the idea :)

## Installation

```gib``` is built as a go module so you can quickly get started by running:

```sh
go get github.com/saferwall/saferwall/pkg/gib
```

## Usage

A quick example for using ```sh gib ``` with the same dataset as published by [nostril](https://github.com/casics/nostril)

```go

import (
  "fmt"
  "github.com/saferwall/saferwall/pkg/gib"  
)


func main() {

  // A couple test cases.
  randomString := "asdqwfbeqbfuilac"
  nonRandomString := "CreateNewUser"
  
  // Create a new gibberish detector
  isGibberish := gib.NewScorer(nil)

  // Will return `True`.
  fmt.Println(isGibberish(randomString))

  // Will return `False`.
  fmt.Println(isGibberish(nonRandomString))
```
## Limitations

```gib``` is strongly based on ```nostril``` and re-uses its training data for filtering strings.
Thus ```gib``` limitations are similar, it only supports American English words and will generate false positives
due to the problem domain.

The idea behind both is to fine tune certain scoring values such as length and repition based on its training corpus.
A custom IDF value is computed for each given string, the default n-gram value is set to 4 similar to the dataset provided.

To make an analogy to how TF-IDF is used in document classification, nonsense strings are those that use unusual n-grams and thus score highly, while real/meaningful strings are those that use more common n-grams and thus score lower. The Nostril package includes a precomputed table of n-gram weights derived by training the system on a large set of strings constructed from concatenated American English words, real text corpora, and other inputs. Parameter values were optimized using the evolutionary algorithm NSGA-II.


## Notes about testing

The score function acts as a classifier returning true if the given string is gibberish false otherwise.
Therefore the rate of false positives points to the number of non-gibberish strings marked as gibberish
a false negative in this case is when a gibberish string is marked as non-gibberish.

The score function is evaluated using accuracy, precision and recall.

## References

- [Nostril : A Nonsense String Evaluator Written in Python](https://www.theoj.org/joss-papers/joss.00596/10.21105.joss.00596.pdf)
