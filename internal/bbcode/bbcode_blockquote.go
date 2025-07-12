package bbcode

import (
	"fmt"
	"time"

	"github.com/frustra/bbcode"

	"phpbb-golang/internal/helper"
	"phpbb-golang/model"
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
	in := node.GetOpeningTag()
	username := ""
	if val, ok := in.Args["user_name"]; ok && val != "" {
		username = val
	}
	userid := model.INVALID_USER_ID
	if val, ok := in.Args["user_id"]; ok && val != "" {
		userid = helper.StrToInt(val, model.INVALID_USER_ID)
	}
	postid := model.INVALID_POST_ID
	if val, ok := in.Args["post_id"]; ok && val != "" {
		postid = helper.StrToInt(val, model.INVALID_POST_ID)
	}
	var unixTime int64
	if val, ok := in.Args["time"]; ok && val != "" {
		unixTime = helper.StrToInt64(val, 0)
	}
	now := time.Now().UTC()
	currentTime := now.Unix()

	blockquoteHtmlTag := bbcode.NewHTMLTag("")
	blockquoteHtmlTag.Name = "blockquote"
	divHtmlTag := bbcode.NewHTMLTag("")
	divHtmlTag.Name = "div"
	citeHtmlTag := bbcode.NewHTMLTag("")
	citeHtmlTag.Name = "cite"
	if username != "" && userid != model.INVALID_USER_ID {
		userLinkTag := bbcode.NewHTMLTag("")
		userLinkTag.Name = "a"
		userLinkTag.Attrs = map[string]string{
			"href": bbcode.ValidURL(fmt.Sprintf("./users?u=%d", userid)),
		}
		userLinkTag.AppendChild(bbcode.NewHTMLTag(username))
		citeHtmlTag.AppendChild(userLinkTag)
		citeHtmlTag.AppendChild(bbcode.NewHTMLTag(" wrote: "))
	} else {
		citeHtmlTag.AppendChild(bbcode.NewHTMLTag("Quote"))
	}
	if postid > 0 {
		postLinkTag := bbcode.NewHTMLTag("")
		postLinkTag.Name = "a"
		postLinkTag.Attrs = map[string]string{
			"href":       bbcode.ValidURL(fmt.Sprintf("./posts?p=%d#p%d", postid, postid)),
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
	}
	if unixTime > 0 && unixTime <= currentTime {
		timeTag := bbcode.NewHTMLTag("")
		timeTag.Name = "span"
		timeTag.Attrs = map[string]string{
			"class": "responsive-hide",
		}
		timeTag.AppendChild(bbcode.NewHTMLTag(helper.UnixTimeToStr(unixTime)))
		citeHtmlTag.AppendChild(timeTag)
	}
	divHtmlTag.AppendChild(citeHtmlTag)

	// Process things within [blockquote]...[/blockquote], including another [blockquote]
	for _, child := range node.Children {
		htmlChild := node.Compiler.CompileTree(child)
		// Inside [blockquote], convert back all <br> tag into new line "\n"
		htmlChild = convertBrTagIntoNewLine(htmlChild)
		divHtmlTag.AppendChild(htmlChild)
	}

	blockquoteHtmlTag.AppendChild(divHtmlTag)
	return blockquoteHtmlTag, false
}

func convertBrTagIntoNewLine(tag *bbcode.HTMLTag) *bbcode.HTMLTag {
	if tag == nil {
		return nil
	}

	if tag.Name == "br" {
		out := bbcode.NewHTMLTag("")
		out.Value = "\n"
		return out
	}

	// Create a copy of the tag, and recursively filter its children
	// Ref: https://github.com/frustra/bbcode/blob/master/html.go
	out := bbcode.NewHTMLTag("")
	out.Name = tag.Name
	out.Value = tag.Value
	out.Attrs = tag.Attrs
	for _, child := range tag.Children {
		out.AppendChild(convertBrTagIntoNewLine(child))
	}
	return out
}
