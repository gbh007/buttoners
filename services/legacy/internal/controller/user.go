package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (cnt Controller) createUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req loginRequest

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	user, err := cnt.userSevice.CreateUser(ctx, req.Login, req.Pass)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     userSessionCookieName,
		Value:    user.Token,
		Path:     "/",
		HttpOnly: true,
	})

	return c.NoContent(http.StatusNoContent)
}

func (cnt Controller) logout(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(userSessionCookieName)
	if err != nil {
		return err
	}

	err = cnt.userSevice.Logout(ctx, cookie.Value)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     userSessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	return c.NoContent(http.StatusNoContent)
}

func (cnt Controller) getUser(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(userSessionCookieName)
	if err != nil {
		return err
	}

	user, err := cnt.userSevice.GetUser(ctx, cookie.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"login":   user.Name,
		"user_id": user.ID,
	})
}
