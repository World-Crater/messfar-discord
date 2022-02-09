package util

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestImageProcessing(t *testing.T) {
	imageBytes, err := DownloadImage("https://storage.googleapis.com/image.blocktempo.com/2021/05/%E6%88%AA%E5%9C%96-2021-05-31-%E4%B8%8A%E5%8D%8811.23.33.png")
	if err != nil {
		return
	}

	resizeImageBuffer, err := ImageResizeByBuffer(bytes.NewBuffer(imageBytes), 600)
	if err != nil {
		logrus.Errorf("resize image failed. error: %+v", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(resizeImageBuffer.Bytes()))
	if err != nil {
		logrus.Fatalln(err)
	}

	out, _ := os.Create("./img.jpeg")
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100

	err = jpeg.Encode(out, img, &opts)
	//jpeg.Encode(out, img, nil)
	if err != nil {
		logrus.Println(err)
	}
}
