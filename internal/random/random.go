// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// package random generate a random file name. The random
// generator use the english disctionary words instead of
// gibberish strings as the malware could detect high
// gibberish strings.
package random

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Ramdomizer is an abstract interface for generating random
// strings.
type Ramdomizer interface {
	Random() string
}

type Service struct {
	words []string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// New creates a new service.
func New(wordsPath string) (Service, error) {

	words, err := readAvailableDictionary(wordsPath)
	if err != nil {
		return Service{}, err
	}
	return Service{words}, nil
}

// Random picks a random strings from the list of english words.
func (s Service) Random() string {
	return s.words[rand.Int()%len(s.words)]
}

func readAvailableDictionary(wordsPath string) ([]string, error) {

	file, err := os.Open(wordsPath)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(bytes), "\n"), nil
}
