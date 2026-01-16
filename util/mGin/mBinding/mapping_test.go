package mBinding

import (
	"reflect"
	"strconv"
	"strings"
	_ "strings"
	"testing"

	"github.com/a5932016/go-ddd-example/util/filters"
)

func TestMTagBinding_Bind(t *testing.T) {
	type testStruct struct {
		Region *filters.StringFilter    `mTag:"region"`
		OrgID  *filters.StringFilter    `mTag:"org_id"`
		Num    *filters.NumberFilter    `mTag:"num"`
		Boo    *filters.BooleanFilter   `mTag:"bool"`
		Time   *filters.TimestampFilter `mTag:"time"`
		SortBy *filters.SortFilter      `mTag:"sortBy"`
	}

	obj := testStruct{}

	gt := 34567
	is := true
	regionIs := "VN"
	timeAfter := int64(97754)

	expect := testStruct{
		Region: &filters.StringFilter{
			Is: &regionIs,
			In: []string{"FGHJ", "GHUIJ", "TFYGUIHO"},
		},
		OrgID: &filters.StringFilter{
			In: []string{"cytgvhjk", "tyguhijlk"},
		},
		Num: &filters.NumberFilter{
			Gt:      &gt,
			Between: []int{12345, 45678},
		},
		Boo: &filters.BooleanFilter{
			Is: &is,
		},
		Time: &filters.TimestampFilter{
			After:   &timeAfter,
			Between: []int64{9722, 5678},
		},
		SortBy: &filters.SortFilter{
			Asc: "XDXD",
		},
	}

	values := map[string][]string{}
	values["region[is]"] = []string{*expect.Region.Is}
	values["region[in]"] = []string{strings.Join(expect.Region.In, ",")}
	values["org_id[in]"] = []string{strings.Join(expect.OrgID.In, ",")}
	values["num[gt]"] = []string{strconv.Itoa(*expect.Num.Gt)}
	values["num[between]"] = []string{numSliceToString(expect.Num.Between)}
	values["bool[is]"] = []string{strconv.FormatBool(*expect.Boo.Is)}
	values["time[after]"] = []string{strconv.Itoa(int(*expect.Time.After))}
	values["time[between]"] = []string{numSliceToString(expect.Time.Between)}
	values["sortBy[asc]"] = []string{expect.SortBy.Asc}

	if err := mapFormByTag(&obj, values, mTag); err != nil {
		t.Error(err)
	}

	objMap := structToMap(obj)
	expMap := structToMap(expect)
	if !reflect.DeepEqual(objMap, expMap) {
		t.Errorf("want: %v\n, get: %v", expMap, objMap)
	}
}

func numSliceToString(input interface{}) string {
	result := []string{}
	switch v := input.(type) {
	case []int:
		for _, i := range v {
			result = append(result, strconv.Itoa(i))
		}
	case []int64:
		for _, i := range v {
			result = append(result, strconv.Itoa(int(i)))
		}
	}
	return strings.Join(result, ",")
}

func structToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	inputValue := reflect.ValueOf(input)
	for i := 0; i < inputValue.NumField(); i++ {
		name := inputValue.Type().Field(i).Name
		value := inputValue.Field(i)
		if value.Kind() == reflect.Ptr {
			subResult := make(map[string]interface{})
			for j := 0; j < value.Elem().NumField(); j++ {
				subName := value.Elem().Type().Field(j).Name
				subValue := value.Elem().FieldByName(subName)
				tmp := getValueWithType(subValue)
				subResult[subName] = tmp
			}
			result[name] = subResult
			continue
		}
		result[name] = value
	}
	return result
}

func getValueWithType(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Invalid:
		return ""
	case reflect.Ptr:
		return getValueWithType(v.Elem())
	default:
		return v.Interface()
	}
}
