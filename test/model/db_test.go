package model

import (
	"testing"

	"phpbb-golang/model"
)

func TestEscape01(t *testing.T) {
	sqlEscapeList := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}
	for sql, expected := range sqlEscapeList {
		actual := model.Escape(sql)
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestEscape02(t *testing.T) {
	sql := `<p>123</p><div><img width="1080" />`
	actual := model.Escape(sql)
	expected := `<p>123</p><div><img width=\"1080\" />`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}
