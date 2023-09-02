package main

import (
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ryunosuke121/muscle-SNS/controller"
	"github.com/ryunosuke121/muscle-SNS/db"
	"github.com/ryunosuke121/muscle-SNS/repository"
	"github.com/ryunosuke121/muscle-SNS/router"
	"github.com/ryunosuke121/muscle-SNS/s3client"
	"github.com/ryunosuke121/muscle-SNS/usecase"
	"github.com/ryunosuke121/muscle-SNS/validator"
)

func main() {
	db := db.NewDB()

	// AWS S3の設定
	client := s3client.NewS3Client()
	presignS3client := s3client.NewPresignS3Client(client)
	userValidator := validator.NewUserValidator()
	userRepository := repository.NewUserRepository(db, client, presignS3client)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	userController := controller.NewUserController(userUsecase)
	trainingController := controller.NewTrainingController(client, presignS3client)
	groupController := controller.NewGroupController(client, presignS3client)
	e := router.NewRouter(userController, *trainingController, *groupController)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
