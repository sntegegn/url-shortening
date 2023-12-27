package main

import (
	"io"
	"log/slog"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/sntegegn/url-shortening/internal/mocks"
	"go.opentelemetry.io/otel"
)

func newTestApplication(t *testing.T) *application {
	formDecoder := form.NewDecoder()
	tr := otel.Tracer("teset")
	app := application{
		config:      config{},
		logger:      slog.New(slog.NewJSONHandler(io.Discard, nil)),
		formDecoder: formDecoder,
		URLModel:    &mocks.URLModel{},
		tracer:      tr,
	}
	return &app
}

/* type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}
*/
