package main

import (
	"github.com/google/wire"

	"github.com/a5932016/go-ddd-example/usecase"
)

var (
	usecaseProvider = wire.NewSet(
		usecase.NewHandler,
		wire.Bind(new(usecase.Handler), new(*usecase.HandlerConstructor)),
	)
)
