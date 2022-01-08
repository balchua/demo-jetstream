package dtrace

import (
	"context"
	"io"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

type Dtrace struct {
	tp  *sdktrace.TracerProvider
	exp sdktrace.SpanExporter
	f   *os.File
}

func SetupTracer() *Dtrace {
	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		return nil
	}
	exp, err := newExporter(f)
	if err != nil {
		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &Dtrace{
		tp:  tp,
		exp: exp,
		f:   f,
	}
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("publisher"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func (d *Dtrace) Close() {
	if err := d.tp.Shutdown(context.Background()); err != nil {
		zap.S().Fatal(err)
	}
}

func (d *Dtrace) Flush(ctx context.Context) {
	d.tp.ForceFlush(ctx)
}
