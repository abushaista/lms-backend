package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	CtxLoggerKey = "logger"
	CtxCorrKey   = "correlation_id"
	CtxUserKey   = "user" // if you store user in context using middleware.UserContext
)

func WithRequestLogger(root zerolog.Logger, c echo.Context) zerolog.Logger {
	// correlation id
	corr, _ := c.Get(CtxCorrKey).(string)
	if corr == "" {
		corr = c.Request().Header.Get("X-Correlation-ID")
	}

	l := root.With().Str("correlation_id", corr)

	// if user context present, also attach user_id and username
	if u := c.Get(CtxUserKey); u != nil {
		if uc, ok := u.(*UserContext); ok {
			l = l.Str("user_id", uc.UserID).Str("username", uc.Username)
		}
	}

	return l.Logger()
}

type UserContext struct {
	UserID   string
	Username string
}

// GetLogger returns request-scoped logger stored in context, or root logger.
func GetLogger(c echo.Context, root zerolog.Logger) zerolog.Logger {
	if v := c.Get(CtxLoggerKey); v != nil {
		if l, ok := v.(zerolog.Logger); ok {
			return l
		}
	}
	// create one on the fly (not stored)
	return WithRequestLogger(root, c)
}
