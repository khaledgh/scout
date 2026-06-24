package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func SaveUpload(file *multipart.FileHeader, uploadDir string, allowedMimes string, maxSizeMB int64) (string, error) {
	if file.Size > maxSizeMB*1024*1024 {
		return "", fmt.Errorf("file size exceeds %dMB limit", maxSizeMB)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := make([]byte, 512)
	if _, err = src.Read(buf); err != nil {
		return "", err
	}
	if _, err = src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	randomName, err := randomHex(16)
	if err != nil {
		return "", err
	}
	filename := randomName + ext

	if err = os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
