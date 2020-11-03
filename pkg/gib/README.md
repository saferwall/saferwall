# gib : A String Gibberish Detector in Golang

*gib* is a package to detect gibberish strings in Golang.
This utility is useful when analyzing textual artefacts in compiled software such
as malformed section names in the case of PE Files or malicious Macros.

## Install

```gib``` is built as a go module so you can quickly get started by running

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

  // a couple test cases
  randomString := "asdqwfbeqbfuilac"
  
  nonRandomString := "CreateNewUser"
  
  // path to the provided default dataset
  pathToDataset := "/home/user/go/github.com/saferwall/saferwall/pkg/gib/data/ngram.json"
  
  // load dataset as an ngram score table
  defaultDataset := gib.LoadDataset(pathToDataset)

  // create a new gibberish detector
  isGibberish := gib.NewScorer(defaultDataset)

  fmt.Println(isGibberish(randomString))

  fmt.Println(isGibberish(nonRandomString))

```

## References

- [Nostril : A Nonsense String Evaluator Written in Python](https://www.theoj.org/joss-papers/joss.00596/10.21105.joss.00596.pdf)
