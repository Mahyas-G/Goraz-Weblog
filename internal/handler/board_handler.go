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
	"weblog/internal/upload"
	"weblog/internal/validation"
)

type BoardHandler struct {
	boardService   *service.BoardService
	commentService *service.CommentService
}

func NewBoardHandler(boardService *service.BoardService, commentService *service.CommentService) *BoardHandler {
	return &BoardHandler{boardService: boardService, commentService: commentService}
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
	return c.Render(http.StatusOK, "weblog/create.html", map[string]any{
		"Error":           "",
		"Title":           "",
		"Content":         "",
		"Privacy":         "",
		"SharedUsernames": "",
	})
}

func (h *BoardHandler) Create(c echo.Context) error {
	user := middleware.CurrentUser(c)

	title := c.FormValue("title")
	content := c.FormValue("content")
	privacy := c.FormValue("privacy")
	rawSharedUsernames := c.FormValue("shared_usernames")
	sharedUsernames := strings.Split(rawSharedUsernames, ",")

	imageFile, ferr := c.FormFile("image")
	if ferr != nil {
		if !errors.Is(ferr, http.ErrMissingFile) {
			return c.Render(http.StatusInternalServerError, "weblog/create.html", map[string]any{"Error": "failed to read uploaded image"})
		}
		imageFile = nil
	}

	board, err := h.boardService.Create(title, content, imageFile, user.ID, privacy, sharedUsernames)
	if err != nil {
		formData := map[string]any{
			"Title":           title,
			"Content":         content,
			"Privacy":         privacy,
			"SharedUsernames": rawSharedUsernames,
		}
		if isBoardValidationError(err) {
			formData["Error"] = err.Error()
			return c.Render(http.StatusUnprocessableEntity, "weblog/create.html", formData)
		}
		formData["Error"] = "failed to create post"
		return c.Render(http.StatusInternalServerError, "weblog/create.html", formData)
	}

	return c.Redirect(http.StatusSeeOther, "/weblog/"+strconv.Itoa(board.ID))
}

func (h *BoardHandler) Detail(c echo.Context) error {
	user := middleware.CurrentUser(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	data, err := buildBoardDetailData(h.boardService, h.commentService, id, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			return c.String(http.StatusNotFound, "board not found")
		}
		return c.String(http.StatusInternalServerError, "failed to load board")
	}

	return c.Render(http.StatusOK, "weblog/detail.html", data)
}

func buildBoardDetailData(boardService *service.BoardService, commentService *service.CommentService, boardID, userID int) (map[string]any, error) {
	board, err := boardService.Detail(boardID, userID)
	if err != nil {
		return nil, err
	}

	isAuthor := board.AuthorID == userID

	var sharedWith []string
	if isAuthor && board.Privacy == model.PrivacyPrivate {
		sharedWith, err = boardService.SharedUsernames(board.ID)
		if err != nil {
			return nil, err
		}
	}

	comments, err := commentService.List(board.ID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"Board":          board,
		"IsAuthor":       isAuthor,
		"SharedWith":     sharedWith,
		"Comments":       comments,
		"CommentError":   "",
		"CommentContent": "",
	}, nil
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
		errors.Is(err, validation.ErrInvalidPrivacy) ||
		errors.Is(err, upload.ErrFileTooLarge) ||
		errors.Is(err, upload.ErrUnsupportedFileType)
}
