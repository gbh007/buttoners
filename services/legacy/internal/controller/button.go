package controller

import (
	"errors"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

func (c *Controller) PressButton(ctx *fasthttp.RequestCtx) {
	token := string(ctx.Request.Header.Cookie(userSessionCookieName))

	user, err := c.userSevice.GetUser(ctx, token)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.SetStatusCode(http.StatusUnauthorized)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: "user not found",
		})

		return
	}

	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	b, err := c.buttonService.PressButton(ctx, user)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	ctx.SetStatusCode(http.StatusOK)
	_ = jsoniter.NewEncoder(ctx).Encode(b)
}

func (c *Controller) Buttons(ctx *fasthttp.RequestCtx) {
	token := string(ctx.Request.Header.Cookie(userSessionCookieName))

	user, err := c.userSevice.GetUser(ctx, token)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.SetStatusCode(http.StatusUnauthorized)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: "user not found",
		})

		return
	}

	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	buttons, err := c.buttonService.Buttons(ctx, user)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	ctx.SetStatusCode(http.StatusOK)
	_ = jsoniter.NewEncoder(ctx).Encode(buttons)
}

func (c *Controller) ButtonPower(ctx *fasthttp.RequestCtx) {
	id, err := ctx.QueryArgs().GetUint("user")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	err = c.buttonService.ButtonBadge(ctx, ctx, id)
	if err != nil {
		ctx.ResetBody()
		ctx.SetStatusCode(http.StatusInternalServerError)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType("image/svg+xml")
}
