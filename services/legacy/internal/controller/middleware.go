package controller

import (
	"time"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/services/legacy/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
)

func (cnt *Controller) logActivity() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := time.Now()
			key := ksuid.New().String()

			var user domain.User

			cookie, _ := c.Cookie(userSessionCookieName)
			if cookie != nil {
				user, _ = cnt.userSevice.GetUser(c.Request().Context(), cookie.Value)
			}

			err := next(c)

			data := dto.KafkaLogData{
				Addr:         c.Request().RemoteAddr,
				UserID:       int64(user.ID),
				SessionToken: user.Token,
				Action:       c.Path(),
				RequestTime:  t,
				RealIP:       c.RealIP(),
				ForwardedFor: c.Request().Header.Values("X-Forwarded-For"),
			}

			if err != nil {
				data.ErrorText = err.Error()
			}

			_ = cnt.kafkaLogClient.Write(c.Request().Context(), key, data)

			return err
		}
	}
}
