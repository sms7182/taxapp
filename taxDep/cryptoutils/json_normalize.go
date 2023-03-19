package cryptoutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

var ErrNotJsonString = errors.New("input is not a JSON string")

func NormalizeJson(jsonString []byte) (string, error) {
	jMap := make(map[string]interface{})

	if err := json.Unmarshal(jsonString, &jMap); err != nil {
		return "", ErrNotJsonString
	}

	return NormalizeJsonObj(jMap)
}

func NormalizeJsonObj(jMap map[string]interface{}) (string, error) {
	return normalizeVals(sortJsonMap(jMap, "")), nil
}

type kv struct {
	k string
	v interface{}
}

func sortJsonMap(m map[string]interface{}, prefix string) []kv {
	var result []kv

	for k, v := range m {
		key := prefix
		if len(key) > 0 {
			key += k
		} else {
			key = k
		}

		if v == nil {
			result = append(result, kv{
				k: key,
				v: v,
			})
			continue
		}

		if obj, ok := isJsonObject(v); ok {
			result = append(result, sortJsonMap(obj, prefix+k+".")...)
			continue
		}

		if array, ok := isJsonObjectArray(v); ok {
			for _, obj := range array {
				result = append(result, sortJsonMap(obj.(map[string]interface{}), prefix+k+".")...)
			}
			continue
		}

		if array, ok := isJsonArray(v); ok {
			for _, val := range array {
				result = append(result, kv{
					k: key,
					v: val,
				})
			}
			continue
		}

		result = append(result, kv{
			k: key,
			v: v,
		})

	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].k < result[j].k
	})

	return result
}

func isJsonObject(v interface{}) (map[string]interface{}, bool) {
	obj, ok := v.(map[string]interface{})
	return obj, ok
}

func isJsonObjectArray(v interface{}) ([]interface{}, bool) {
	array, ok := isJsonArray(v)
	if !ok {
		return nil, false
	}

	_, ok = array[0].(map[string]interface{})

	return array, ok
}

func isJsonArray(v interface{}) ([]interface{}, bool) {
	array, ok := v.([]interface{})
	return array, ok
}

func normalizeVals(kvs []kv) string {
	var results []string

	for _, v := range kvs {
		results = append(results, normalizeVal(v.v))
	}

	return strings.Join(results, "#")
}

func normalizeVal(v interface{}) string {
	if v == nil {
		return "#"
	}

	switch val := v.(type) {
	case float64:
		return float64String(val)
	case string:
		if len(val) == 0 {
			return "#"
		}
		return strings.ReplaceAll(val, "#", "##")
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return ""
	}

}

func float64String(f float64) string {
	i, frac := math.Modf(f)
	if frac == 0 {
		return fmt.Sprintf("%d", int64(i))
	}
	return strings.TrimRight(fmt.Sprintf("%f", f), "0")
}
