package helper

import "encoding/json"

func JsonDumps(jsonObj any) string {
	jsonBytes, _ := json.Marshal(jsonObj)
	jsonString := string(jsonBytes)
	return jsonString
}
