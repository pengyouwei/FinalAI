package router

import (
	"context"
	"finalai/internal/handler"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.Recover())
	// e.Use(middleware.RequestLogger())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogStatus:   true,
		LogURI:      true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			}

			if v.Error == nil {
				e.Logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST", attrs...)
			} else {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
				e.Logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST", attrs...)
			}
			return nil
		},
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	// 测试路由
	e.GET("/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	v1 := e.Group("/api/v1")

	registerUserRoutes(v1.Group("/user"))
	registerChatRoutes(v1.Group("/chat"))
	registerImageRoutes(v1.Group("/image"))
}

func registerUserRoutes(g *echo.Group) {
	userHandler := handler.NewUserHandler()
	g.POST("/register", userHandler.Register)
	g.POST("/login", userHandler.Login)
}

func registerChatRoutes(g *echo.Group) {
}

func registerImageRoutes(g *echo.Group) {
}
