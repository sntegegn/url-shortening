package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
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

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	t.Helper()

	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return res.StatusCode, res.Header, string(body)
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	t.Helper()

	res, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}
