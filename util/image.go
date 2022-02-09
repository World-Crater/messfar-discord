package util

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"regexp"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

var (
	isImageReg *regexp.Regexp
)

func Init() error {
	reg, err := regexp.Compile(`(\.png|\.jpeg|\.jpg)$`)
	if err != nil {
		return errors.Wrap(err, "create regexp failed")
	}
	isImageReg = reg
	return nil
}

func IsImage(url string) bool {
	return isImageReg.MatchString(url)
}

func DownloadImage(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request fail")
	}
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request fail")
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		return nil, errors.Wrap(err, "decode file to image failed")
	}

	return buf.Bytes(), nil
}

func GetSize(imageBuffer *bytes.Buffer) (uint, uint, error) {
	m, _, err := image.Decode(imageBuffer)
	if err != nil {
		return 0, 0, err
	}
	g := m.Bounds()
	height := g.Dy()
	width := g.Dx()
	return uint(width), uint(height), nil
}

func ImageResizeByBuffer(file *bytes.Buffer, width uint) (*bytes.Buffer, error) {
	imageTypeString := http.DetectContentType(file.Bytes())

	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	img = imaging.Resize(img, int(width), 0, imaging.Lanczos)
	imageBuffer := bytes.NewBuffer(nil)

	fmt.Println("york", imageTypeString)

	switch imageTypeString {
	case "image/jpeg":
		if err := jpeg.Encode(imageBuffer, img, nil); err != nil {
			return nil, errors.Wrap(err, "encode image failed")
		}
	case "image/png":
		if err := png.Encode(imageBuffer, img); err != nil {
			return nil, errors.Wrap(err, "encode image failed")
		}
	}

	return imageBuffer, nil
}

//yorktodo
// func ImageProcessing(buffer []byte, quality int) ([]byte, error) {
// 	imageTypeString := http.DetectContentType(buffer)
// 	var imageType bimg.ImageType

// 	switch imageTypeString {
// 	case "image/jpeg":
// 		imageType = bimg.JPEG
// 	case "image/png":
// 		imageType = bimg.PNG
// 	}

// 	converted, err := bimg.NewImage(buffer).Convert(imageType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return processed, nil
// }
