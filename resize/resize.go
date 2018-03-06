package resize

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/aracki/gohexis/gohexis/api"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

type ImageFile struct {
	Image    image.Image
	FileName string
	FullPath string
}

func Resize(imgFile *ImageFile, dims []api.Dimension) ([]string, error) {

	var files []string

	for _, d := range dims {

		folder := "/tmp/" +
			strconv.Itoa(d.W) + "x" +
			strconv.Itoa(d.H) + "/"

		_ = os.MkdirAll(folder, os.ModePerm)
		imgFile.FullPath = folder + imgFile.FileName

		dstImage128 := imaging.Resize(imgFile.Image, d.W, d.H, imaging.Lanczos)

		toImg, err := os.Create(imgFile.FullPath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Could not create new file"))
		}
		defer toImg.Close()

		if err := jpeg.Encode(toImg, dstImage128, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Could not encode image"))
		} else {
			fmt.Printf("%s resized and saved\n", imgFile.FileName)
			files = append(files, imgFile.FullPath)
		}
	}

	return files, nil
}
