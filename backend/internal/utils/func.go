package utils

import (
	"encoding/json"
	"memora/internal/errors"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

const defaultLimitSize = 20

// ParseLimit parses a limit string and returns it as an integer.
// Returns an error if the string is not a valid integer.
func ParseLimit(limitStr string) (int, error) {
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}

	// If the limit is less than 1, default value of 20 used
	if limit < 1 {
		limit = defaultLimitSize
	}

	return limit, nil
}

// ParseFilter parses the query based on URI encoding
// and splits it into a slice of fields.
// Returns the slice of fields or an error if parsing fails.
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

// CheckIfUserCanAccessDeck checks if the user making the request
// is either the owner of the deck or has been granted access via shared emails.
// Returns true if the user can access the deck, false otherwise.
func CheckIfUserCanAccessDeck(c *gin.Context, ownerID string, sharedEmails []string) bool {
	uid, err := GetUID(c)
	if err != nil {
		return false
	}
	email, err := GetEmail(c)
	if err != nil {
		return false
	}
	if uid != ownerID && !slices.Contains(sharedEmails, email) {
		return false
	}
	return true
}

// GetUID retrieves the user ID (UID) from the Gin context.
// Returns the UID as a string or an error if not found.
func GetUID(c *gin.Context) (string, error) {
	uid, ok := c.Get("uid")
	if !ok {
		return "", errors.ErrUnauthorized
	}
	return uid.(string), nil
}

// GetEmail retrieves the user email from the Gin context.
// Returns the email as a string or an error if not found.
func GetEmail(c *gin.Context) (string, error) {
	email, ok := c.Get("email")
	if !ok {
		return "", errors.ErrUnauthorized
	}
	return email.(string), nil
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
	case time.Time:
		return !val.IsZero()
	default:
		return v != nil
	}
}
