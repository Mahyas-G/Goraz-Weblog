package upload

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"
)

const Dir = "web/static/uploads"

const maxFileSize = 5 * 1024 * 1024

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var ErrFileTooLarge = errors.New("image must be smaller than 5MB")
var ErrUnsupportedFileType = errors.New("image must be jpg, png, or webp")

func Save(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > maxFileSize {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return "", ErrUnsupportedFileType
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if _, _, err := image.DecodeConfig(src); err != nil {
		return "", ErrUnsupportedFileType
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	if err := os.MkdirAll(Dir, 0o755); err != nil {
		return "", err
	}

	filename, err := randomFilename(ext)
	if err != nil {
		return "", err
	}

	dst, err := os.Create(filepath.Join(Dir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return "/static/uploads/" + filename, nil
}

func Delete(urlPath string) error {
	filename := filepath.Base(urlPath)
	return os.Remove(filepath.Join(Dir, filename))
}

func randomFilename(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + ext, nil
}
