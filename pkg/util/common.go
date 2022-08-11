package util

import "encoding/json"

func Marshal(m map[string]interface{}) string {
	if data, err := json.Marshal(m); err != nil {
		//log
		return ""
	} else {
		return string(data)
	}
}
