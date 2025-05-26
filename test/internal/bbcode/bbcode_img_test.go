package bbcode

import (
	"testing"

	"phpbb-golang/internal/bbcode"
)

func TestConvertBbcodeToHtml_ImgBasic(t *testing.T) {
	{
		// https://en.wikipedia.org/wiki/BBCode
		bbcodeStr := `[img]https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png[/img]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<img referrerpolicy="no-referrer" src="https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png">`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		bbcodeStr := `[img=https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png]This is just an example[/img]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<img alt="This is just an example" referrerpolicy="no-referrer" src="https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png" title="This is just an example">`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestConvertBbcodeToHtml_ImgXss(t *testing.T) {
	{
		// https://en.wikipedia.org/wiki/BBCode
		bbcodeStr := `[img]https://upload.wikimedia.org/<script>wikipedia/commons/7/70/Example.png[/img]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<img referrerpolicy="no-referrer" src="https://upload.wikimedia.org/%3Cscript%3Ewikipedia/commons/7/70/Example.png">`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		bbcodeStr := `[img=https://upload.wikimedia.org/<script>wikipedia/commons/7/70/Example.png]This is just a <script> example[/img]`
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<img alt="This is just a &lt;script&gt; example" referrerpolicy="no-referrer" src="https://upload.wikimedia.org/%3Cscript%3Ewikipedia/commons/7/70/Example.png" title="This is just a &lt;script&gt; example">`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
