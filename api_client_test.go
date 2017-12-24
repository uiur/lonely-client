package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
func TestUploadImageToS3(t *testing.T) {
	token := "foobar"

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(
			w, "ok",
		)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))

	apiClient := &ApiClient{Host: ts.URL, Token: token}

	{
		err := apiClient.UploadImageToS3(ts.URL, "./fixtures/0.png")

		if !(err != nil && strings.Contains(err.Error(), "unexpected content type")) {
			t.Errorf("expected content type error, got: %v", err)
		}
	}

	{
		err := apiClient.UploadImageToS3(ts.URL, "./fixtures/empty.jpg")
		if !(err != nil && strings.Contains(err.Error(), "image is empty")) {
			t.Errorf("expected empty error, got: %v", err)
		}
	}

	{
		err := apiClient.UploadImageToS3(ts.URL, "./fixtures/0.jpg")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
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
