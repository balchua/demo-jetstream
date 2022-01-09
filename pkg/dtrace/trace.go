package dtrace

import (
	"context"

	"github.com/balchua/demo-jetstream/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

type Dtrace struct {
	tp  *sdktrace.TracerProvider
	exp sdktrace.SpanExporter
}

func SetupTracer(tracingConfig config.Tracing) *Dtrace {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(tracingConfig.JaegerUrl)))
	if err != nil {
		return nil
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(newResource(tracingConfig.ServiceName)),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &Dtrace{
		tp:  tp,
		exp: exp,
	}
}

func newResource(serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("v0.1.0"),
		attribute.String("environment", "demo"),
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
