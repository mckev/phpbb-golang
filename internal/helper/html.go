package helper

import "strings"

func FormatAttributeValue(str string) string {
	// Ref: editor.js at formatAttributeValue()
	// WARNING: This is not safe! Please always use this in conjunction with Go HTML Template whenever using this function.

	// Check for any of: space, ', ", \, ]
	if !strings.ContainsAny(str, " \"'\\]") {
		return str
	}

	// Escape for single-quoted version
	singleEscaped := strings.ReplaceAll(str, `\`, `\\`)
	singleEscaped = strings.ReplaceAll(singleEscaped, `'`, `\'`)
	singleQuoted := "'" + singleEscaped + "'"

	// Escape for double-quoted version
	doubleEscaped := strings.ReplaceAll(str, `\`, `\\`)
	doubleEscaped = strings.ReplaceAll(doubleEscaped, `"`, `\"`)
	doubleQuoted := `"` + doubleEscaped + `"`

	// Return the shorter one
	if len(singleQuoted) < len(doubleQuoted) {
		return singleQuoted
	} else {
		return doubleQuoted
	}
}
