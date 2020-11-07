# gib : A String Gibberish Detector in Golang

*gib* is a package to detect gibberish strings in Golang. This utility is useful when analyzing textual artefacts in malwares. To cite how we're using this package:
  - Macro documents/scripts usually obfuscate their code by randomizing variables and function names, applying `gib` over extracted tokens is a simnple way to tell when  is a simple heuritic/ ML feature to raise.
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
  
  // path to the provided default dataset
  datasetPath := "/home/user/go/github.com/saferwall/saferwall/pkg/gib/data/ngram.json"
  
  // load dataset as an ngram score table
  defaultDataset := gib.LoadDataset(pathToDataset)

  // create a new gibberish detector
  isGibberish := gib.NewScorer(defaultDataset)

  // Will return `True`.
  fmt.Println(isGibberish(randomString))

  // Will return `False`.
  fmt.Println(isGibberish(nonRandomString))
```

## Notes about testing

The score function acts as a classifier returning true if the given string is gibberish false otherwise.
Therefore the rate of false positives points to the number of non-gibberish strings marked as gibberish
a false negative in this case is when a gibberish string is marked as non-gibberish.

The score function is evaluated using accuracy, precision and recall.

## References

- [Nostril : A Nonsense String Evaluator Written in Python](https://www.theoj.org/joss-papers/joss.00596/10.21105.joss.00596.pdf)
