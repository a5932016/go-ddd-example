package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

const (
	ExporterTypeGcp             = "gcp"
	ExporterTypeJaegerCollector = "jaeger-collector"
	ExporterTypeJaegerAgent     = "jaeger-agent"
)

// Attribute keys that can be added to a span.
const (
	ReadBytesKey  = attribute.Key("http.read_bytes")  // if anything was read from the request body, the total number of bytes read
	ReadErrorKey  = attribute.Key("http.read_error")  // If an error occurred while reading a request, the string of the error (io.EOF is not recorded)
	WroteBytesKey = attribute.Key("http.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
	WriteErrorKey = attribute.Key("http.write_error") // if an error occurred while writing a reply, the string of the error (io.EOF is not recorded)

	SpanTypeKey   = attribute.Key("type")
	VMProviderKey = attribute.Key("vm.provider")
)

var (
	SpanTypeVM             = SpanTypeKey.String("vm")
	SpanTypeDNS            = SpanTypeKey.String("dns")
	SpanTypeAWSMarketPlace = SpanTypeKey.String("aws_marketplace")
)

// Server HTTP metrics
const (
	RequestCount          = "http.server.request_count"           // Incoming request count total
	RequestContentLength  = "http.server.request_content_length"  // Incoming request bytes total
	ResponseContentLength = "http.server.response_content_length" // Incoming response bytes total
	ServerLatency         = "http.server.duration"                // Incoming end to end duration, microseconds
)

// Filter is a predicate used to determine whether a given http.request should
// be traced. A Filter must return true if the request should be traced.
type Filter func(*http.Request) bool
