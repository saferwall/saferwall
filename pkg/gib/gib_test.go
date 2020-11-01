package gib

import (
	"bufio"
	"os"
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
		res := ngrams(tt.input, tt.n)

		for i, s := range tt.output {
			if s != res[i] {
				t.Errorf("bad ngram exepceted %s got %s", s, res[i])
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
			t.Fatalf("expected nonsense(%s) to be %t got %t", tt.text, tt.expected, isNonsense)
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
	_, err := loadDataset("./dataset/ngram.json")
	if err != nil {
		t.Fatal("failed to load dataset with error : ", err)
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
