package helper

import (
	"testing"

	"phpbb-golang/internal/helper"
)

func TestIsEmailValid(t *testing.T) {
	if !helper.IsEmailValid("user@example.com") {
		t.Errorf("Expected true")
	}
	if helper.IsEmailValid("") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("user") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("@example.com") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("user@") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("user@.com") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("u ser@.com") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("user@example,com") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("user@example@example.com") {
		t.Errorf("Expected false")
	}
	if helper.IsEmailValid("user@example.com.") {
		t.Errorf("Expected false")
	}
}

func TestIsPasswordValid(t *testing.T) {
	if !helper.IsPasswordValid("Password1") {
		t.Errorf("Expected true")
	}
	if helper.IsPasswordValid("password1") {
		// No uppercase
		t.Errorf("Expected false")
	}
	if helper.IsPasswordValid("PASSWORD1") {
		// No lowercase
		t.Errorf("Expected false")
	}
	if helper.IsPasswordValid("Password") {
		// No digit
		t.Errorf("Expected false")
	}
	if helper.IsPasswordValid("Pass1") {
		// Too short
		t.Errorf("Expected false")
	}
}
