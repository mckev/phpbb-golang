package bbcode

import (
	"testing"

	"phpbb-golang/internal/bbcode"
)

func TestConvertBbcodeToHtml_UrlBasic(t *testing.T) {
	{
		bbcodeStr := `[url]https://www.example.com/[/url]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<a href="./redirect?url=https%3A%2F%2Fwww.example.com%2F" target="_blank">https://www.example.com/</a>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		bbcodeStr := `[url=https://www.example.com/]click here[/url]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<a href="./redirect?url=https%3A%2F%2Fwww.example.com%2F" target="_blank">click here</a>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestConvertBbcodeToHtml_UrlComplex(t *testing.T) {
	{
		bbcodeStr := `[url]https://www.google.com/search?q=how+to+make+a+raspberry+pi+web+server&hl=en&source=hp&ei=abcdef[/url]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<a href="./redirect?url=https%3A%2F%2Fwww.google.com%2Fsearch%3Fq%3Dhow%2Bto%2Bmake%2Ba%2Braspberry%2Bpi%2Bweb%2Bserver%26hl%3Den%26source%3Dhp%26ei%3Dabcdef" target="_blank">https://www.google.com/search?q=how+to+make+a+raspberry+pi+web+server&amp;hl=en&amp;source=hp&amp;ei=abcdef</a>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		bbcodeStr := `[url=https://www.google.com/search?q=how+to+make+a+raspberry+pi+web+server&hl=en&source=hp&ei=abcdef]click here[/url]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<a href="./redirect?url=https%3A%2F%2Fwww.google.com%2Fsearch%3Fq%3Dhow%2Bto%2Bmake%2Ba%2Braspberry%2Bpi%2Bweb%2Bserver%26hl%3Den%26source%3Dhp%26ei%3Dabcdef" target="_blank">click here</a>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestConvertBbcodeToHtml_UrlXss(t *testing.T) {
	{
		bbcodeStr := `[url]https://www.<script>example.com/[/url]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<a href="./redirect?url=https%3A%2F%2Fwww.%3Cscript%3Eexample.com%2F" target="_blank">https://www.&lt;script&gt;example.com/</a>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		bbcodeStr := `[url=https://www.<script>example.com/]click <script> here[/url]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<a href="./redirect?url=https%3A%2F%2Fwww.%3Cscript%3Eexample.com%2F" target="_blank">click &lt;script&gt; here</a>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
