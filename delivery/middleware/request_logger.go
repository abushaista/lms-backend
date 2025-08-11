package middleware

import (
	"github.com/abushaista/lms-backend/delivery/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func AttachRequestLogger(root zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			l := utils.WithRequestLogger(root, c)
			// store logger into context for handlers/repositories to fetch
			c.Set(utils.CtxLoggerKey, l)
			return next(c)
		}
	}
}
