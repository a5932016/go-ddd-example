package filters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Adaptor adaptor
type Adaptor interface {
	ToSQL(column string) (queries []string, args []interface{})
}

func GenQueryParamsByStruct(obj interface{}) map[string]string {
	result := map[string]string{}
	value := reflect.ValueOf(obj)

	if value.IsZero() || !value.IsValid() || value.IsNil() {
		return result
	}

	if value.Type().Kind() == reflect.Ptr {
		MergeQueryMap(result, SerializeListParams("", value.Elem().Interface()))
	} else {
		MergeQueryMap(result, SerializeListParams("", value.Interface()))
	}

	return result
}

func MergeQueryMap(query, item map[string]string) map[string]string {
	for k, v := range item {
		query[k] = v
	}
	return query
}


// SerializeListParams is to used to serialize the inputParams of list request.
// Input: prefix="ttl", params=NumberFilter{Is: types.Int(12),Gt: types.Int(11)}
// output: map[string]string{"ttl[gt]": "11", "ttl[is]": "12"}
func SerializeListParams(prefix string, params interface{}) map[string]string {
	queryParams, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	data := json.NewDecoder(strings.NewReader(string(queryParams)))
	data.UseNumber()
	var m map[string]interface{}
	if err := data.Decode(&m); err != nil {
		panic(err)
	}

	serListParams := make(map[string]interface{})
	parseMapListParams(m, prefix, serListParams)
	body := map[string]string{}
	for k, v := range serListParams {
		switch val := v.(type) {
		case []interface{}:
			value := "[\""
			str := []string{}
			for _, element := range val {
				str = append(str, fmt.Sprintf("%v", element))
			}
			value = value + strings.Join(str, "\",\"")
			value = value + "\"]"
			body[k] = value
		default:
			body[k] = fmt.Sprintf("%v", v)
		}
	}
	return body

}

func parseMapListParams(aMap map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, val := range aMap {
		switch value := val.(type) {
		case map[string]interface{}:
			parseMapListParams(val.(map[string]interface{}), key, result)
		default:
			if prefix != "" {
				k := prefix + "[" + key + "]"
				result[k] = value
			} else {
				result[key] = value
			}

		}
	}
}
