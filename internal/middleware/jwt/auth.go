package jwtauth

import (
	myjwt "finalai/pkg/jwt"
	"net/http"
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
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code": -1,
					"msg":  "缺少 Authorization 头",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code": -1,
					"msg":  "Authorization 格式错误，应为 Bearer <token>",
				})
			}

			tokenStr := strings.TrimSpace(parts[1])
			claims, err := myjwt.ParseToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code": -1,
					"msg":  "无效或过期的 token",
				})
			}

			c.Set(ContextKeyUsername, claims.Username)
			c.Set(ContextKeyClaims, claims)

			return next(c)
		}
	}
}
