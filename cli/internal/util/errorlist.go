package util

import "errors"

// ErrorList is a wrapper for multiple errors
// Helpful with stuff like combined validator functions
// where more than one error may occur
type ErrorList interface {
	error
	Errors() []error
	Is(error) bool
}

// NewErrorList constructs a new ErrorList wrapper type for multiple errors
//
// If an empty list is passed in, retuns nil
func NewErrorList(errlist []error) ErrorList {
	if len(errlist) == 0 {
		return nil
	}
	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errorList(errlist)
}

// errorList is the internal, private error list. Not exposing it
// to prevent empty error lists
type errorList []error

// Error stringifies all errors in the list and returns the string
// Implements the error interface
func (el errorList) Error() string {
	if len(el) == 0 {
		return ""
	}
	if len(el) == 1 {
		return el[0].Error()
	}
	seenerrs := NewSet()
	res := ""
	el.visit(func(err error) bool {
		msg := err.Error()
		if seenerrs.Has(msg) {
			return false
		}
		seenerrs.Insert(msg)
		if len(seenerrs) > 1 {
			res += ", "
		}
		res += msg
		return false
	})
	if len(seenerrs) == 1 {
		return res
	}
	return "[" + res + "]"
}

// Errors returns the list of all errors in the list
func (el errorList) Errors() []error {
	return []error(el)
}

// Is implements the error comparison where it matches if at least one error in the list matches the compared one
func (el errorList) Is(target error) bool {
	return el.visit(func(err error) bool {
		return errors.Is(err, target)
	})
}

// visit implements the visitor pattern on the error list
func (el errorList) visit(f func(err error) bool) bool {
	for _, err := range el {
		if match := f(err); match {
			return match
		}
	}
	return false
}
