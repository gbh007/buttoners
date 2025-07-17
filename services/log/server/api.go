package server

import (
	"net/http"
	"time"

	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Activity(c *fiber.Ctx) error {
	var req logclient.ActivityRequest

	err := c.BodyParser(&req)
	if err != nil {
		c.Set(fiber.HeaderContentType, logclient.ContentType)

		return c.
			Status(http.StatusBadRequest).
			JSON(logclient.ErrorResponse{
				Code:    "parse",
				Details: err.Error(),
			})
	}

	count, last, err := s.db.SelectCompressedUserLogByUserID(c.UserContext(), req.UserID)
	if err != nil {
		c.Set(fiber.HeaderContentType, logclient.ContentType)

		return c.
			Status(http.StatusInternalServerError).
			JSON(logclient.ErrorResponse{
				Code:    "logic",
				Details: err.Error(),
			})
	}

	c.Set(fiber.HeaderContentType, logclient.ContentType)

	return c.
		Status(http.StatusOK).
		JSON(logclient.ActivityResponse{
			RequestCount: count,
			LastRequest:  last,
		})
}

func (s *Server) authHandler(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("Authorization"))

	if token == "" {
		ctx.Set(fiber.HeaderContentType, logclient.ContentType)

		return ctx.
			Status(http.StatusUnauthorized).
			JSON(logclient.ErrorResponse{
				Code:    "unauthorized",
				Details: "empty token",
			})
	}

	if token != s.cfg.SelfToken {
		ctx.Set(fiber.HeaderContentType, logclient.ContentType)

		return ctx.
			Status(http.StatusForbidden).
			JSON(logclient.ErrorResponse{
				Code:    "forbidden",
				Details: "invalid token",
			})
	}

	return ctx.Next()
}
func (s *Server) observabilityHandler(ctx *fiber.Ctx) error {
	tStart := time.Now()

	s.httpServerMetrics.IncActive(string(ctx.Request().Host()), string(ctx.Request().URI().Path()), string(ctx.Request().Header.Method()))
	defer s.httpServerMetrics.DecActive(string(ctx.Request().Host()), string(ctx.Request().URI().Path()), string(ctx.Request().Header.Method()))

	defer func() {
		s.httpServerMetrics.AddHandle(string(ctx.Request().Host()), string(ctx.Request().URI().Path()), string(ctx.Request().Header.Method()), ctx.Response().StatusCode(), time.Since(tStart))
	}()

	defer observability.LogFastHTTPData(ctx.UserContext(), s.l, "log server request", ctx.Request(), ctx.Response())

	return ctx.Next()
}
