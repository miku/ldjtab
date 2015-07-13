package ldjtab

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const Version = "0.1.3"

var (
	errKeyNotFound   = fmt.Errorf("key not found")
	errValueNotFound = fmt.Errorf("value not found")
	errInvalidType   = fmt.Errorf("invalid type")
)

// StringValue returns the value for a given key in dot notation. Work
// resursively for objects, but not lists of objects.
func StringValue(key string, doc map[string]interface{}) (string, error) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return "", errKeyNotFound
	}

	val, ok := doc[keys[0]]
	if !ok {
		return "", errKeyNotFound
	}

	switch val.(type) {
	case string:
		return val.(string), nil
	case map[string]interface{}:
		if len(keys) < 2 {
			return "", errValueNotFound
		}
		return StringValue(keys[1], val.(map[string]interface{}))
	case json.Number:
		return fmt.Sprintf("%s", val), nil
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', 6, 64), nil
	case int:
		return strconv.Itoa(val.(int)), nil
	default:
		return "", errInvalidType
	}

	return "", errValueNotFound
}
