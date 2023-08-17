package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryunosuke121/muscle-SNS/s3client"
)

func main() {

	// AWS S3の設定
	s3client.InitS3Client()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
