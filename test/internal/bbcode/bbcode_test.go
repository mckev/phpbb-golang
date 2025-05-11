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

func TestConvertBbcodeToHtml_BlockQuote(t *testing.T) {
	bbcodeStr := `[blockquote user_name="User" user_id="123" post_id="456"]text[/blockquote]`
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<blockquote><div><cite><a href="/users?u=123">User</a> wrote: <a aria-label="View quoted post" href="/topics?p=456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a></cite>text</div></blockquote>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestConvertBbcodeToHtml_BlockQuoteXss(t *testing.T) {
	bbcodeStr := `[blockquote user_name="User<script>" user_id="123<script>" post_id="<script>456" <script>="<script>"]text <script> text[/blockquote]`
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<blockquote><div><cite><a href="/users?u=123&lt;script&gt;">User&lt;script&gt;</a> wrote: <a aria-label="View quoted post" href="/topics?p=&lt;script&gt;456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a></cite>text &lt;script&gt; text</div></blockquote>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}
