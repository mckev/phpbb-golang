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
