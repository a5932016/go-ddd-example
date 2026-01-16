package paging

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// LimitKeyName is the request key name.
	LimitKeyName = "pageSize"

	// PageKeyName is the request page key name.
	PageKeyName = "pageCurrent"

	// DefaultLimit is the default number of items per page.
	DefaultLimit = 40

	// DefaultMaxLimit is the default number of max items per page.
	DefaultMaxLimit = 5000
)

// Paginator page data
type Paginator struct {
	TotalCount int `json:"recordCount"`
	TotalPage  int `json:"pageCount"`
	Page       int `json:"pageCurrent"`
	Limit      int `json:"pageSize"`
	Offset     int `json:"-"`
}

func GetPaginatorFromContext(c *gin.Context) Paginator {
	var (
		page  int
		limit int
		err   error
	)
	pageStr, _ := c.GetQuery(PageKeyName)
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
	} else {
		page = 1
	}
	limitStr, _ := c.GetQuery(LimitKeyName)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit == 0 || limit > DefaultMaxLimit {
			limit = DefaultLimit
		}
	} else {
		limit = DefaultLimit
	}

	offset := (page - 1) * limit
	return Paginator{
		Limit:  limit,
		Page:   page,
		Offset: offset,
	}
}

func (p *Paginator) SetTotalCount(count int) {
	p.TotalCount = count
	if p.Limit < 0 {
		// unlimit
		p.TotalPage = 1
	} else {
		p.TotalPage = p.TotalCount / p.Limit
	}
	if p.TotalCount%p.Limit > 0 || p.TotalPage == 0 {
		p.TotalPage++
	}
}
