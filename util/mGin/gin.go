package mGin

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/a5932016/go-ddd-example/util/filters"
	"github.com/a5932016/go-ddd-example/util/log"
	"github.com/a5932016/go-ddd-example/util/mError"
	"github.com/a5932016/go-ddd-example/util/mGin/mBinding"
	"github.com/a5932016/go-ddd-example/util/paging"
)

type HandlerFunc func(*Context)

type SortFilter struct {
	SortKey   string `json:"sortKey" form:"sortKey"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

type Meta struct {
	*paging.Paginator
	*SortFilter
	Code        int         `json:"code"`
	Status      string      `json:"status"` // success or fail
	Message     string      `json:"message"`
	Details     []detail    `json:"details,omitempty"`
	Errors      interface{} `json:"errors,omitempty"`
	DeclineCode string      `json:"decline_code,omitempty"`
}

type Wrap struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type detail struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type Context struct {
	*gin.Context
	wrap          Wrap
	err           error
	code          int
	detailFields  DetailFields
	fieldErrInfos []mBinding.FieldErrInfo
	customError   *CustomError
}

type DetailFields map[string]interface{}

var prefixCode *int

func SetResponseCodePrefix(prefix int) {
	prefixCode = &prefix
}

func NewContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
		wrap:    Wrap{},
	}
}

func (mGinFun HandlerFunc) GinFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContext(c)
		mGinFun(ctx)
	}
}

func (c *Context) GetSort() *filters.SortFilter {
	var (
		sortKey   string
		sortOrder string
	)

	sortKey, _ = c.GetQuery("sortKey")
	sortOrder, _ = c.GetQuery("sortOrder")
	switch strings.ToLower(sortOrder) {
	case "asc":
		return &filters.SortFilter{
			Asc: sortKey,
		}
	case "desc":
		return &filters.SortFilter{
			Desc: sortKey,
		}
	}

	return nil
}

func (c *Context) WithSort(sort *filters.SortFilter) *Context {
	if sort != nil {
		var (
			sortKey   string
			sortOrder string
		)

		if len(sort.Asc) > 0 {
			sortOrder = "asc"
			sortKey = sort.Asc
		} else {
			sortOrder = "desc"
			sortKey = sort.Desc
		}

		c.wrap.Meta.SortFilter = &SortFilter{
			SortKey:   sortKey,
			SortOrder: sortOrder,
		}
	}
	return c
}

func (c *Context) GetPaginator() paging.Paginator {
	var (
		page  int
		limit int
		err   error
	)
	pageStr, _ := c.GetQuery(paging.PageKeyName)
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
	} else {
		page = 1
	}
	limitStr, _ := c.GetQuery(paging.LimitKeyName)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit == 0 || limit > paging.DefaultMaxLimit {
			limit = paging.DefaultLimit
		}
	} else {
		limit = paging.DefaultLimit
	}

	offset := (page - 1) * limit
	return paging.Paginator{
		Limit:  limit,
		Page:   page,
		Offset: offset,
	}
}

// Response serializes the given struct as JSON into the response body.
func (c *Context) Response(httpCode int, msg string) {
	c.beforeResponse()
	if c.customError != nil {
		c.ResponseWithCustomError(*c.customError)
		return
	}
	if c.code == 0 {
		c.WithCode(httpCode)
	}
	c.wrap.Meta.Code = c.code
	c.wrap.Meta.Message = Tprintf(msg, c.detailFields)
	c.wrap.Meta.Status = c.status(httpCode)
	if c.fieldErrInfos != nil {
		c.wrap.Meta.Errors = c.fieldErrInfos
	}

	c.logError(httpCode)

	c.JSON(httpCode, c.wrap)
	c.Abort()
}

// ResponseWithCustomError response with custom error
func (c *Context) ResponseWithCustomError(cErr CustomError) {
	c.beforeResponse()
	c.customError = &cErr
	c.wrap.Meta.Code = c.customError.Code
	c.wrap.Meta.Message = c.customError.Message
	c.wrap.Meta.Status = c.status(c.customError.HTTPCode)
	c.wrap.Meta.Errors = c.customError.ErrorInfo
	c.wrap.Meta.DeclineCode = c.customError.DeclineCode

	c.logError(c.customError.HTTPCode)

	c.JSON(c.customError.HTTPCode, c.wrap)
	c.Abort()
}

// WithPaginator set paginator
func (c *Context) WithPaginator(page paging.Paginator) *Context {
	c.wrap.Meta.Paginator = &page
	return c
}

// WithData set response data
func (c *Context) WithData(data interface{}) *Context {
	c.wrap.Data = data
	return c
}

// WithError set error
func (c *Context) WithError(err error) *Context {
	c.err = err
	return c
}

// WithCode set code
func (c *Context) WithCode(code int) *Context {
	if code >= 1000 {
		log.Errorf("Invalid Error Code %d", code)
	}
	if prefixCode != nil {
		code = (*prefixCode * 1000) + code
	}
	c.code = code
	return c
}

// WithDetail TODO: deprecate
func (c *Context) WithDetail(key string, val interface{}) *Context {
	if c.detailFields != nil {
		c.detailFields[key] = val
	} else {
		c.detailFields = DetailFields{
			key: val,
		}
	}
	return c
}

// WithDetails TODO: deprecate
func (c *Context) WithDetails(details DetailFields) *Context {
	c.detailFields = details
	return c
}

func (c *Context) status(httpCode int) (status string) {
	if isSuccessHttpCode(httpCode) {
		status = StatusSuccess
	} else {
		status = StatusFail
	}
	return
}

func (c *Context) logError(httpCode int) {
	logLevel := logrus.InfoLevel
	if httpCode >= http.StatusInternalServerError {
		logLevel = logrus.ErrorLevel
	} else if httpCode >= http.StatusBadRequest {
		logLevel = logrus.WarnLevel
	} else {
		return
	}

	msg := c.wrap.Meta.Message
	_, file, line, _ := runtime.Caller(2)
	log := log.FromContext(c).WithFields(log.Fields{
		"httpSource":     fmt.Sprintf("%s:%d", file, line),
		"responseStatus": httpCode,
		"requestMethod":  c.Request.Method,
		"requestURL":     c.Request.URL.String(),
		// "requestBody":    c.getRequestBody(),
		"requestHeader": c.Request.Header,
		"err":           c.err,
		"fieldErr":      c.wrap.Meta.Errors,
	})

	log.Log(logLevel, msg)
}

func (c *Context) getDetails() []detail {
	fields := c.detailFields
	if fields == nil {
		return nil
	}
	ds := []detail{}
	if fields != nil {
		for k, v := range fields {
			t, err := findType(v)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"key":   k,
					"value": v,
					"err":   c.err,
					"code":  c.code,
				}).Error("Parse Response Detail Fail")
			}
			ds = append(ds, detail{
				Name:  k,
				Type:  t,
				Value: v,
			})
		}
	}
	return ds
}

func (c *Context) beforeResponse() *Context {
	c.handleError()
	return c
}

func (c *Context) handleError() {
	err := c.err
	if err != nil {
		err = mError.Unwrap(err)
	}
	switch mErr := err.(type) {
	case mBinding.FieldErrors:
		for _, err := range mErr {
			c.fieldErrInfos = append(c.fieldErrInfos, err.GetInfofs()...)
		}
		if c.code == 0 {
			c.WithCode(JsonValidationFail)
		}
	case mBinding.FieldError:
		c.fieldErrInfos = mErr.GetInfofs()
		if c.code == 0 {
			c.WithCode(JsonValidationFail)
		}
	case CustomError:
		c.customError = &mErr
		c.code = mErr.Code
	}
}

// Done invoke parent Done method
func (c *Context) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err invoke parent Err method
func (c *Context) Err() error {
	return c.Context.Err()
}

// Value invoke parent Value method
func (c *Context) Value(key interface{}) interface{} {
	if c.Request != nil {
		ctx := c.Request.Context()
		if v := ctx.Value(key); v != nil {
			return v
		}
	}
	return c.Context.Value(key)
}

func (c *Context) ShouldBindMQuery(obj interface{}) error {
	return c.ShouldBindWith(obj, mBinding.Query)
}

func (c *Context) setRequestBody(body string) {
	c.Set(_ContextKeyRequestBody, string(body))
}

func (c *Context) getRequestBody() string {
	val, exist := c.Get(_ContextKeyRequestBody)
	if !exist {
		return ""
	}

	if valStr, ok := val.(string); ok {
		return valStr
	}
	return ""
}
