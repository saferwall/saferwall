// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package ml

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	staticPeEndpoint      = "/api/static/pe"
	staticStringsEndpoint = "/api/static/strings"
	clientTimeout         = time.Second * 15
)

// ClassifierPrediction represents the classifier prediction result.
type ClassifierPrediction struct {
	Class       string  `json:"predicted_class,omitempty"`
	Probability float64 `json:"predicted_probability,omitempty"`
	Score       string  `json:"predicted_score,omitempty"`
	SHA256      string  `json:"sha256,omitempty"`
}

// StringsRanker represents the string ranker results.
type StringsRanker struct {
	Strings []string `json:"strings,omitempty"`
	SHA256  string   `json:"sha256,omitempty"`
}

// PEClassPrediction returns the ML PE classifier results.
func PEClassPrediction(server string, buff []byte) (
	ClassifierPrediction, error) {

	res := ClassifierPrediction{}
	url := server + staticPeEndpoint

	client := &http.Client{}
	client.Timeout = clientTimeout

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buff))
	if err != nil {
		return res, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(d, &res)
	return res, err
}

// RankStrings applies the String ranker model to a list of strings.
func RankStrings(server string, buff []byte) (StringsRanker, error) {

	client := &http.Client{}
	client.Timeout = clientTimeout

	res := StringsRanker{}
	url := server + staticStringsEndpoint

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buff))
	if err != nil {
		return res, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(d, &res)
	return res, err
}
