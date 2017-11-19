package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type PostUploadsResponse struct {
	Timestamp    int64
	PresignedUrl string `json:"presigned_url"`
}

type ApiClient struct {
	Host  string
	Token string
}

func (client *ApiClient) setApiHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+client.Token)
}

func (client *ApiClient) uploadImageToS3(presignedUrl string, imagePath string) error {
	buf, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodPut, presignedUrl, bytes.NewReader(buf))
	req.Header.Add("Content-Type", "image/jpeg")
	_, err = client.httpClient().Do(req)

	return err
}

func (client *ApiClient) httpClient() *http.Client {
	return &http.Client{Timeout: 20 * time.Second}
}

func (client *ApiClient) createUpload() (*PostUploadsResponse, error) {
	req, _ := http.NewRequest(http.MethodPost, client.Host+"/api/uploads", nil)
	client.setApiHeaders(req)

	resp, err := client.httpClient().Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &PostUploadsResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (client *ApiClient) createImage(timestamp int64) error {
	payload := fmt.Sprintf(`{"timestamp": %d}`, timestamp)

	req, err := http.NewRequest(http.MethodPost, client.Host+"/api/images", strings.NewReader(payload))
	if err != nil {
		return err
	}

	client.setApiHeaders(req)

	_, err = client.httpClient().Do(req)

	return err
}
