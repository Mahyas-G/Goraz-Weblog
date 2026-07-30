package main

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"weblog/internal/auth"
	"weblog/internal/config"
	"weblog/internal/db"
	"weblog/internal/handler"
	"weblog/internal/logger"
	appmw "weblog/internal/middleware"
	"weblog/internal/repository"
	"weblog/internal/service"
)

type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := template.ParseFiles("web/templates/layouts/base.html", "web/templates/"+name)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

func main() {
	log := logger.New()

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		log.Error("config error", "error", err)
		os.Exit(1)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Error("database connection error", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	sessionStore := auth.NewSessionStore(database)
	authService := service.NewAuthService(userRepo, sessionStore)
	authHandler := handler.NewAuthHandler(authService)

	boardRepo := repository.NewBoardRepository(database)
	shareRepo := repository.NewShareRepository(database)
	boardService := service.NewBoardService(boardRepo, shareRepo, userRepo)
	boardHandler := handler.NewBoardHandler(boardService)

	e := echo.New()
	e.Renderer = &TemplateRenderer{}
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(appmw.LoadCurrentUser(authService))

	e.Static("/static", "web/static")

	e.GET("/healthz", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx, database); err != nil {
			return c.String(http.StatusServiceUnavailable, "db unreachable: "+err.Error())
		}
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/signup", authHandler.ShowSignupForm)
	e.POST("/signup", authHandler.Signup)
	e.GET("/login", authHandler.ShowLoginForm)
	e.POST("/login", authHandler.Login)
	e.POST("/logout", authHandler.Logout, appmw.RequireAuth)

	e.GET("/weblog", boardHandler.Feed, appmw.RequireAuth)
	e.GET("/weblog/new", boardHandler.ShowCreateForm, appmw.RequireAuth)
	e.POST("/weblog", boardHandler.Create, appmw.RequireAuth)
	e.GET("/weblog/:id", boardHandler.Detail, appmw.RequireAuth)
	e.POST("/weblog/:id/delete", boardHandler.Delete, appmw.RequireAuth)

	e.GET("/", func(c echo.Context) error {
		if appmw.CurrentUser(c) != nil {
			return c.Redirect(http.StatusSeeOther, "/weblog")
		}
		return c.Redirect(http.StatusSeeOther, "/login")
	})

	log.Info("server starting", "port", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}
