package tracing

import (
	"fmt"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type TraceOption struct {
	ServiceName string
	ExporterOption
}

type ExporterOption struct {
	ExporterType string
	GcpProjectID string
	JaegerHost   string
}

func NewTracerProvider(opt TraceOption) (*trace.TracerProvider, error) {
	exporter, err := newExporter(opt.ExporterOption)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(opt.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func newExporter(opt ExporterOption) (exporter trace.SpanExporter, err error) {
	switch opt.ExporterType {
	case ExporterTypeGcp:
		exporter, err = texporter.New(texporter.WithProjectID(opt.GcpProjectID))
	case ExporterTypeJaegerCollector:
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(fmt.Sprintf("%s/api/traces", opt.JaegerHost))))
	case ExporterTypeJaegerAgent:
		exporter, err = jaeger.New(jaeger.WithAgentEndpoint())
	}

	return
}
