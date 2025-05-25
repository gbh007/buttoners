package server

import (
	"net/http"

	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/services/log/internal/pb"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type pbServer struct {
	pb.UnimplementedLogServer

	db *storage.Database
}

func (s *pbServer) Activity(c *fiber.Ctx) error {
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

	count, last, err := s.db.SelectCompressedUserLogByUserID(c.Context(), req.UserID)
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
