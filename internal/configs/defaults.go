package configs

import "encoding/json"

var defaults = map[string]json.RawMessage{
	KeySwaggerEnabled: json.RawMessage(`false`),
}

func Bool(key string) bool {
	v, ok := Configs.Get(key)
	if !ok {
		return false
	}
	var b bool
	if err := json.Unmarshal(v.Value, &b); err != nil {
		return false
	}
	return b
}
