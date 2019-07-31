package trid

import (
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	// Command to invoke TriD scanner
	Command = "trid"
)

// Scan a file using TRiD Scanner
// This will execute trid command line tool and read the stdout
func Scan(FilePath string) ([]string, error) {

	args := []string{FilePath}
	output, err := utils.ExecCommand(Command, args...)
	if err != nil {
		return []string{}, err
	}
	return parseOutput(output), nil

}

// parseOutput parse TriD stdout, returns an array of strings
func parseOutput(tridout string) []string {

	keepLines := []string{}
	lines := strings.Split(tridout, "\n")
	if utils.StringInSlice("Error: found no file(s) to analyze!", lines) {
		return nil
	}
	lines = lines[6:]

	for _, line := range lines {
		if len(strings.TrimSpace(line)) != 0 {
			keepLines = append(keepLines, strings.TrimSpace(line))
		}
	}

	return keepLines
}
