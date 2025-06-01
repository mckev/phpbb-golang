package bbcode

import (
	"net/url"

	"github.com/frustra/bbcode"
)

func urlBBTagHandler(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
	// Override default tag [url]
	// Ref: https://github.com/frustra/bbcode, https://github.com/frustra/bbcode/blob/master/compiler.go
	out := bbcode.NewHTMLTag("")
	out.Name = "a"
	value := node.GetOpeningTag().Value
	if value == "" {
		text := bbcode.CompileText(node)
		if len(text) > 0 {
			out.Attrs["href"] = bbcode.ValidURL("./redirect?url=" + url.QueryEscape(text))
		}
	} else {
		out.Attrs["href"] = bbcode.ValidURL("./redirect?url=" + url.QueryEscape(value))
	}
	out.Attrs["target"] = "_blank"
	return out, true
}
