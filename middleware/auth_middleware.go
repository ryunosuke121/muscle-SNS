package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authClient *auth.Client
}

func NewAuthClient() *AuthMiddleware {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil
	}

	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (am *AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenHeader := c.Request().Header.Get("Authorization")
		splitToken := strings.Split(tokenHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			return c.String(http.StatusUnauthorized, "invalid token")
		}
		accessToken := splitToken[1]
		decodedToken, err := am.authClient.VerifyIDToken(c.Request().Context(), accessToken)
		if err != nil {
			return c.String(http.StatusUnauthorized, "invalid token")
		}

		log.Println(decodedToken)
		return next(c)
	}
}
