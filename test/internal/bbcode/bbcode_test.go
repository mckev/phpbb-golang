package bbcode

import (
	"testing"

	"phpbb-golang/internal/bbcode"
)

func TestConvertBbcodeToHtml_Basic(t *testing.T) {
	bbcodeStr := "[b]Hello World[/b]"
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := "<b>Hello World</b>"
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestConvertBbcodeToHtml_Quote(t *testing.T) {
	bbcodeStr := "[quote name=Somebody]text[/quote]"
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := "<blockquote><cite>Somebody said:</cite>text</blockquote>"
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestConvertBbcodeToHtml_UnmatchedClosingTags(t *testing.T) {
	bbcodeStr := "[center][b]text[/center]"
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<div style="text-align: center;">[b]text</div>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestConvertBbcodeToHtml_Ambiguous(t *testing.T) {
	// https://en.wikipedia.org/wiki/BBCode at "ambiguities"
	{
		bbcodeStr := "[quote=[b]text[/b][/quote]"
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `[quote=<b>text</b>[/quote]`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		bbcodeStr := "[quote=text[/quote]"
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `[quote=text[/quote]`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
