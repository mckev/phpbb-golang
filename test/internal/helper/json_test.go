package helper

import (
	"encoding/json"
	"testing"

	"phpbb-golang/internal/helper"
)

const testJsonString = `{"name":"Michael", "age":18, "athlete": true, "GPA": null, "hobbies": ["martial arts", "piano"]}`

func TestJsonDumps(t *testing.T) {
	jsonObj := map[string]any{}
	err := json.Unmarshal([]byte(testJsonString), &jsonObj)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	actual := helper.JsonDumps(jsonObj)
	expected := `{"GPA":null,"age":18,"athlete":true,"hobbies":["martial arts","piano"],"name":"Michael"}`
	if actual != expected {
		t.Errorf("Got %v, wanted %v", actual, expected)
	}
}
