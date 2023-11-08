package main

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ryunosuke121/muscle-SNS/router"
	"github.com/ryunosuke121/muscle-SNS/src/application"
	"github.com/ryunosuke121/muscle-SNS/src/controller"
	"github.com/ryunosuke121/muscle-SNS/src/repository"
	"github.com/ryunosuke121/muscle-SNS/src/repository/db"
	"github.com/ryunosuke121/muscle-SNS/src/repository/redis"
	"github.com/ryunosuke121/muscle-SNS/src/repository/s3client"
	echoValidator "github.com/ryunosuke121/muscle-SNS/utils/validator"
)

func main() {
	db := db.NewDB()

	// AWS S3の設定
	client := s3client.NewS3Client()
	redisClient := redis.NewRedisClient()
	userRepository := repository.NewUserRepository(db, client)
	userService := application.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	postRepository := repository.NewPostRepository(db, client, redisClient)
	postService := application.NewPostService(postRepository)
	postController := controller.NewPostController(postService)
	rankingService := application.NewRankingService(userRepository, postRepository)
	rankingController := controller.NewRankingController(rankingService)
	e := router.NewRouter(userController, postController, rankingController)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Validator = echoValidator.New(validator.New())

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
