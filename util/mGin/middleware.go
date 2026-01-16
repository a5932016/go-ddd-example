package mGin

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/a5932016/go-ddd-example/util/tracing"
)

const (
	_ContextKeyRequestBody = "requestBody"
	tracerKey              = "otel-go-contrib-tracer"
)

func RequestBodyToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)

		mCtx := NewContext(c)
		mCtx.setRequestBody(string(body))

		c.Next()
	}
}

// reference: go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin@v0.28.0/gintrace.go
func TracingMiddleware(service string) gin.HandlerFunc {
	cfg := tracing.NewConfig(
		tracing.WithSpanOptions(trace.WithSpanKind(trace.SpanKindClient)),
	)
	tracer := cfg.Tracer

	return func(c *gin.Context) {
		c.Set(tracerKey, tracer)
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		ctx := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.FullPath(), c.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		spanName := fmt.Sprintf("HTTP %s %s", c.Request.Method, c.Request.URL.Host+c.Request.URL.Path)
		if c.FullPath() == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		var b []byte
		if c.Request.Body != nil {
			b, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
		}

		tracing.SetSpanHttpAttritubes(
			span,
			tracing.WithRequestHeader(c.Request.Header),
			tracing.WithRequestBody(b),
			tracing.WithRequestQuery(c.Request.URL),
		)

		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)

		if c.Request.Response != nil {
			tracing.SetSpanHttpAttritubes(
				span,
				tracing.WithResponseBody(c.Request.Response),
				tracing.WithResponseHeader(c.Request.Response.Header),
			)
		}

		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}
