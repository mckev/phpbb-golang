package bbcode

import (
	"github.com/frustra/bbcode"

	"phpbb-golang/internal/helper"
)

func blockquoteBBTagHandler(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
	// Custom tag [blockquote]
	// Ref: https://www.phpbb.com/community/viewtopic.php?t=2649439: Quotes are formatted like this in the database:  [quote="User" post_id="???" time="???" userid="???"]
	// Hierarchy:
	//   - blockquoteHtmlTag <blockquote>
	//       - divHtmlTag <div>
	//           - citeHtmlTag <cite>
	//               - userLinkTag <a href=<user_id>>
	//                   - <user_name>
	//               - " wrote: "
	//               - postLinkTag <a href=<post_id>>
	//                   - postIconTag <i class="icon fa-arrow-circle-up fa-fw">
	//               - timeTag <span>
	//                   - <time>
	//           - <node.Children>
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
	var unixTime int64
	if val, ok := in.Args["time"]; ok && val != "" {
		unixTime = helper.StrToInt64(val, 0)
	}
	divHtmlTag := bbcode.NewHTMLTag("")
	divHtmlTag.Name = "div"
	citeHtmlTag := bbcode.NewHTMLTag("")
	citeHtmlTag.Name = "cite"
	userLinkTag := bbcode.NewHTMLTag("")
	userLinkTag.Name = "a"
	userLinkTag.Attrs = map[string]string{
		"href": bbcode.ValidURL("./users?u=" + userid),
	}
	userLinkTag.AppendChild(bbcode.NewHTMLTag(username))
	citeHtmlTag.AppendChild(userLinkTag)
	citeHtmlTag.AppendChild(bbcode.NewHTMLTag(" wrote: "))
	postLinkTag := bbcode.NewHTMLTag("")
	postLinkTag.Name = "a"
	postLinkTag.Attrs = map[string]string{
		"href":       bbcode.ValidURL("./posts?p=" + postid + "#p" + postid),
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
	timeTag := bbcode.NewHTMLTag("")
	timeTag.Name = "span"
	timeTag.Attrs = map[string]string{
		"class": "responsive-hide",
	}
	timeTag.AppendChild(bbcode.NewHTMLTag(helper.UnixTimeToStr(unixTime)))
	citeHtmlTag.AppendChild(timeTag)
	divHtmlTag.AppendChild(citeHtmlTag)
	// Process things within [blockquote]...[/blockquote], including another [blockquote]
	for _, child := range node.Children {
		htmlChild := node.Compiler.CompileTree(child)
		divHtmlTag.AppendChild(htmlChild)
	}
	blockquoteHtmlTag.AppendChild(divHtmlTag)
	return blockquoteHtmlTag, false
}
