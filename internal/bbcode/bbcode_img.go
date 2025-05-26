package bbcode

import (
	"github.com/frustra/bbcode"
)

func imgBBTagHandler(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
	// Override default tag [img]
	// Ref: https://github.com/frustra/bbcode, https://github.com/frustra/bbcode/blob/master/compiler.go
	out := bbcode.NewHTMLTag("")
	out.Name = "img"
	value := node.GetOpeningTag().Value
	if value == "" {
		out.Attrs["src"] = bbcode.ValidURL(bbcode.CompileText(node))
	} else {
		out.Attrs["src"] = bbcode.ValidURL(value)
		text := bbcode.CompileText(node)
		if len(text) > 0 {
			out.Attrs["alt"] = text
			out.Attrs["title"] = out.Attrs["alt"]
		}
	}
	out.Attrs["referrerpolicy"] = "no-referrer"
	return out, false
}
