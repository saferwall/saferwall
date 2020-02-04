// Copyright 2019 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
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

// RemoveStringFromSlice removes a string item from a list of strings.
func RemoveStringFromSlice(s []string, r string) []string {
    for i, v := range s {
        if v == r {
            return append(s[:i], s[i+1:]...)
        }
    }
    return s
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

// GetStructFields retrieve json struct fields names
func GetStructFields(i interface{}) []string {

	val := reflect.ValueOf(i)
	var temp string

	var listFields []string
	for i := 0; i < val.Type().NumField(); i++ {
		temp = val.Type().Field(i).Tag.Get("json")
		temp = strings.Replace(temp, ",omitempty", "", -1)
		listFields = append(listFields, temp)
	}

	return listFields
}

func fieldSet(fields ...string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, s := range fields {
		set[s] = true
	}
	return set
}

// SelectFields execlude sensitive fields
func SelectFields(i interface{}, fields ...string) map[string]interface{} {
	fs := fieldSet(fields...)
	rt, rv := reflect.TypeOf(i), reflect.ValueOf(i)
	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonKey := field.Tag.Get("json")
		if fs[jsonKey] {
			out[jsonKey] = rv.Field(i).Interface()
		}
	}
	return out
}

// JSONTime alias for timne.Time
type JSONTime time.Time

// MarshalJSON will marshall date
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}