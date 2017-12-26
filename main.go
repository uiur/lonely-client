package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "lonely"
	app.Usage = "lonely client"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "Captures a current image and uploads it",
			Action: func(c *cli.Context) error {
				err := checkEnv()
				if err != nil {
					return err
				}

				dir, err := ioutil.TempDir("", "lonely-")
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				defer os.RemoveAll(dir)

				imagePath := dir + "/snapshot.jpg"

				err = capture(imagePath)
				if err != nil {
					return err
				}

				err = upload(imagePath)

				return err
			},
		},
		{
			Name:      "upload",
			Usage:     "Uploads an image",
			ArgsUsage: "[image_path]",
			Action: func(c *cli.Context) error {
				err := checkEnv()
				if err != nil {
					return err
				}

				err = upload(c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func checkEnv() error {
	token := os.Getenv("LONELY_DEVICE_TOKEN")
	if token == "" {
		return cli.NewExitError("$LONELY_DEVICE_TOKEN is required", 1)
	}

	host := os.Getenv("LONELY_SERVER_HOST")
	if host == "" {
		return cli.NewExitError("$LONELY_SERVER_HOST is required", 1)
	}

	return nil
}

// Upload an image to lonely server
func upload(imagePath string) error {
	token := os.Getenv("LONELY_DEVICE_TOKEN")
	host := os.Getenv("LONELY_SERVER_HOST")

	apiClient := &ApiClient{Host: host, Token: token}

	uploadResponse, err := apiClient.createUpload()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error: during a request to create upload: %v \n", err), 1)
	}

	err = apiClient.UploadImageToS3(uploadResponse.PresignedUrl, imagePath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error: during a request to upload image to s3: %v \n", err), 1)
	}

	err = apiClient.createImage(uploadResponse.Timestamp)

	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error: during a request to register image: %v", err), 1)
	}

	return nil
}

// Capture an image from webcam
//
// requires following commands:
//   linux: `fswebcam`
//   mac:   `imagesnap`
func capture(imagePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("fswebcam", "--resolution", "1280x720", imagePath)
	case "darwin":
		cmd = exec.Command("imagesnap", "-w", "3", imagePath)
	}

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
