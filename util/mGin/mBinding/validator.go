package mBinding

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// refer: https://github.com/go-playground/validator/blob/master/_examples/gin-upgrading-overriding/v8_to_v9.go
var (
	Validator binding.StructValidator = newValidator()

	// use a single instance of Validate, it caches struct info
	AppValidate = validator.New()

	hostnameRegexRFC1123AndWildcardString = `^(\*\.)*([a-z0-9]+([-]+[a-z0-9]+)*\.)+[a-z]{2,}$`
	hostnameRegexRFC1123AndWildcard       = regexp.MustCompile(hostnameRegexRFC1123AndWildcardString)

	fqdnDomainRegexString = `^([a-z0-9]+([-]+[a-z0-9]+)*\.)+[a-z]{2,}$`
	fqdnDomainRegex       = regexp.MustCompile(fqdnDomainRegexString)
)

func newValidator() *DefaultValidator {
	return &DefaultValidator{}
}

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	valueType := kindOfData(obj)

	switch valueType {
	case reflect.Slice, reflect.Array:
		v.lazyinit()
		errMap := FieldErrors{}

		value := reflect.ValueOf(obj)
		valueType := value.Kind()
		if valueType == reflect.Ptr {
			value = value.Elem()
		}

		l := value.Len()
		for i := 0; i < l; i++ {
			item := value.Index(i).Interface()
			if err := v.ValidateStruct(item); err != nil {
				vErr, ok := err.(FieldError)
				if !ok {
					return err
				}
				errMap[i] = vErr
			}

		}

		if len(errMap) > 0 {
			errMap.SetIndex()
			return errMap
		}

	case reflect.Struct, reflect.Interface:
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			if vErrs, ok := err.(validator.ValidationErrors); ok {
				return NewFieldError(vErrs)
			}
			return err
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = AppValidate
		v.validate.SetTagName("binding")
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// add any custom validations etc. here
		v.validate.RegisterValidation("zone-domain", ValidateZoneDomain)
		v.validate.RegisterValidation("fqdn-domain", ValidateFQDNDomain)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// ValidateZoneDomain implements validator.Func
func ValidateZoneDomain(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return hostnameRegexRFC1123AndWildcard.MatchString(val)
}

// ValidateFQDNDomain implements validator.Func
func ValidateFQDNDomain(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return fqdnDomainRegex.MatchString(val)
}
