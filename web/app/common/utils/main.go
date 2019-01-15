// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// IsStringInSlice check if a string exist in a list of strings
func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

// GetQueryParamsFields retrieve `fields`` so we can filter them in GET/
func GetQueryParamsFields(c echo.Context) []string {
	var filters []string
	fields := c.QueryParam("fields")
	if fields != "" {
		filters = strings.Split(fields, ",")
	}

	return filters
}

// IsFilterAllowed check if we are allowed to filter GET with fields
func IsFilterAllowed(allowed []string, filters []string) bool {
	for _, filter := range filters {
		if !IsStringInSlice(filter, allowed) {
			return false
		}
	}
	return true
}

// JSONTime alias for timne.Time
type JSONTime time.Time

// MarshalJSON will marshall date
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}
