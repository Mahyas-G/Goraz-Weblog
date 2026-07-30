package service

import (
	"errors"
	"fmt"
	"strings"

	"weblog/internal/model"
	"weblog/internal/repository"
	"weblog/internal/validation"
)

var ErrForbidden = errors.New("you do not have permission to perform this action")

type ErrUsernameNotFound struct {
	Username string
}

func (e *ErrUsernameNotFound) Error() string {
	return fmt.Sprintf("user %q not found", e.Username)
}

type BoardService struct {
	boards *repository.BoardRepository
	shares *repository.ShareRepository
	users  *repository.UserRepository
}

func NewBoardService(boards *repository.BoardRepository, shares *repository.ShareRepository, users *repository.UserRepository) *BoardService {
	return &BoardService{boards: boards, shares: shares, users: users}
}

func (s *BoardService) Create(title, content string, imagePath *string, authorID int, privacy string, sharedUsernames []string) (*model.Board, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	privacy = strings.TrimSpace(privacy)

	if err := validation.ValidateTitle(title); err != nil {
		return nil, err
	}
	if err := validation.ValidateContent(content); err != nil {
		return nil, err
	}
	if err := validation.ValidatePrivacy(privacy); err != nil {
		return nil, err
	}

	var sharedUserIDs []int
	if privacy == model.PrivacyPrivate {
		ids, err := s.resolveUsernames(authorID, sharedUsernames)
		if err != nil {
			return nil, err
		}
		sharedUserIDs = ids
	}

	return s.boards.Create(title, content, imagePath, authorID, privacy, sharedUserIDs)
}

func (s *BoardService) Feed(userID int) ([]model.Board, error) {
	return s.boards.ListVisibleToUser(userID)
}

func (s *BoardService) Detail(boardID, userID int) (*model.Board, error) {
	return s.boards.FindVisibleToUser(boardID, userID)
}

func (s *BoardService) SharedUsernames(boardID int) ([]string, error) {
	return s.shares.ListUsernames(boardID)
}

func (s *BoardService) Delete(boardID, userID int) error {
	board, err := s.boards.FindVisibleToUser(boardID, userID)
	if err != nil {
		return err
	}
	if board.AuthorID != userID {
		return ErrForbidden
	}

	deleted, err := s.boards.DeleteOwned(boardID, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return repository.ErrBoardNotFound
	}
	return nil
}

func (s *BoardService) resolveUsernames(authorID int, usernames []string) ([]int, error) {
	seen := map[int]bool{authorID: true}
	var ids []int

	for _, raw := range usernames {
		username := strings.TrimSpace(raw)
		if username == "" {
			continue
		}

		user, err := s.users.FindByUsername(username)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return nil, &ErrUsernameNotFound{Username: username}
			}
			return nil, err
		}

		if !seen[user.ID] {
			seen[user.ID] = true
			ids = append(ids, user.ID)
		}
	}

	return ids, nil
}
