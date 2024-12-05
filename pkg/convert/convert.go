package convert

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("convert - StructToMap - json.Marshal: %w", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil, fmt.Errorf("convert - StructToMap - json.Unmarshal: %w", err)
	}

	cleanedMap := make(map[string]interface{})
	for key, value := range m {
		if value != nil {
			cleanedMap[key] = value
		}
	}

	return cleanedMap, nil
}
