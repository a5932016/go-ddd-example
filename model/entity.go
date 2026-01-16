package model

import "github.com/a5932016/go-ddd-example/util/filters"

type EntityOption struct {
	Keyword *string

	SortBy         *filters.SortFilter
	Op             string
	IncludeDeleted bool
	Offset         *int
	Limit          *int
}
