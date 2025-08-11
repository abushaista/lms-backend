package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const HeaderCorrelationID = "X-Correlation-ID"

func CorrelationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// try header first
		corr := c.Request().Header.Get(HeaderCorrelationID)
		if corr == "" {
			corr = uuid.NewString()
		}

		// store in echo context so other middlewares/handlers can read it
		c.Set("correlation_id", corr)

		// expose to client
		c.Response().Header().Set(HeaderCorrelationID, corr)

		return next(c)
	}
}
