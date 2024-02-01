package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var cfg = oauth2.Config{
	RedirectURL:  fmt.Sprintf("https://%s/api/auth/callback", viper.GetString("domain")),
	ClientID:     viper.GetString("passport.client_id"),
	ClientSecret: viper.GetString("passport.client_secret"),
	Scopes:       []string{"openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:   fmt.Sprintf("%s/auth/o/connect", viper.GetString("passport.endpoint")),
		TokenURL:  fmt.Sprintf("%s/api/auth/token", viper.GetString("passport.endpoint")),
		AuthStyle: oauth2.AuthStyleInParams,
	},
}

func doLogin(c *fiber.Ctx) error {
	url := cfg.AuthCodeURL(uuid.NewString())

	return c.JSON(fiber.Map{
		"target": url,
	})
}

func doPostLogin(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to exchange token: %q", err))
	}

	agent := fiber.
		Get(fmt.Sprintf("%s/api/users/me", viper.GetString("passport.endpoint"))).
		Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.AccessToken))

	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to get userinfo: %q", errs))
	}

	var userinfo services.PassportUserinfo
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse userinfo: %q", err))
	}

	account, err := services.LinkAccount(userinfo)
	access, refresh, err := services.GetToken(account)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to get token: %q", err))
	}

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func doRefreshToken(c *fiber.Ctx) error {
	var data struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	access, refresh, err := services.RefreshToken(data.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to get token: %q", err))
	}

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
