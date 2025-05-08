package helper

import "encoding/json"

func JsonDumps(jsonObj any) string {
	jsonBytes, _ := json.Marshal(jsonObj)
	// Note: Character &, <, and > are escaped into \u0026, \u003c, and \u003e to avoid certain safety problems that can arise when embedding JSON in HTML. Ref: https://pkg.go.dev/encoding/json#Encoder.SetEscapeHTML
	jsonString := string(jsonBytes)
	return jsonString
}
