package mGin

import (
	"github.com/a5932016/go-ddd-example/util/mGin/mBinding"

	"github.com/gin-gonic/gin/binding"
)

// refer: https://github.com/go-playground/validator/blob/master/_examples/gin-upgrading-overriding/v8_to_v9.go
var (
	_ binding.StructValidator = newValidator()
)

func newValidator() *DefaultValidator {
	return &DefaultValidator{}
}

type DefaultValidator struct {
	mBinding.DefaultValidator
}
