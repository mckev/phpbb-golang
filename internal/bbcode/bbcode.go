package bbcode

import (
	"github.com/frustra/bbcode"
)

func ConvertBbcodeToHtml(bbcodeStr string) string {
	// Ref: https://github.com/frustra/bbcode
	compiler := bbcode.NewCompiler(false, false)

	// Make HTML attributes deterministic for unit test
	compiler.SortOutputAttributes = true

	// Custom BB tags
	compiler.SetTag("blockquote", blockquoteBBTagHandler)

	html := compiler.Compile(bbcodeStr)
	return html
}

func blockquoteBBTagHandler(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
	// Custom tag [blockquote]
	// Ref: https://www.phpbb.com/community/viewtopic.php?t=2649439: Quotes are formatted like this in the database:  [quote="User" post_id="???" time="???" userid="???"]
	blockquoteHtmlTag := bbcode.NewHTMLTag("")
	blockquoteHtmlTag.Name = "blockquote"
	in := node.GetOpeningTag()
	username := ""
	if val, ok := in.Args["user_name"]; ok && val != "" {
		username = val
	}
	userid := ""
	if val, ok := in.Args["user_id"]; ok && val != "" {
		userid = val
	}
	postid := ""
	if val, ok := in.Args["post_id"]; ok && val != "" {
		postid = val
	}
	divHtmlTag := bbcode.NewHTMLTag("")
	divHtmlTag.Name = "div"
	citeHtmlTag := bbcode.NewHTMLTag("")
	citeHtmlTag.Name = "cite"
	userLinkTag := bbcode.NewHTMLTag("")
	userLinkTag.Name = "a"
	userLinkTag.Attrs = map[string]string{
		"href": bbcode.ValidURL("/users?u=" + userid),
	}
	userLinkTag.AppendChild(bbcode.NewHTMLTag(username))
	citeHtmlTag.AppendChild(userLinkTag)
	citeHtmlTag.AppendChild(bbcode.NewHTMLTag(" wrote: "))
	postLinkTag := bbcode.NewHTMLTag("")
	postLinkTag.Name = "a"
	postLinkTag.Attrs = map[string]string{
		"href":       bbcode.ValidURL("/topics?p=" + postid),
		"aria-label": "View quoted post",
	}
	postIconTag := bbcode.NewHTMLTag("")
	postIconTag.Name = "i"
	postIconTag.Attrs = map[string]string{
		"class":       "icon fa-arrow-circle-up fa-fw",
		"aria-hidden": "true",
	}
	postIconTag.AppendChild(bbcode.NewHTMLTag(""))
	postLinkTag.AppendChild(postIconTag)
	citeHtmlTag.AppendChild(postLinkTag)
	divHtmlTag.AppendChild(citeHtmlTag)
	text := bbcode.CompileText(node) // The text within [blockquote]...[/blockquote]
	divHtmlTag.AppendChild(bbcode.NewHTMLTag(text))
	blockquoteHtmlTag.AppendChild(divHtmlTag)
	return blockquoteHtmlTag, false
}
