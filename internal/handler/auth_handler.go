package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"weblog/internal/middleware"
	"weblog/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	cookieSecure bool
}

func NewAuthHandler(authService *service.AuthService, cookieSecure bool) *AuthHandler {
	return &AuthHandler{authService: authService, cookieSecure: cookieSecure}
}

func (h *AuthHandler) ShowSignupForm(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/signup.html", map[string]any{"Error": ""})
}

func (h *AuthHandler) Signup(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	_, err := h.authService.Signup(username, password)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, "auth/signup.html", map[string]any{"Error": err.Error()})
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (h *AuthHandler) ShowLoginForm(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/login.html", map[string]any{"Error": ""})
}

func (h *AuthHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	session, err := h.authService.Login(username, password)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, "auth/login.html", map[string]any{"Error": err.Error()})
	}

	cookie := new(http.Cookie)
	cookie.Name = middleware.SessionCookieName
	cookie.Value = session.Token
	cookie.Expires = session.ExpiresAt
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = h.cookieSecure
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/weblog")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie(middleware.SessionCookieName)
	if err == nil {
		_ = h.authService.Logout(cookie.Value)
	}

	expired := new(http.Cookie)
	expired.Name = middleware.SessionCookieName
	expired.Value = ""
	expired.Path = "/"
	expired.Expires = time.Unix(0, 0)
	c.SetCookie(expired)

	return c.Redirect(http.StatusSeeOther, "/login")
}
