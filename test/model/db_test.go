package model

import (
	"fmt"
	"path"
	"testing"

	"phpbb-golang/model"
)

func TestSqlEscape01(t *testing.T) {
	tests := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}
	for sql, expected := range tests {
		actual := model.SqlEscape(sql)
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestSqlEscape02(t *testing.T) {
	sql := `<p>123</p><div><img width="1080" />`
	actual := model.SqlEscape(sql)
	expected := `<p>123</p><div><img width=\"1080\" />`
	if actual != expected {
		t.Errorf("Got %s, wanted %s", actual, expected)
		return
	}
}

func TestPathEscape(t *testing.T) {
	{
		tableName := "../../evil.txt"
		actual := fmt.Sprintf("file:./model/db/%s.db", path.Base(tableName))
		expected := `file:./model/db/evil.txt.db`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		tableName := "/../../"
		actual := fmt.Sprintf("file:./model/db/%s.db", path.Base(tableName))
		expected := `file:./model/db/...db`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
