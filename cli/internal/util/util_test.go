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

import "testing"

func TestAddLeadingZero(t *testing.T) {
	t.Run("assignment=1", func(t *testing.T) {
		an := AddLeadingZero(1)
		expected := "01"
		if an != expected {
			t.Errorf("expected %s, found %s", expected, an)
		}
	})
	t.Run("assignment=0", func(t *testing.T) {
		an := AddLeadingZero(0)
		expected := "00"
		if an != expected {
			t.Errorf("expected %s, found %s", expected, an)
		}
	})
	t.Run("assignment=9", func(t *testing.T) {
		an := AddLeadingZero(9)
		expected := "09"
		if an != expected {
			t.Errorf("expected %s, found %s", expected, an)
		}
	})
	t.Run("assignment=10", func(t *testing.T) {
		an := AddLeadingZero(10)
		expected := "10"
		if an != expected {
			t.Errorf("expected %s, found %s", expected, an)
		}
	})
	t.Run("assignment=11", func(t *testing.T) {
		an := AddLeadingZero(11)
		expected := "11"
		if an != expected {
			t.Errorf("expected %s, found %s", expected, an)
		}
	})
	t.Run("assignment=MAX_UINT32", func(t *testing.T) {
		an := AddLeadingZero(4294967295)
		expected := "4294967295"
		if an != expected {
			t.Errorf("expected %s, found %s", expected, an)
		}
	})
}

func TestAssignmentNumberFromFilename(t *testing.T) {
	t.Run("assignment-00", func(t *testing.T) {
		a := "assignment-00"
		expected := "00"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err != nil {
			t.Fatal(err)
		}
		if r != expected {
			t.Errorf("expected %s, found %s", expected, r)
		}
	})
	t.Run("assignment-1", func(t *testing.T) {
		a := "assignment-1"
		expected := "01"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err != nil {
			t.Fatal(err)
		}
		if r != expected {
			t.Errorf("expected %s, found %s", expected, r)
		}
	})
	t.Run("assignment-02", func(t *testing.T) {
		a := "assignment-02"
		expected := "02"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err != nil {
			t.Fatal(err)
		}
		if r != expected {
			t.Errorf("expected %s, found %s", expected, r)
		}
	})
	t.Run("assignment-10", func(t *testing.T) {
		a := "assignment-10"
		expected := "10"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err != nil {
			t.Fatal(err)
		}
		if r != expected {
			t.Errorf("expected %s, found %s", expected, r)
		}
	})
	t.Run("assignment-100", func(t *testing.T) {
		a := "assignment-100"
		expected := "100"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err != nil {
			t.Fatal(err)
		}
		if r != expected {
			t.Errorf("expected %s, found %s", expected, r)
		}
	})
	t.Run("assignment-010", func(t *testing.T) {
		a := "assignment-010"
		expected := "10"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err != nil {
			t.Fatal(err)
		}
		if r != expected {
			t.Errorf("expected %s, found %s", expected, r)
		}
	})
	t.Run("assignment-a010", func(t *testing.T) {
		a := "assignment-a010"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err == nil {
			t.Errorf("expected error, found %s", r)
		}
	})
	t.Run("assignment--1", func(t *testing.T) {
		a := "assignment--1"
		r, err := AssignmentNumberFromRegex(AssignmentDirectoryPattern, a)
		if err == nil {
			t.Errorf("expected error, found %s", r)
		}
	})
}
