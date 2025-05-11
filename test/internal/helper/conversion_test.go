package helper

import (
	"phpbb-golang/internal/helper"
	"testing"
)

func TestStrToInt64(t *testing.T) {
	numstr := "9223372036854775807"
	actual := helper.StrToInt64(numstr, -1)
	expected := int64(9223372036854775807)
	if actual != expected {
		t.Errorf("Got %d, wanted %d", actual, expected)
		return
	}
}

func TestStrToInt64_Overflow(t *testing.T) {
	numstr := "9223372036854775808"
	actual := helper.StrToInt64(numstr, -1)
	expected := int64(-1)
	if actual != expected {
		t.Errorf("Got %d, wanted %d", actual, expected)
		return
	}
}

func TestStrToInt64_Invalid(t *testing.T) {
	numstr := "a223372036854775807"
	actual := helper.StrToInt64(numstr, -1)
	expected := int64(-1)
	if actual != expected {
		t.Errorf("Got %d, wanted %d", actual, expected)
		return
	}
}
