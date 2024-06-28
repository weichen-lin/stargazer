package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agoda-com/opentelemetry-go/otelzap"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdklog "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

func InitTracerHTTP() *sdktrace.TracerProvider {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	endpoint := os.Getenv("TRACING_ENDPOINT")
	token := os.Getenv("OTEL_AUTH_TOKEN")
	env := os.Getenv("APP_ENV")

	if endpoint == "" || token == "" || env == "" {
		panic("Missing required otel environment variables")
	}

	otlptracehttp.NewClient()

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": token,
		}),
	)

	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter: ", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String("kafka-service"),
		semconv.ServiceVersionKey.String("0.0.1"),
		attribute.String("environment", env),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpHTTPExporter),
	)

	logExporter, _ := otlplogs.NewExporter(context.Background())

	otlplogshttp.NewClient()

	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithBatcher(logExporter),
		sdklog.WithResource(res),
	)

	logger := zap.New(otelzap.NewOtelCore(logProvider))
	zap.ReplaceGlobals(logger)
	otel.SetTracerProvider(tp)

	return tp
}
