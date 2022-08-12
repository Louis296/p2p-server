package util

import (
	"encoding/json"
	"github.com/louis296/p2p-server/pkg/log"
)

func Marshal(m map[string]interface{}) string {
	if data, err := json.Marshal(m); err != nil {
		//log
		return ""
	} else {
		return string(data)
	}
}

func Unmarshal(str string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return data, nil
}
