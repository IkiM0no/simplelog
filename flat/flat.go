package flat

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func FlatMap(m map[string]interface{}) string {
	var buf bytes.Buffer
	for k, v := range m {
		switch v.(type) {
		case string, int, int32, int64, bool, float32, float64, uint64, error:
			buf.WriteString(fmt.Sprintf(`"%s"="%s" `, k, interfaceToString(v)))
		default:
			log.Println("not implemented. k: %v v: %v. type: %T", k, v, v)
		}
	}
	return buf.String()
}

func Flatten(nested map[string]interface{}, prefix string) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := _flatten(true, flatmap, nested, prefix)
	if err != nil {
		return nil, err
	}
	return flatmap, nil
}

var InvalidInputError = errors.New("invalid input: must be map or slice")

func _flatten(first bool, flatMap map[string]interface{}, nested interface{}, prefix string) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			if err := _flatten(false, flatMap, v, newKey); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := key(first, prefix, k)
			assign(newKey, v)
		}
	case []interface{}:
		for i, v := range nested.([]interface{}) {
			newKey := key(first, prefix, strconv.Itoa(i))
			assign(newKey, v)
		}
	default:
		return InvalidInputError
	}

	return nil
}

func key(first bool, prefix, subkey string) string {
	key := prefix

	if first {
		key += subkey
	} else {
		key += "_" + subkey
	}

	return key
}

func interfaceToString(inf interface{}) string {
	switch i := inf.(type) {
	case nil:
		return ""
	case string:
		return i
	case int:
		return strconv.Itoa(i)
	case int32:
		return strconv.FormatInt(int64(i), 10)
	case int64:
		return strconv.FormatInt(i, 10)
	case bool:
		return strconv.FormatBool(i)
	case float32:
		return strconv.FormatFloat(float64(i), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(i, 'f', -1, 64)
	case uint64:
		return strconv.FormatUint(i, 10)
	case error:
		return inf.(error).Error()
	default:
		return "-"
	}
}
