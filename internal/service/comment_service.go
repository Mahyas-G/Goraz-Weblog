package service

import (
	"strings"

	"weblog/internal/model"
	"weblog/internal/repository"
	"weblog/internal/validation"
)

type CommentService struct {
	comments *repository.CommentRepository
	boards   *repository.BoardRepository
}

func NewCommentService(comments *repository.CommentRepository, boards *repository.BoardRepository) *CommentService {
	return &CommentService{comments: comments, boards: boards}
}

func (s *CommentService) Add(boardID, authorID int, content string) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if err := validation.ValidateComment(content); err != nil {
		return nil, err
	}

	if _, err := s.boards.FindVisibleToUser(boardID, authorID); err != nil {
		return nil, err
	}

	return s.comments.Create(boardID, authorID, content)
}

func (s *CommentService) List(boardID int) ([]model.Comment, error) {
	return s.comments.ListByBoard(boardID)
}
