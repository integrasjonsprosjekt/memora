package utils

import (
	"encoding/json"
	"net/url"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func ParseFilter(filter string) ([]string, error) {
	var result []string
	decoded, err := url.QueryUnescape(filter)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitSeq(decoded, ",")
	for part := range parts {
		part = strings.TrimSpace(part)

		result = append(result, part)
	}
	return result, nil
}

func ReadDataFromIterator[T any](iter *firestore.DocumentIterator) ([]T, error) {
	var results []T

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		var item T
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

// StructToUpdate converts a struct to a slice of Firestore updates.
// It ignores zero-value fields to prevent overwriting existing data with empty values.
// Returns the slice of updates or an error if the conversion fails.
func StructToUpdate(data any) ([]firestore.Update, error) {
	// Marshal the struct into a JSON byte array
	var m map[string]any
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON byte array into a map
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	// Create Firestore updates, ignoring zero-value fields
	updates := make([]firestore.Update, 0, len(m))
	for k, v := range m {
		// Skip zero-value fields
		if validateValue(v) {
			updates = append(updates, firestore.Update{Path: k, Value: v})
		}
	}

	return updates, nil
}

// validateValue checks if a value is non-zero based on its type.
// Returns true if the value is non-zero, false otherwise.
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
