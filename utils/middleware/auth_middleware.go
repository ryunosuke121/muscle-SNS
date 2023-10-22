package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"github.com/ryunosuke121/muscle-SNS/src/domain"
	"google.golang.org/api/option"
)

type decodedTokenKey struct{}

type AuthMiddleware struct {
	AuthClient *auth.Client
}

func NewAuthClient() (*AuthMiddleware, error) {
	opt := option.WithCredentialsFile("/app/eco-lane-398113-firebase-adminsdk-j03hg-9d416d1188.json")
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

func GetDecodedToken(c context.Context) (*auth.Token, error) {
	token, ok := c.Value(decodedTokenKey{}).(*auth.Token)
	if !ok {
		return nil, fmt.Errorf("failed to get decoded token")
	}
	return token, nil
}

func GetUserId(c context.Context) (domain.UserID, error) {
	decodedToken, err := GetDecodedToken(c)
	if err != nil || decodedToken == nil {
		return "", err
	}

	user_id, ok := (*decodedToken).Claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get user_id")
	}

	return domain.UserID(user_id), nil
}

func GetEmail(c context.Context) (string, error) {
	decodedToken, err := GetDecodedToken(c)
	if err != nil || decodedToken == nil {
		return "", err
	}

	email, ok := (*decodedToken).Claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get email")
	}

	return email, nil
}
