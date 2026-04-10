package jwtauth

import (
	"finalai/internal/common/apperror"
	"finalai/internal/controller/response"
	myjwt "finalai/pkg/jwt"
	"strings"

	"github.com/labstack/echo/v5"
)

const (
	ContextKeyUsername = "username"
	ContextKeyClaims   = "jwt_claims"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, apperror.ErrTokenMissing)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				return response.Error(c, apperror.ErrTokenFormat)
			}

			tokenStr := strings.TrimSpace(parts[1])
			claims, err := myjwt.ParseToken(tokenStr)
			if err != nil {
				return response.Error(c, apperror.ErrTokenInvalid)
			}

			c.Set(ContextKeyUsername, claims.Username)
			c.Set(ContextKeyClaims, claims)

			return next(c)
		}
	}
}
