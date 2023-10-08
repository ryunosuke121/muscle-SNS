package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type decodedTokenKey struct{}

func GetDecodedToken(c context.Context) *auth.Token {
	token, ok := c.Value(decodedTokenKey{}).(*auth.Token)
	if !ok {
		return nil
	}
	return token
}

type AuthMiddleware struct {
	AuthClient *auth.Client
}

func NewAuthClient() (*AuthMiddleware, error) {
	opt := option.WithCredentialsFile("/Users/itouryuunosuke/Project/go/muscle-SNS/eco-lane-398113-firebase-adminsdk-j03hg-9d416d1188.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		AuthClient: authClient,
	}, err
}

func (am *AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenHeader := c.Request().Header.Get("Authorization")
		splitToken := strings.Split(tokenHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			return c.String(http.StatusUnauthorized, "invalid token")
		}
		accessToken := splitToken[1]
		if am.AuthClient == nil {
			log.Println("authClient is nil")
		}
		decodedToken, err := am.AuthClient.VerifyIDToken(c.Request().Context(), accessToken)
		if err != nil {
			return c.String(http.StatusUnauthorized, "invalid token")
		}
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), decodedTokenKey{}, decodedToken)))

		return next(c)
	}
}
