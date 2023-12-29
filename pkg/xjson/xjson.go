package xjson

import (
	"encoding/json"

	"github.com/liushuangls/go-server-template/pkg/xstrings"
)

func MarshalString(val any) (string, error) {
	if val == nil {
		return "", nil
	}

	bs, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return xstrings.BytesToString(bs), nil
}

func MarshalStringNotE(val any) string {
	result, _ := MarshalString(val)
	return result
}

func UnmarshalString(buf string, val interface{}) error {
	return json.Unmarshal(xstrings.StringToBytes(buf), val)
}
