/*
Copyright 2022 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	// Regex pattern for matching filenames against
	AssignmentPattern          string = "^assignment-([0-9][0-9]+).pdf$"
	AssignmentDirectoryPattern string = "^assignment-([0-9][0-9]+)$"
)

// addLeadingZero prepends numbers smaller than 10 with a leading zero
func AddLeadingZero(assignment uint32) string {
	if assignment < 10 {
		return fmt.Sprintf("0%d", assignment)
	}
	return fmt.Sprintf("%d", assignment)
}

// assignmentNumber extracts the assignment number (with a leading zero already,
// thus string-typed) from the target directory of the form "assignment-*"
func AssignmentNumberFromRegex(pattern string, input string) (string, error) {
	r, _ := regexp.Compile(pattern)
	if !r.MatchString(input) {
		return "", fmt.Errorf("target directory %s does not match the common pattern %s", input, AssignmentPattern)
	}
	number := r.ReplaceAllString(input, "$1")
	i, err := strconv.Atoi(number)
	if err != nil {
		return "", err
	}
	return AddLeadingZero(uint32(i)), nil
}
