// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package gib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNGrams(t *testing.T) {

	testCases := []struct {
		input  string
		output []string
		n      int
	}{
		{
			"hello world",
			[]string{"hel", "ell", "llo", "lo ", "o w", " wo", "wor", "orl", "rld"},
			3,
		}, {
			"hello world",
			[]string{"hell", "ello", "llo ", "lo w", "o wo", " wor", "worl", "orld"},
			4,
		}, {
			"saferwall",
			[]string{"safer", "aferw", "ferwa", "erwal", "rwall"},
			5,
		}, {
			"",
			[]string{},
			5,
		},
	}

	for _, tt := range testCases {
		res := ngramsFromString(tt.input, tt.n)

		for i, s := range tt.output {
			if s != res[i] {
				t.Errorf("bad ngram expected %s got %s", s, res[i])
			}
		}
	}
}

func TestAllNGrams(t *testing.T) {
	expectedThreeNgrams, err := readLines("./testdata/all-three-ngrams.txt")
	if err != nil {
		t.Fatal("failed to read file with ", err)
	}
	allPossibleZeroNgrams := allNgrams(0)
	if len(allPossibleZeroNgrams) != 0 {
		t.Fatalf("expected 0-gram but got %d", len(allPossibleZeroNgrams))
	}
	allPossibleOneNgrams := allNgrams(1)
	for i, c := range allPossibleOneNgrams {
		if c != lowerCaseLetters[i] {
			t.Fatalf("wrong possible ngram expected %s got %s", c, lowerCaseLetters[i])
		}
	}
	allPossibleThreeNgrams := allNgrams(3)
	for i, c := range expectedThreeNgrams {
		if c != allPossibleThreeNgrams[i] {
			t.Fatalf("wrong possible ngram expected %s got %s", c, allPossibleThreeNgrams[i])
		}
	}
}

func TestIDF(t *testing.T) {
	totalStrings := 12030.
	stringFreq := 430.
	totalFreq := 840.
	maxFreq := 0.12

	expected := 4.80280496297434
	actual := ngramIDFValue(totalStrings, stringFreq, totalFreq, maxFreq)

	if expected != actual {
		t.Log(actual)
		t.Fatalf("bad idf value expected %f got %f", expected, actual)
	}
}

func TestNGramValues(t *testing.T) {
	corpus := []string{"hello world!", "saferwall is great", "is this gibberish"}
	n := 3
	ngrams := nGramValues(corpus, n, true)

	expectedIDF := 1.
	expectedFreq := 3.

	actualIDF := highestIDF(ngrams)
	actualFreq := highestFreq(ngrams)

	if actualIDF != expectedIDF || actualFreq != expectedFreq {
		t.Fatalf("bad idf values expected %f got %f | expected freq %f got %f", expectedIDF, actualIDF, expectedFreq, actualFreq)
	}

	expectedScore := 0.7443419746550483
	tfIDFScore := TFIDFScoreFunction(ngrams, 3, 25., 1.365, 1.159)

	if tfIDFScore("popopo") != expectedScore {
		t.Fatalf("expected score %f got %f", expectedScore, tfIDFScore("popopo"))
	}
}

func TestHeuristic(t *testing.T) {
	testCases := []struct {
		text     string
		expected bool
	}{
		{
			text:     "MMMMMMXVIII",
			expected: true,
		}, {
			text:     "pqxyzwww",
			expected: true,
		}, {
			text:     "aaaaaaaaaa",
			expected: true,
		}, {
			text:     "asdqwfbeqbfuilac",
			expected: true,
		},
	}

	for _, tt := range testCases {
		isNonsense := simpleNonSense(tt.text)
		if isNonsense != tt.expected {
			t.Fatalf("expected nonsense(%s) to be %t got %t",
				tt.text, tt.expected, isNonsense)
		}
	}
}

func TestSanitize(t *testing.T) {

	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "!hello.;",
			output: "hello",
		}, {
			input:  "@HEY.$",
			output: "hey",
		}, {
			input:  "sAfErWaLl<=>?@[\\]^_`{|}~ ",
			output: "saferwall",
		}, {
			input:  "Bartók's",
			output: "bartks",
		},
	}

	for _, tt := range testCases {
		san := sanitize(tt.input)
		if san != tt.output {
			t.Fatalf("sanitizing failed expected %s got %s", tt.output, san)
		}
	}
}

func TestDataset(t *testing.T) {
	_, err := loadDataset("./data/ngram.json")
	if err != nil {
		t.Fatal("failed to load dataset with error : ", err)
	}
}

func TestScorer(t *testing.T) {
	testCases := []string{
		"lakdfqtajaklj",
		"AaBbCcDdEeFGgHhIiJjKkLlMmNnOoPpQqRrSsTtU",
		"AcoGQMJyIapivScpnfuXUDMtgTtvAACYdAyABnSLpoABhzZWAAVvAYAAAnqUAFTPo",
		"BCDEFGHIJKLMNOPQRSTUVWXYZ",
		"CDjjiJJTbvFWaSdEtUygGMoGl",
		"CQgHAwIDFQIDAxYCAQIeAQIXgAAKCRC",
		"CgKDQpPUkdBTklaSUHIENPTUJVFRFRQKLStLStLStLStLStLStLStLStLStLSt",
		"aoaoesuouooeueooeuoaeuoeou",
		"iuewrofahgalkfgaufpiupqrjf",
		"ieeoienkjadfakj",
		"lalalaalkjuogaajfajlfal",
	}
	_, err := NewScorer(nil)
	if err != nil {
		t.Fatal("failed to create new score function with error :", err)
	}
	nonsense, err := NewScorer(&Options{Dataset: Dataset})
	if err != nil {
		t.Fatal("failed to create new score function with error :", err)
	}
	for _, tt := range testCases {
		isNonsense, _ := nonsense(tt)
		if isNonsense != true {
			t.Fatalf("expected %t but got %t", true, isNonsense)
		}
	}
}

