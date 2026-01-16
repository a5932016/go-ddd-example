package mGin

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"text/template"
)

// Tprintf passed template string is formatted usign its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func Tprintf(tmpl string, data map[string]interface{}) string {
	if data == nil {
		return tmpl
	}
	t := template.Must(template.New("tmpl").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func findType(v interface{}) (string, error) {
	switch t := v.(type) {
	default:
		return "", errors.New(fmt.Sprintf("unexpected type %T", t))
	case bool:
		return "bool", nil
	case string:
		return "string", nil
	case int, int32, int64, float32, float64:
		return "number", nil
	}
}

func isSuccessHttpCode(httpCode int) bool {
	return httpCode < http.StatusBadRequest
}
