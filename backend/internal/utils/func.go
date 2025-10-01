package utils

import (
	"encoding/json"

	"cloud.google.com/go/firestore"
)

func StructToUpdate(data any) ([]firestore.Update, error) {
	var m map[string]any
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	updates := make([]firestore.Update, 0, len(m))
	for k, v := range m {
		if validateValue(v) {
			updates = append(updates, firestore.Update{Path: k, Value: v})
		}
	}

	return updates, nil
}

func validateValue(v any) bool {
	switch val := v.(type) {
	case string:
		return val != ""
	case float64:
		return val != 0
	case []any:
		return len(val) > 0
	case map[string]any:
		return len(val) > 0
	case bool:
		return true
	default:
		return v != nil
	}
}
