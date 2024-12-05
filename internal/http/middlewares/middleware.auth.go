package middlewares

import (
	"backend/internal/constants"
	"backend/internal/http/handlers"
	"backend/pkg/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return handlers.NewErrorResponse(ctx, http.StatusUnauthorized, "token not found")
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			return handlers.NewErrorResponse(ctx, http.StatusBadRequest, "invalid header format")
		}

		if headerParts[0] != "Bearer" {
			return handlers.NewErrorResponse(ctx, http.StatusUnauthorized, "token must content Bearer")
		}

		claims, err := jwt.ParseToken(headerParts[1])
		if err != nil {
			return handlers.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid token")
		}

		ctx.Set(constants.CtxAuthenticatedUserKey, claims)
		return next(ctx)
	}
}
