package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
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

var templatePages = []string{
	"auth/signup.html",
	"auth/login.html",
	"weblog/feed.html",
	"weblog/create.html",
	"weblog/detail.html",
}

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func NewTemplateRenderer() (*TemplateRenderer, error) {
	templates := make(map[string]*template.Template)

	for _, page := range templatePages {
		tmpl, err := template.ParseFiles(
			"web/templates/layouts/base.html",
			"web/templates/partials/navbar.html",
			"web/templates/"+page,
		)
		if err != nil {
			return nil, fmt.Errorf("parsing template %q: %w", page, err)
		}
		templates[page] = tmpl
	}

	return &TemplateRenderer{templates: templates}, nil
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template %q not registered", name)
	}

	if m, ok := data.(map[string]any); ok {
		m["CurrentUser"] = appmw.CurrentUser(c)
		if token, ok := c.Get("csrf").(string); ok {
			m["CSRFToken"] = token
		}
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

	renderer, err := NewTemplateRenderer()
	if err != nil {
		log.Error("template error", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(database)
	sessionStore := auth.NewSessionStore(database)
	authService := service.NewAuthService(userRepo, sessionStore)
	authHandler := handler.NewAuthHandler(authService, cfg.CookieSecure)

	boardRepo := repository.NewBoardRepository(database)
	shareRepo := repository.NewShareRepository(database)
	boardService := service.NewBoardService(boardRepo, shareRepo, userRepo, log)

	commentRepo := repository.NewCommentRepository(database)
	commentService := service.NewCommentService(commentRepo, boardRepo)
	commentHandler := handler.NewCommentHandler(commentService, boardService)

	boardHandler := handler.NewBoardHandler(boardService, commentService)

	go runSessionCleanup(sessionStore, log)

	e := echo.New()
	e.Renderer = renderer
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.BodyLimit("10M"))
	e.Use(echomw.CSRFWithConfig(echomw.CSRFConfig{
		TokenLookup: "form:csrf_token",
	}))
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
	e.POST("/weblog/:id/comments", commentHandler.Create, appmw.RequireAuth)

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

func runSessionCleanup(sessionStore *auth.SessionStore, log *slog.Logger) {
	cleanup := func() {
		deleted, err := sessionStore.DeleteExpired()
		if err != nil {
			log.Error("session cleanup failed", "error", err)
			return
		}
		if deleted > 0 {
			log.Info("session cleanup completed", "deleted", deleted)
		}
	}

	cleanup()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		cleanup()
	}
}
