package utils

import (
	"io"
	"mime/multipart"
	"net/http"
)

func InspectFileMimeType(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(fileBytes)
	return contentType, nil
}

func GetFileExtension(mime_type string) string {
	switch mime_type {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}
