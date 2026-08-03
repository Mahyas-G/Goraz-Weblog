package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"weblog/internal/middleware"
	"weblog/internal/repository"
	"weblog/internal/service"
	"weblog/internal/validation"
)

type CommentHandler struct {
	commentService *service.CommentService
	boardService   *service.BoardService
}

func NewCommentHandler(commentService *service.CommentService, boardService *service.BoardService) *CommentHandler {
	return &CommentHandler{commentService: commentService, boardService: boardService}
}

func (h *CommentHandler) Create(c echo.Context) error {
	user := middleware.CurrentUser(c)

	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid board id")
	}

	content := c.FormValue("content")

	if _, err := h.commentService.Add(boardID, user.ID, content); err != nil {
		switch {
		case errors.Is(err, validation.ErrCommentEmpty), errors.Is(err, validation.ErrCommentTooLong):
			return h.renderWithCommentError(c, boardID, user.ID, content, err.Error())
		case errors.Is(err, repository.ErrBoardNotFound):
			return c.String(http.StatusNotFound, "board not found")
		default:
			return c.String(http.StatusInternalServerError, "failed to add comment")
		}
	}

	return c.Redirect(http.StatusSeeOther, "/weblog/"+strconv.Itoa(boardID))
}

func (h *CommentHandler) renderWithCommentError(c echo.Context, boardID, userID int, commentContent, commentError string) error {
	data, err := buildBoardDetailData(h.boardService, h.commentService, boardID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			return c.String(http.StatusNotFound, "board not found")
		}
		return c.String(http.StatusInternalServerError, "failed to load board")
	}

	data["CommentError"] = commentError
	data["CommentContent"] = commentContent

	return c.Render(http.StatusUnprocessableEntity, "weblog/detail.html", data)
}
