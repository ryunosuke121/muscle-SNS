package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/ryunosuke121/muscle-SNS/router"
	"github.com/ryunosuke121/muscle-SNS/s3client"
	"github.com/ryunosuke121/muscle-SNS/src/controller"
	"github.com/ryunosuke121/muscle-SNS/src/db"
	"github.com/ryunosuke121/muscle-SNS/src/repository"
	"github.com/ryunosuke121/muscle-SNS/src/usecase"
	"github.com/ryunosuke121/muscle-SNS/src/validator"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
