package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"weblog/internal/middleware"
	"weblog/internal/model"
	"weblog/internal/repository"
	"weblog/internal/service"
	"weblog/internal/validation"
)

type BoardHandler struct {
	boardService *service.BoardService
}

func NewBoardHandler(boardService *service.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

func (h *BoardHandler) Feed(c echo.Context) error {
	user := middleware.CurrentUser(c)

	boards, err := h.boardService.Feed(user.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to load feed")
	}

	return c.Render(http.StatusOK, "weblog/feed.html", map[string]any{"Boards": boards})
}

func (h *BoardHandler) ShowCreateForm(c echo.Context) error {
	return c.Render(http.StatusOK, "weblog/create.html", map[string]any{"Error": ""})
}

func (h *BoardHandler) Create(c echo.Context) error {
	user := middleware.CurrentUser(c)

	title := c.FormValue("title")
	content := c.FormValue("content")
	privacy := c.FormValue("privacy")
	sharedUsernames := strings.Split(c.FormValue("shared_usernames"), ",")

	board, err := h.boardService.Create(title, content, nil, user.ID, privacy, sharedUsernames)
	if err != nil {
		if isBoardValidationError(err) {
			return c.Render(http.StatusUnprocessableEntity, "weblog/create.html", map[string]any{"Error": err.Error()})
		}
		return c.Render(http.StatusInternalServerError, "weblog/create.html", map[string]any{"Error": "failed to create post"})
	}

	return c.Redirect(http.StatusSeeOther, "/weblog/"+strconv.Itoa(board.ID))
}

func (h *BoardHandler) Detail(c echo.Context) error {
	user := middleware.CurrentUser(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	board, err := h.boardService.Detail(id, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			return c.String(http.StatusNotFound, "board not found")
		}
		return c.String(http.StatusInternalServerError, "failed to load board")
	}

	isAuthor := board.AuthorID == user.ID

	var sharedWith []string
	if isAuthor && board.Privacy == model.PrivacyPrivate {
		sharedWith, err = h.boardService.SharedUsernames(board.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load board")
		}
	}

	return c.Render(http.StatusOK, "weblog/detail.html", map[string]any{
		"Board":      board,
		"IsAuthor":   isAuthor,
		"SharedWith": sharedWith,
	})
}

func (h *BoardHandler) Delete(c echo.Context) error {
	user := middleware.CurrentUser(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	if err := h.boardService.Delete(id, user.ID); err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			return c.String(http.StatusForbidden, "you do not have permission to delete this post")
		case errors.Is(err, repository.ErrBoardNotFound):
			return c.String(http.StatusNotFound, "board not found")
		default:
			return c.String(http.StatusInternalServerError, "failed to delete board")
		}
	}

	return c.Redirect(http.StatusSeeOther, "/weblog")
}

func isBoardValidationError(err error) bool {
	var usernameErr *service.ErrUsernameNotFound
	if errors.As(err, &usernameErr) {
		return true
	}
	return errors.Is(err, validation.ErrTitleEmpty) ||
		errors.Is(err, validation.ErrTitleTooLong) ||
		errors.Is(err, validation.ErrContentEmpty) ||
		errors.Is(err, validation.ErrContentTooLong) ||
		errors.Is(err, validation.ErrInvalidPrivacy)
}
