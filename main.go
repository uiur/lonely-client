package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("pass args\n")
		os.Exit(1)
	}

	imagePath := os.Args[1]

	token := os.Getenv("LONELY_DEVICE_TOKEN")
	if token == "" {
		fmt.Println("$LONELY_DEVICE_TOKEN is required")
		os.Exit(1)
	}

	host := os.Getenv("LONELY_SERVER_HOST")
	if host == "" {
		fmt.Println("$LONELY_SERVER_HOST is required")
		os.Exit(1)
	}

	apiClient := &ApiClient{Host: host, Token: token}

	uploadResponse, err := apiClient.createUpload()
	if err != nil {
		fmt.Printf("error: during a request to create upload: %v \n", err)
		os.Exit(1)
	}

	err = apiClient.uploadImageToS3(uploadResponse.PresignedUrl, imagePath)
	if err != nil {
		fmt.Printf("error: during a request to upload image to s3: %v \n", err)
		os.Exit(1)
	}

	err = apiClient.createImage(uploadResponse.Timestamp)

	if err != nil {
		fmt.Printf("error: during a request to register image: %v", err)
		os.Exit(1)
	}
}
