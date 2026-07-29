package service

import (
	"errors"

	"weblog/internal/model"
	"weblog/internal/repository"
	"weblog/internal/validation"
)

var ErrForbidden = errors.New("you do not have permission to perform this action")

type BoardService struct {
	boards *repository.BoardRepository
}

func NewBoardService(boards *repository.BoardRepository) *BoardService {
	return &BoardService{boards: boards}
}

func (s *BoardService) Create(title, content string, imagePath *string, authorID int, privacy string) (*model.Board, error) {
	if err := validation.ValidateTitle(title); err != nil {
		return nil, err
	}
	if err := validation.ValidateContent(content); err != nil {
		return nil, err
	}
	if err := validation.ValidatePrivacy(privacy); err != nil {
		return nil, err
	}

	return s.boards.Create(title, content, imagePath, authorID, privacy)
}

func (s *BoardService) Feed(userID int) ([]model.Board, error) {
	return s.boards.ListVisibleToUser(userID)
}

func (s *BoardService) Detail(boardID, userID int) (*model.Board, error) {
	return s.boards.FindVisibleToUser(boardID, userID)
}

func (s *BoardService) Delete(boardID, userID int) error {
	board, err := s.boards.FindVisibleToUser(boardID, userID)
	if err != nil {
		return err
	}
	if board.AuthorID != userID {
		return ErrForbidden
	}
	return s.boards.Delete(boardID)
}
