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

func TestIsStringNFKCNormal(t *testing.T) {
	// Ambiguous characters are rejected
	if !helper.IsStringNFKCNormalized("ABC") {
		t.Errorf("Expected true")
	}
	if helper.IsStringNFKCNormalized("ＡＢＣ") {
		t.Errorf("Expected false")
	}
	if !helper.IsStringNFKCNormalized("Admin") {
		t.Errorf("Expected true")
	}
	if helper.IsStringNFKCNormalized("Ａｄｍｉｎ") {
		t.Errorf("Expected false")
	}
	if !helper.IsStringNFKCNormalized("Hello") {
		t.Errorf("Expected true")
	}
	if helper.IsStringNFKCNormalized("ℌ𝔢𝔩𝔩𝔬") {
		t.Errorf("Expected false")
	}
	if helper.IsStringNFKCNormalized("𝔥𝔢𝔩𝔩𝔬") {
		t.Errorf("Expected false")
	}
	if !helper.IsStringNFKCNormalized("LOL") {
		t.Errorf("Expected true")
	}
	if helper.IsStringNFKCNormalized("Ⓛⓞⓛ") {
		t.Errorf("Expected false")
	}

	// All below are okay
	if !helper.IsStringNFKCNormalized("Café") {
		t.Errorf("Expected true")
	}
	if !helper.IsStringNFKCNormalized("Straße") {
		t.Errorf("Expected true")
	}
	// Chinese "hello"
	if !helper.IsStringNFKCNormalized("你好") {
		t.Errorf("Expected true")
	}
	// Japanese kanji "Tokyo"
	if !helper.IsStringNFKCNormalized("東京") {
		t.Errorf("Expected true")
	}
	// Russian "good day"
	if !helper.IsStringNFKCNormalized("Добрый день") {
		t.Errorf("Expected true")
	}
	// Arabic "hello"
	if !helper.IsStringNFKCNormalized("مرحبا") {
		t.Errorf("Expected true")
	}
	// Hindi "computer"
	if !helper.IsStringNFKCNormalized("कंप्यूटर") {
		t.Errorf("Expected true")
	}
}
