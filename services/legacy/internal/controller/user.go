package controller

import (
	"errors"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

const userSessionCookieName = "baton-session"

func (c Controller) GetUser(ctx *fasthttp.RequestCtx) {
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

	ctx.SetStatusCode(http.StatusOK)
	_ = jsoniter.NewEncoder(ctx).Encode(user)
}

func (c Controller) CreateUser(ctx *fasthttp.RequestCtx) {
	user, err := c.userSevice.CreateUser(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		_ = jsoniter.NewEncoder(ctx).Encode(&errorModel{
			Message: err.Error(),
		})

		return
	}

	cookie := fasthttp.Cookie{}
	cookie.SetHTTPOnly(true)
	cookie.SetKey(userSessionCookieName)
	cookie.SetPath("/")
	cookie.SetValue(user.Token)
	ctx.Response.Header.SetCookie(&cookie)

	ctx.SetStatusCode(http.StatusNoContent)
}
