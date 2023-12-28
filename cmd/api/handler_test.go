package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/sntegegn/url-shortening/internal/assert"
)

func TestShortenURL(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.route())
	defer ts.Close()

	const formTag = `<form action="/shorten" method="POST" class="max-w-sm mx-auto">`

	/* _, _, body := ts.get(t, "/shorten") */
	tests := []struct {
		name            string
		longURL         string
		ExpectedCode    int
		ExpectedFormTag string
	}{
		{
			name:         "Valid Submission",
			longURL:      "http://example.com",
			ExpectedCode: http.StatusOK,
		},
		{
			name:            "Empty LongURL",
			longURL:         "",
			ExpectedCode:    http.StatusBadRequest,
			ExpectedFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("longURL", tt.longURL)

			code, _, body := ts.postForm(t, "/shorten", form)
			assert.Equal(t, code, tt.ExpectedCode)

			if formTag != "" {
				assert.StringContains(t, body, tt.ExpectedFormTag)
			}

		})
	}
}

func TestExpandURL(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.route())
	defer ts.Close()

	tests := []struct {
		name            string
		shortKey        string
		expectedLongURL string
		statusCode      int
	}{
		{
			name:            "Valid ShortKey",
			shortKey:        "mUV4W2",
			expectedLongURL: "http://example.com",
			statusCode:      http.StatusOK,
		},
		{
			name:       "Invalid ShortKey",
			shortKey:   "abcdef",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.get(t, "/s/"+tt.shortKey)
			assert.Equal(t, code, tt.statusCode)
		})
	}
}