func TestScoreFunctionOnLabeledData(t *testing.T) {

	type TestCase struct {
		input     string
		knownReal bool
		expected  bool
	}
	testCases := make([]TestCase, 0)
	// populate test cases
	labeledCases, err := readLines("./testdata/labeledCases.csv")
	if err != nil {
		t.Fatal("failed to read test data with error :", err)
	}

	for _, c := range labeledCases {
		var input string
		var knownReal bool
		var expected bool
		testCase := strings.Split(c, ",")

		input = testCase[1]
		res := testCase[0]

		if res == "y" {
			knownReal = true
			expected = false
		} else if res == "n" {
			knownReal = false
			expected = true
		}

		testCases = append(testCases, TestCase{
			input:     input,
			knownReal: knownReal,
			expected:  expected,
		})
	}

	isGibberish, err := NewScorer(nil)
	if err != nil {
		t.Fatal("failed to create new score function with error :", err)
	}
	var trueNegatives int
	var truePositives int

	var falseNegatives []string
	var falsePositives []string

	for _, tt := range testCases {
		actual, err := isGibberish(tt.input)
		if err != nil {
			t.Fatal("failed to score string with error :", err)
		}
		knownReal := tt.knownReal
		labeledAsGibberish := actual

		if knownReal {
			if labeledAsGibberish {
				falsePositives = append(falsePositives, tt.input)
			} else {
				trueNegatives++
			}
		} else {
			if labeledAsGibberish {
				truePositives++
			} else {
				falseNegatives = append(falseNegatives, tt.input)
			}
		}
	}

	fpCount := len(falsePositives)
	fnCount := len(falseNegatives)

	precision := Precision(truePositives, fpCount) * 100
	recall := Recall(truePositives, fnCount) * 100
	accuracy := Accuracy(truePositives, fpCount, trueNegatives, fnCount) * 100
	testResults := fmt.Sprintf("Test Results : \n Accuracy %f %% Precision : %f %% \t Recall %f %% \n True Positives : %d \t True Negatives : %d \n False Positives : %d \t False Negatives : %d\n ", accuracy, precision, recall, truePositives, trueNegatives, fpCount, fnCount)
	t.Log(testResults)
}

func TestScoreFunctionOnRealData(t *testing.T) {

	// all the test cases are real string read from /usr/share/dict/words
	testCases, err := readLines("/usr/share/dict/words")
	if err != nil {
		t.Fatal("failed to read test cases with error :", err)
	}

	isGibberish, err := NewScorer(nil)
	if err != nil {
		t.Fatal("failed to create new score function with error ", err)
	}
	fp := 0
	for _, tt := range testCases {
		if len(sanitize(tt)) <= 6 {
			continue
		}
		gibberish, err := isGibberish(tt)
		if err != nil {
			t.Fatal("failed to score string with error : ", err)
		}
		if gibberish != false {
			fp++
		}
	}
	t.Logf("Found %d False Positives (Real String Marked as Gibberish", fp)

}

func TestScoreFunctionOnLudiso(t *testing.T) {

	// this test runs on ludiso dataset similar to the above test
	// we expect all test cases to be real strings
	ludiso, err := readLines("./testdata/ludiso.txt")
	if err != nil {
		t.Fatal("could not read data file failed with error :", err)
	}

	isGibberish, err := NewScorer(nil)
	if err != nil {
		t.Fatal("failed to create new score function with error :", err)
	}
	var fpRate int
	var tpRate int
	var fnRate int
	var tnRate int

	for _, s := range ludiso {

		isRandom, err := isGibberish(s)
		if err != nil {
			t.Fatalf("failed to detect %s with error : %s", s, err.Error())
		}
		// the given string was real (negative class) but was tagged as positive class
		tnRate += boolToInt(!isRandom)
		fpRate += boolToInt(isRandom)

	}

	precision := Precision(tpRate, fpRate) * 100
	recall := Recall(tpRate, fnRate) * 100
	accuracy := Accuracy(tpRate, fpRate, tnRate, fnRate) * 100
	testResults := fmt.Sprintf("Test Results : \n Accuracy %f %% Precision : %f %% \t Recall %f %% \n True Positives : %d \t True Negatives : %d \n False Positives : %d \t False Negatives : %d\n ", accuracy, precision, recall, tpRate, tnRate, fpRate, fnRate)
	t.Log(testResults)
}

func TestOnMacroDocs(t *testing.T) {

	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "H67oooeewxpd8ll",
			expected: true,
		}, {
			input:    "IGwkqQGAL(lAwPHFBmE + lAwPHFBmE)",
			expected: true,
		}, {
			input:    "Y0sd1ec1f2mgj4",
			expected: true,
		}, {
			input:    "RtlMoveMemory",
			expected: false,
		}, {
			input:    "OEACCDH",
			expected: true,
		}, {
			input:    "QcNNKBE",
			expected: true,
		}, {
			input:    "Ogilhvhap",
			expected: true,
		},
	}

	isGibberish, err := NewScorer(nil)

	if err != nil {
		t.Fatal("failed to create new score function with error :", err)
	}
	for _, tt := range testCases {
		israndom, err := isGibberish(tt.input)
		if err != nil {
			t.Fatalf("failed to test on string %s with error %v",
				tt.input, err)
		}
		if israndom != tt.expected {
			t.Logf("ambigious on test case %s expected %t got %t",
				tt.input, tt.expected, israndom)
		}
	}

}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// boolToInt
func boolToInt(b bool) int {

	if b {
		return 1
	}

	return 0

}
