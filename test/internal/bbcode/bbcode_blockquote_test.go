package bbcode

import (
	"testing"

	"phpbb-golang/internal/bbcode"
)

func TestConvertBbcodeToHtml_BlockQuote(t *testing.T) {
	// 1234567890 is 2009-02-13 23:31:30 +0000 UTC
	bbcodeStr := `[blockquote user_name="User" user_id="123" post_id="456" time="1234567890"]text[/blockquote]`
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<blockquote><div><cite><a href="/users?u=123">User</a> wrote: <a aria-label="View quoted post" href="/topics?p=456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a><span class="responsive-hide">13 Feb 09 23:31 UTC</span></cite>text</div></blockquote>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestConvertBbcodeToHtml_BlockQuoteXss(t *testing.T) {
	bbcodeStr := `[blockquote user_name="User<script>" user_id="123<script>" post_id="<script>456" time="<script>" <script>="<script>"]text <script> text[/blockquote]`
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<blockquote><div><cite><a href="/users?u=123&lt;script&gt;">User&lt;script&gt;</a> wrote: <a aria-label="View quoted post" href="/topics?p=&lt;script&gt;456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a><span class="responsive-hide">01 Jan 70 00:00 UTC</span></cite>text &lt;script&gt; text</div></blockquote>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}
