// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package gib metrics.go implements accuracy metrics for rating detection on
// several test cases.
package gib

// Labels : positive class is gibberish and negative class is correct string
// True Positive : A gibberish string labeled by the score function as gibberish
// False Positive : A correct (non-gibberish) string labeled by the score function as gibberish
// True Negative : A correct string labeled by the score function as non-gibberish
// False Negative : A gibberish string labeled by the score function as non-gibberish

// Accuracy is the fraction of predictions gib made correctly.
func Accuracy(truePositiveCount, falsePositiveCount, trueNegativeCount,
	falseNegativeCount int) float64 {

	// just cast to float64
	tpCount := float64(truePositiveCount)
	fpCount := float64(falsePositiveCount)
	tnCount := float64(trueNegativeCount)
	fnCount := float64(falseNegativeCount)

	return (tpCount + tnCount) / (tpCount + fpCount + tnCount + fnCount)
}

// Precision is a metric to answer the question:
// "What proportion of positive identifications was actually correct".
func Precision(truePositiveCount, falsePositiveCount int) float64 {

	tpCount := float64(truePositiveCount)
	fpCount := float64(falsePositiveCount)

	// if a model has zero false positives it's precision is 1.
	return tpCount / (tpCount + fpCount)
}

// Recall is a metric to answer the question:
// "What proportion of actual positives was identified correctly?"
func Recall(truePositiveCount, falseNegativeCount int) float64 {

	tpCount := float64(truePositiveCount)
	fnCount := float64(falseNegativeCount)

	return tpCount / (tpCount + fnCount)
}
