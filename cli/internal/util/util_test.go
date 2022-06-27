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
