package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var totalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "total_http_requests",
		Help: "Number of requests handled by the urlshortner",
	},
	[]string{"path"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time",
		Help: "Duration of HTTP requests",
	},
	[]string{"path"},
)

var expandCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "expand_count",
		Help: "Number of expands",
	},
	[]string{"success"},
)

const (
	service     = "url_shortner"
	environment = "developement"
	version     = "1"
	id          = 1
)

func tracerProvider(ctx context.Context, endPoint, urlPath string) (*tracesdk.TracerProvider, error) {
	//Create the Jaeger exporter
	headers := map[string]string{
		"Content-Type": "application/x-protobuf",
	}
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endPoint),
		otlptracehttp.WithURLPath(urlPath),
		otlptracehttp.WithHeaders(headers),
	)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}
	rsc, rErr :=
		resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(service),
				semconv.ServiceVersionKey.String(version),
				attribute.String("environment", environment),
			),
		)

	if rErr != nil {
		panic(rErr)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(rsc),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)
	return tp, nil
}
