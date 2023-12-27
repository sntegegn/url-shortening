package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func AddMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path))
		defer timer.ObserveDuration()
		totalRequest.WithLabelValues(r.URL.Path).Inc()
		next.ServeHTTP(w, r)
	})
}

func AddSpan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		span.SetAttributes(attribute.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
