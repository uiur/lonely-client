package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUpload(t *testing.T) {
	token := "foobar"

	handler := func(w http.ResponseWriter, r *http.Request) {
		actual := r.Header.Get("Authorization")
		expected := "Bearer foobar"
		if actual != expected {
			t.Errorf("expected `Authorization` header: `%v`, got `%v`", expected, actual)
		}

		io.WriteString(
			w,
			`{"timestamp": 1000, "presigned_url": "https://aws.com/foobar"}`,
		)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))

	apiClient := &ApiClient{Host: ts.URL, Token: token}
	resp, err := apiClient.createUpload()

	if err != nil {
		t.Error(err)
	}

	if resp.Timestamp != 1000 {
		t.Errorf("expected valid timestamp: %v, actual: %v", 1000, resp.Timestamp)
	}

	expected := "https://aws.com/foobar"
	if resp.PresignedUrl != expected {
		t.Errorf("expected valid presigned url: %v, actual: %v", expected, resp.PresignedUrl)
	}
}

func TestCreateImage(t *testing.T) {
	timestamp := int64(1234)

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(
			w,
			fmt.Sprintf(`{"timestamp": %d}`, timestamp),
		)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))

	apiClient := &ApiClient{Host: ts.URL, Token: "foobar"}
	err := apiClient.createImage(timestamp)
	if err != nil {
		t.Error(err)
	}
}
