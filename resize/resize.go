package resize

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

type ImageFile struct {
	Image    image.Image
	FileName string
	FullPath string
}

type Dimension struct {
	Width  int
	Height int
}

func Resize(imgFile *ImageFile, variations []Dimension) (*os.File, error) {

	dim := variations[0]

	// TODO proper create folders for multiple variations at same time
	folder := "/tmp/" +
		strconv.Itoa(dim.Width) + "x" +
		strconv.Itoa(dim.Height) + "/"

	_ = os.MkdirAll(folder, os.ModePerm)
	imgFile.FullPath = folder + imgFile.FileName

	dstImage128 := imaging.Resize(imgFile.Image, dim.Width, dim.Height, imaging.Lanczos)

	toImg, err := os.Create(imgFile.FullPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not create new file"))
	}
	defer toImg.Close()

	if err := jpeg.Encode(toImg, dstImage128, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not encode image"))
	} else {
		fmt.Printf("%s resized and saved\n", imgFile.FileName)
	}

	return toImg, nil
}
