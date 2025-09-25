package helper

import (
	"testing"

	"phpbb-golang/internal/helper"
)

func TestSha256(t *testing.T) {
	actual := helper.Sha256("The quick brown fox jumps over the lazy dog")
	expected := "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592"
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestCrc32(t *testing.T) {
	{
		actual := helper.Crc32("The quick brown fox jumps over the lazy dog")
		expected := "414fa339"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.Crc32("")
		expected := "00000000"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestGenerateSessionId(t *testing.T) {
	sessionId, err := helper.GenerateSessionId()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if len(sessionId) != 32 {
		t.Errorf("Got %d, wanted %d", len(sessionId), 32)
		return
	}
	if !helper.IsSessionIdValid(sessionId) {
		t.Errorf("Expected true")
		return
	}
}

func TestGenerateRandomAlphanumeric(t *testing.T) {
	randomString, err := helper.GenerateRandomAlphanumeric(16)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if len(randomString) != 16 {
		t.Errorf("Got %d, wanted %d", len(randomString), 16)
		return
	}
}

func TestIsSessionIdValid(t *testing.T) {
	if helper.IsSessionIdValid("") {
		t.Errorf("Expected false")
		return
	}
	if !helper.IsSessionIdValid("11223344556677889900aabbccddeeff") {
		t.Errorf("Expected true")
		return
	}
	if !helper.IsSessionIdValid("11223344556677889900AaBbCcdDEeFF") {
		t.Errorf("Expected true")
		return
	}
	if helper.IsSessionIdValid("11223344556677889900aabbccddee") {
		t.Errorf("Expected false")
		return
	}
	if helper.IsSessionIdValid("11223344556677889900aabbccddeeff11") {
		t.Errorf("Expected false")
		return
	}
	if helper.IsSessionIdValid("11223344556677889900aabbccddeegg") {
		t.Errorf("Expected false")
		return
	}
}

func TestHashPasswordAndIsPasswordCorrect(t *testing.T) {
	{
		actual := helper.HashPassword("Password1", "MySalt1")
		expected := "sha256:MySalt1:9410cf91d1f4512b908e72875e57c091eaeab7bd588636aff4571b8a3b28e80d"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
		if !helper.IsPasswordCorrect("Password1", actual) {
			t.Errorf("Expected true")
			return
		}
		if helper.IsPasswordCorrect("Password2", actual) {
			t.Errorf("Expected false")
			return
		}
	}
	{
		actual := helper.HashPassword("Password1", "MySalt2")
		expected := "sha256:MySalt2:91758427fa428a517b5c8b9e49d0683d702f539a06c5cfef7feb631da2fd1a91"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
		if !helper.IsPasswordCorrect("Password1", actual) {
			t.Errorf("Expected true")
			return
		}
		if helper.IsPasswordCorrect("Password2", actual) {
			t.Errorf("Expected false")
			return
		}
	}
}
