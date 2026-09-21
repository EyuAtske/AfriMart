package product

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

const (
	maxProductImageSize  = 5 << 20
	maxProductImageCount = 10
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func validateProductImage(file *multipart.FileHeader) error {
	if file.Size > maxProductImageSize {
		return fmt.Errorf("image %q exceeds the 5 MB size limit", file.Filename)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("could not open image %q: %w", file.Filename, err)
	}
	defer src.Close()

	buffer := make([]byte, 512)

	n, err := src.Read(buffer)
	if err != nil {
		return fmt.Errorf("could not read image %q: %w", file.Filename, err)
	}

	contentType := http.DetectContentType(buffer[:n])

	if !allowedImageTypes[contentType] {
		return fmt.Errorf(
			"unsupported image type %q for file %q",
			contentType,
			file.Filename,
		)
	}

	return nil
}

func imageExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
