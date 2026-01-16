package mBinding

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/a5932016/go-ddd-example/util/mGin/bytesconv"
)

const filtersPkgPath = "github.com/a5932016/go-ddd-example/util/filters"

var errUnknownType = errors.New("unknown type")

var emptyField = reflect.StructField{}

type setOptions struct {
	isDefaultExists bool
	defaultValue    string
}

type setter interface {
	TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions, isMTag bool) (isSetted bool, err error)
}

type formSource map[string][]string

func (form formSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions, isMTag bool) (isSetted bool, err error) {
	if isMTag {
		return setFilterStruct(value, form, tagValue)
	}
	return setByForm(value, field, form, tagValue, opt)
}

func mapFormByTag(ptr interface{}, form map[string][]string, tag string) error {
	_, err := mapping(reflect.ValueOf(ptr), emptyField, formSource(form), tag)
	return err
}

// Recursive every attribute in target object, setter: query parameters
func mapping(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	if field.Tag.Get(tag) == "-" {
		return false, nil
	}

	var vKind = value.Kind()

	// value is a Ptr, mapping Elem of it
	if vKind == reflect.Ptr {
		var isNew bool
		vPtr := value
		if value.IsNil() {
			isNew = true
			vPtr = reflect.New(value.Type().Elem())
		}
		isSetted, err := mapping(vPtr.Elem(), field, setter, tag)
		if err != nil {
			return false, err
		}
		if isNew && isSetted {
			value.Set(vPtr)
		}
		return isSetted, nil
	}

	// value is filters.XXX => parse query string with []
	if value.Type().PkgPath() == filtersPkgPath {
		if ok, err := tryToSetValue(value, field, setter, tag, true); !ok {
			return false, err
		}
		return true, nil
	}

	// value is Struct, mapping its fields iterate
	if vKind == reflect.Struct {
		return mapFields(value, setter, tag, false)
	}

	// set value by property type
	if vKind != reflect.Struct || !field.Anonymous {
		if ok, err := tryToSetValue(value, field, setter, tag, false); !ok {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// Get tag and default setting, then try to set by query parameters
func tryToSetValue(value reflect.Value, field reflect.StructField, setter setter, tag string, isMTag bool) (bool, error) {
	var tagValue string
	var setOpt setOptions

	tagValue = field.Tag.Get(tag)
	tagValue, opts := head(tagValue, ",")

	if tagValue == "" {
		tagValue = field.Name
	}
	if tagValue == "" {
		return false, nil
	}

	var opt string
	for len(opts) > 0 {
		opt, opts = head(opts, ",")

		if k, v := head(opt, "="); k == "default" {
			setOpt.isDefaultExists = true
			setOpt.defaultValue = v
		}
	}

	return setter.TrySet(value, field, tagValue, setOpt, isMTag)
}

// Get value from query parameters by tag, then set it to specific field
func setByForm(value reflect.Value, field reflect.StructField, form map[string][]string, tagValue string, opt setOptions) (isSetted bool, err error) {
	vs, ok := form[tagValue]
	if !ok && !opt.isDefaultExists {
		return false, nil
	}

	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) > 0 {
			vs = strings.Split(vs[0], ",")
		}
		return true, setSlice(vs, value, field)
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) != value.Len() {
			return false, fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		}
		if len(vs) > 0 {
			vs = strings.Split(vs[0], ",")
		}
		return true, setArray(vs, value, field)
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}

		if len(vs) > 0 {
			val = vs[0]
		}
		return true, setWithProperType(val, value, field)
	}
}

/*
*
Parsing query parameters by regexp rule: {field}[{condition}]
example:

	tagValue = "user_id"
	form = ["user_id[in]": "uid1, uid2, uid3", "user_id[is]": "abc", "xxx[is_not]": "ghbjk"]
	get a map ["in": "uid1, uid2, uid3", "is": "abc"]

then mapping them by json tag at filters.Struct
*/
func setFilterStruct(value reflect.Value, form map[string][]string, tagValue string) (bool, error) {
	regRule := fmt.Sprintf("^%s\\[([a-z|_]+)\\]$", tagValue)
	r, err := regexp.Compile(regRule)
	if err != nil {
		return false, err
	}

	var isSetted bool

	for k, v := range form {
		regResult := r.FindStringSubmatch(k)
		mapTag := map[string][]string{}
		if len(regResult) == 2 {
			for i, value := range v {
				replacer := strings.NewReplacer("[", "", "]", "")
				v[i] = replacer.Replace(value)
			}
			mapTag[regResult[1]] = v

			isSetted, err = mapFields(value, formSource(mapTag), "json", true)
			if err != nil {
				return false, err
			}
		}
	}
	return isSetted, nil
}

// Mapping fields iterate, if onlyBindOnce is true, return when found the target field
func mapFields(value reflect.Value, setter setter, tag string, onlyBindOnce bool) (bool, error) {
	tValue := value.Type()
	var isSetted bool
	for i := 0; i < value.NumField(); i++ {
		sf := tValue.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}

		ok, err := mapping(value.Field(i), tValue.Field(i), setter, tag)
		if err != nil {
			return false, err
		}

		if onlyBindOnce {
			if ok {
				return true, nil
			}
		} else {
			isSetted = isSetted || ok
		}

	}
	return isSetted, nil
}

func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value, field)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		}
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal(bytesconv.StringToBytes(val), value.Addr().Interface())
	default:
		return errUnknownType
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil

	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, value reflect.Value, field reflect.StructField) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
