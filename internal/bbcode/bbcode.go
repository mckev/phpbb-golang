package bbcode

import "github.com/frustra/bbcode"

func ConvertBbcodeToHtml(bbcodeStr string) string {
	// Ref: https://github.com/frustra/bbcode
	compiler := bbcode.NewCompiler(false, false)

	// Make HTML attributes deterministic for unit test
	compiler.SortOutputAttributes = true

	// Custom BB tags
	compiler.SetTag("blockquote", blockquoteBBTagHandler)
	compiler.SetTag("img", imgBBTagHandler)
	compiler.SetTag("url", urlBBTagHandler)

	html := compiler.Compile(bbcodeStr)
	return html
}
