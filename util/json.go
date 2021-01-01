package util

import "encoding/json"

// Map2BYTES はマップを json バイト列に変換します。
func Map2BYTES(data map[string]string) ([]byte, error) {
	return json.Marshal(data)
}

// JSON2Map は json をマップに変換します。
func JSON2Map(jsonmap string) (map[string]string, error) {
	data := map[string]string{}
	if err := json.Unmarshal([]byte(jsonmap), &data); err != nil {
		return nil, err
	}
	return data, nil
}
