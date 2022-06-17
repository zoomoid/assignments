package util

import "fmt"

// addLeadingZero prepends numbers smaller than 10 with a leading zero
func AddLeadingZero(assignment uint32) string {
	if assignment < 10 {
		return fmt.Sprintf("0%d", assignment)
	}
	return fmt.Sprintf("%d", assignment)
}
