package mBinding

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// TODO: need to be implement more validator tag
	// https://godoc.org/gopkg.in/go-playground/validator.v9
	fieldInputErrorCode = map[string]string{
		"required":         "The value is required",
		"eq":               "The value must be equal %s",
		"ne":               "The value must be not equal %s",
		"gt":               "The value must be greater than %s",
		"gte":              "The value must be greater than or equal %s",
		"lt":               "The value must be less than %s",
		"lte":              "The value must be less than or equal %s",
		"oneof":            "The value must be one of the values in %s",
		"zone-domain":      "The value must be FQDN or begin with (*.)",
		"unique":           "Values must be unique",
		"hostname_rfc1123": "The value must be a valid Hostname according to RFC 1123",
		"ipv4":             "The value must be a v4 IP Address",
		"fqdn":             "The value must be a valid FQDN",
	}
)

type FieldErrInfo struct {
	Index   string   `json:"index,omitempty"`
	Field   string   `json:"field"`
	Code    string   `json:"code"`
	Param   []string `json:"param,omitempty"`
	Message string   `json:"message"`
	value   interface{}
}

// FieldError field error
type FieldError struct {
	errs  validator.ValidationErrors
	infos []FieldErrInfo
}

type FieldErrors map[int]FieldError

func NewFieldError(errs validator.ValidationErrors) FieldError {
	fErrs := FieldError{errs: errs}
	for _, err := range errs {
		ferr := FieldErrInfo{
			Field: err.Field(),
			Code:  err.ActualTag(),
			value: err.Value(),
		}
		if len(err.Param()) > 0 {
			ferr.Param = strings.Split(err.Param(), " ")
		}
		if msg, ok := fieldInputErrorCode[err.ActualTag()]; ok {
			ferr.Message = msg
			if ferr.Param != nil {
				paramStr := fmt.Sprintf("{%s}", strings.Join(ferr.Param, ", "))
				ferr.Message = fmt.Sprintf(msg, paramStr)
			}
		} else {
			ferr.Message = "Invalid Input Value"
		}

		fErrs.infos = append(fErrs.infos, ferr)
	}
	return fErrs
}

// AppendErrorInfo append field error info
func (q *FieldError) AppendErrorInfo(field, code, message string) {
	q.infos = append(q.infos, FieldErrInfo{
		Field:   field,
		Code:    code,
		Message: message,
	})
}

func (q FieldError) Error() string {
	return "Invalid Input Value"
}

func (q FieldError) GetInfofs() []FieldErrInfo {
	return q.infos
}

func (q *FieldError) SetIndex(i int) {
	for k := range q.infos {
		q.infos[k].Index = strconv.Itoa(i)
	}
}

func (q FieldErrors) Error() string {
	return "Invalid Input Value"
}

func (q *FieldErrors) SetIndex() {
	for i, v := range *q {
		v.SetIndex(i)
	}
}
