// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GetFileSize returns file size in bytes
func GetFileSize(FilePath string) int64 {
	f, err := os.Stat(FilePath)
	if err != nil {
		log.Fatal(err)
	}
	size := f.Size()
	return size
}

// GetCurrentTime as Time object in UTC
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// SliceContainsString returns if slice contains substring
func SliceContainsString(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(b, a) {
			return true
		}
	}
	return false
}

// StringInSlice returns whether or not a string exists in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ExecCommand runs cmd on file
func ExecCommand(name string, args ...string) (string, error) {

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, name, args...)

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
		return string(out), err
	}

	// If there's no context error, we know the command completed (or errored).
	return string(out), err
}
