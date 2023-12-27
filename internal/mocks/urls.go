package mocks

import "github.com/sntegegn/url-shortening/internal/models"

type URLModel struct{}

func (m *URLModel) Insert(shortKey, longURL string) error {
	return nil
}

func (m *URLModel) Get(shortKey string) (string, error) {
	if shortKey == "mUV4W2" {
		return "http://example.com", nil
	}
	return "", models.ErrRecordNotFound
}
