// Copyright 2021 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package consumer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	mlStaticPeEndpoint      = "/api/static/pe"
	mlStaticStringsEndpoint = "/api/static/strings"
)

type MlClassifierPrediction struct {
	Class       string  `json:"predicted_class,omitempty"`
	Probability float64 `json:"predicted_probability,omitempty"`
	Score       string  `json:"predicted_score,omitempty"`
	Sha256      string  `json:"sha256,omitempty"`
}

type MlStringRanker struct {
	Strings []string `json:"strings,omitempty"`
	Sha256  string   `json:"sha256,omitempty"`
}

func mlPEClassPredResult(buff []byte) (MlClassifierPrediction, error) {

	client := &http.Client{}
	res := MlClassifierPrediction{}
	client.Timeout = time.Second * 15
	url := mlEndpoint + mlStaticPeEndpoint

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(buff))
	if err != nil {
		contextLogger.Errorf("http.NewRequest() failed with: %v", err)
		return res, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		contextLogger.Errorf("client.Do() failed with: %v", err)
		return res, err
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		contextLogger.Errorf("ioutil.ReadAll() failed with: %v", err)
		return res, err
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		contextLogger.Errorf("json.Unmarshal() failed with: %v", err)
		return res, err
	}

	contextLogger.Infof("ML PE Classifier success: status code: %d, resp: %s",
		resp.StatusCode, string(d))

	return res, nil
}

func mlStringRanker(buff []byte) (MlStringRanker, error) {

	client := &http.Client{}
	res := MlStringRanker{}
	client.Timeout = time.Second * 15
	url := mlEndpoint + mlStaticStringsEndpoint

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(buff))
	if err != nil {
		contextLogger.Errorf("http.NewRequest() failed with: %v", err)
		return res, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		contextLogger.Errorf("client.Do() failed with: %v", err)
		return res, err
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		contextLogger.Errorf("ioutil.ReadAll() failed with: %v", err)
		return res, err
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		contextLogger.Errorf("json.Unmarshal() failed with: %v", err)
		return res, err
	}

	contextLogger.Infof("ML String Ranker success: status code: %d, resp: %s",
		resp.StatusCode, string(d))

	return res, nil
}
