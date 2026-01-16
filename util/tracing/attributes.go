package tracing

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var omitKey = []string{
	"key", "Key", "Authorization",
}

func SetSpanHttpAttritubes(span trace.Span, opts ...func(trace.Span)) {
	for _, o := range opts {
		o(span)
	}
}

func WithRequestBody(b []byte) func(trace.Span) {
	return func(span trace.Span) {
		if len(b) == 0 {
			return
		}

		span.SetAttributes(attribute.String("body", string(b)))
	}
}

func WithRequestQuery(reqURL *url.URL) func(trace.Span) {
	return func(span trace.Span) {
		if reqURL == nil || len(reqURL.RawQuery) <= 0 {
			return
		}
		if parsedQuery, err := url.ParseQuery(reqURL.RawQuery); err == nil {
			for k, v := range parsedQuery {
				if isNeedOmit(k) {
					continue
				}
				span.SetAttributes(attribute.StringSlice(fmt.Sprintf("params.%s", k), v))
			}
		}
	}
}

func WithRequestHeader(header http.Header) func(trace.Span) {
	return func(span trace.Span) {
		for k, v := range header {
			if isNeedOmit(k) {
				continue
			}
			span.SetAttributes(attribute.StringSlice(fmt.Sprintf("request.header.%s", k), v))
		}
	}
}

func isNeedOmit(key string) bool {
	for _, t := range omitKey {
		if strings.Contains(key, t) {
			return true
		}
	}

	return false
}

func WithResponseHeader(header http.Header) func(trace.Span) {
	return func(span trace.Span) {
		for k, v := range header {
			span.SetAttributes(attribute.StringSlice(fmt.Sprintf("response.header.%s", k), v))
		}
	}
}

func WithResponseBody(res *http.Response) func(trace.Span) {
	return func(span trace.Span) {
		if res == nil {
			return
		}
		bodyBytes, _ := io.ReadAll(res.Body)
		res.Body.Close() //  must close
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		span.SetAttributes(attribute.String("response.body", string(bodyBytes)))
	}
}
