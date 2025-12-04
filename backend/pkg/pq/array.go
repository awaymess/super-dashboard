// Package pq provides PostgreSQL-specific types for GORM.
package pq

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

// StringArray represents a PostgreSQL text array.
type StringArray []string

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}

	// Format as PostgreSQL array literal
	quoted := make([]string, len(a))
	for i, s := range a {
		// Escape double quotes and backslashes
		escaped := strings.ReplaceAll(s, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		quoted[i] = `"` + escaped + `"`
	}
	return "{" + strings.Join(quoted, ",") + "}", nil
}

// Scan implements the sql.Scanner interface.
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("unsupported type for StringArray")
	}

	// Parse PostgreSQL array format: {val1,val2,val3}
	if str == "" || str == "{}" {
		*a = []string{}
		return nil
	}

	// Remove braces
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	// Simple parsing for unquoted values
	if !strings.Contains(str, `"`) {
		*a = strings.Split(str, ",")
		return nil
	}

	// Parse quoted values
	var result []string
	var current strings.Builder
	inQuote := false
	escaped := false

	for _, r := range str {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		switch r {
		case '\\':
			escaped = true
		case '"':
			inQuote = !inQuote
		case ',':
			if inQuote {
				current.WriteRune(r)
			} else {
				result = append(result, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	result = append(result, current.String())

	*a = result
	return nil
}

// MarshalJSON implements json.Marshaler.
func (a StringArray) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]string(a))
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *StringArray) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*a = arr
	return nil
}
