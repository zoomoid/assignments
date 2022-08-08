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
