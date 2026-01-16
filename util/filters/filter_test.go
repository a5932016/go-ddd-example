package filters

import (
	"reflect"
	"testing"

	"github.com/a5932016/go-ddd-example/util/types"
)

func TestSerializeListParams(t *testing.T) {
	type args struct {
		params interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"ok",
			args{NumberFilter{
				Is: types.Int(12),
				Gt: types.Int(11),
			}},
			map[string]string{
				"ttl[gt]": "11",
				"ttl[is]": "12",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SerializeListParams("ttl", tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializeListParams() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGenQueryParamsByStruct(t *testing.T) {
	tests := struct {
		args interface{}
		want map[string]string
	}{
		args: struct {
			A string       `json:"a"`
			B int          `json:"b"`
			C bool         `json:"c"`
			D StringFilter `json:"d"`
			E SortFilter   `json:"e"`
			F NumberFilter `json:"f"`
		}{
			A: "aaaaa",
			B: 1,
			C: true,
			D: StringFilter{
				IsNot: types.String("QAQ"),
			},
			E: SortFilter{
				Asc: "createdAt",
			},
			F: NumberFilter{
				In: []int{999, 234},
			},
		},
		want: map[string]string{
			"a":         "aaaaa",
			"b":         "1",
			"c":         "true",
			"d[is_not]": "QAQ",
			"e[asc]":    "createdAt",
			"e[desc]":   "",
			"f[in]":     "[\"999\",\"234\"]",
		},
	}

	if got := GenQueryParamsByStruct(tests.args); !reflect.DeepEqual(got, tests.want) {
		t.Errorf("GenQueryParamsByStruct() = %v, want %v", got, tests.want)
	}
}
