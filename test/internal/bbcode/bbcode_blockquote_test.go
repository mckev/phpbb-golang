package bbcode

import (
	"bytes"
	"html/template"
	"testing"

	"phpbb-golang/internal/bbcode"
	"phpbb-golang/model"
)

func TestConvertBbcodeToHtml_BlockQuote(t *testing.T) {
	// 1234567890 is 2009-02-13 23:31:30 +0000 UTC
	bbcodeStr := `[blockquote user_name="User" user_id="123" post_id="456" time="1234567890"]text[/blockquote]`
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<blockquote><div><cite><a href="./users?u=123">User</a> wrote: <a aria-label="View quoted post" href="./posts?p=456#p456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a><span class="responsive-hide">13 Feb 09 23:31 UTC</span></cite>text</div></blockquote>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestConvertBbcodeToHtml_BlockQuoteBBAttack(t *testing.T) {
	userName := `"]an escape[/blockquote]<script>alert('Test XSS User name')</script>`

	// Use Go HTML Template to simulate rendered HTML
	const templateString = `[blockquote user_name="{{ .User.UserName }}" user_id="{{ .User.UserId }}" post_id="456" time="1234567890"]text[/blockquote]`
	templateOutput, err := template.New("").Parse(templateString)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	type TestPageData struct {
		User model.User
	}
	testPageData := TestPageData{
		User: model.User{
			UserName: userName,
			UserId:   123,
		},
	}
	var buffer bytes.Buffer
	err = templateOutput.Execute(&buffer, testPageData)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	bbcodeStr := buffer.String()
	{
		actual := bbcodeStr
		expected := `[blockquote user_name="&#34;]an escape[/blockquote]&lt;script&gt;alert(&#39;Test XSS User name&#39;)&lt;/script&gt;" user_id="123" post_id="456" time="1234567890"]text[/blockquote]`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
		expected := `<blockquote><div><cite><a href="./users?u=123">&amp;#34;]an escape[/blockquote]&amp;lt;script&amp;gt;alert(&amp;#39;Test XSS User name&amp;#39;)&amp;lt;/script&amp;gt;</a> wrote: <a aria-label="View quoted post" href="./posts?p=456#p456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a><span class="responsive-hide">13 Feb 09 23:31 UTC</span></cite>text</div></blockquote>`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestConvertBbcodeToHtml_BlockQuoteXss(t *testing.T) {
	bbcodeStr := `[blockquote=<script>alert('Test BB Attack')</script> user_name="User<script>alert('Test BB Attack')</script>" <script> user_id="123<script>alert('Test BB Attack')</script>" post_id="<script>alert('Test BB Attack')</script>456" time="<script>alert('Test BB Attack')</script>" <script>="<script>"]a <script>alert('Test BB Attack')</script> test[/blockquote]`
	actual := bbcode.ConvertBbcodeToHtml(bbcodeStr)
	expected := `<blockquote><div><cite><a href="./users?u=123&lt;script&gt;alert(&#39;Test BB Attack&#39;)&lt;/script&gt;">User&lt;script&gt;alert(&#39;Test BB Attack&#39;)&lt;/script&gt;</a> wrote: <a aria-label="View quoted post" href="./posts?p=&lt;script&gt;alert(&#39;Test BB Attack&#39;)&lt;/script&gt;456#p%3Cscript%3Ealert(%27Test%20BB%20Attack%27)%3C/script%3E456"><i aria-hidden="true" class="icon fa-arrow-circle-up fa-fw"></i></a><span class="responsive-hide">01 Jan 70 00:00 UTC</span></cite>a &lt;script&gt;alert(&#39;Test BB Attack&#39;)&lt;/script&gt; test</div></blockquote>`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}
