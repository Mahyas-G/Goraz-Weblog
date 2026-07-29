package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"weblog/internal/model"
	"weblog/internal/service"
)

const SessionCookieName = "session_token"
const contextUserKey = "current_user"

func LoadCurrentUser(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(SessionCookieName)
			if err != nil {
				return next(c)
			}

			user, err := authService.CurrentUser(cookie.Value)
			if err != nil {
				return next(c)
			}

			c.Set(contextUserKey, user)
			return next(c)
		}
	}
}

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if CurrentUser(c) == nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return next(c)
	}
}

func CurrentUser(c echo.Context) *model.User {
	user, ok := c.Get(contextUserKey).(*model.User)
	if !ok {
		return nil
	}
	return user
}
