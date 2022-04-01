package sfm

import (
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"
)

// GetImage - reads image from
// path  - to image file relative manager root directory
func (s *SiteFileManager) GetImage(path string, w uint, h uint, thumbnail bool) (*image.Image, string, error) {
	if !s.Drive.Exists(path) {
		return nil, "", fmt.Errorf("file does not exist: %s", s.RootPath+path)
	}
	file, err := s.Drive.Read(path)
	if err != nil {
		return nil, "", fmt.Errorf("can not read file: %s, %v", s.RootPath+path, err)
	}

	var result image.Image

	format := strings.ToLower(filepath.Ext(path))[1:]
	switch format {
	case "png":
		result, err = png.Decode(file)
		break
	case "jpg", "jpeg":
		result, err = jpeg.Decode(file)
		break
	case "webp":
		result, err = webp.Decode(file)
		format = "png"
		break
	default:
		result, format, err = image.Decode(file)
		break
	}
	if err != nil {
		return nil, "", fmt.Errorf("decodin %s: %v", format, err)
	}
	if w != 0 || h != 0 {
		return imageResize(&result, w, h, thumbnail), format, nil
	}
	return &result, format, nil
}

//imageResize
// if with or height equal to 0 image aspect saved original
// if thumbnail true image resized with reduced quality to get the smallest file

func imageResize(image *image.Image, width uint, height uint, thumbnail bool) *image.Image {
	if thumbnail {
		result := resize.Thumbnail(width, height, *image, resize.Lanczos3)
		return &result
	}
	result := resize.Resize(width, height, *image, resize.Lanczos3)
	return &result
}
