package controller

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cnt Controller) pressButton(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(userSessionCookieName)
	if err != nil {
		return err
	}

	user, err := cnt.userSevice.GetUser(ctx, cookie.Value)
	if err != nil {
		return err
	}

	b, err := cnt.buttonService.PressButton(ctx, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (cnt Controller) buttons(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(userSessionCookieName)
	if err != nil {
		return err
	}

	user, err := cnt.userSevice.GetUser(ctx, cookie.Value)
	if err != nil {
		return err
	}

	b, err := cnt.buttonService.Buttons(ctx, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func (cnt Controller) buttonPower(c echo.Context) error {
	ctx := c.Request().Context()

	var req buttonPowerRequest

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)

	err = cnt.buttonService.ButtonBadge(ctx, buff, req.User)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/svg+xml", buff.Bytes())
}
