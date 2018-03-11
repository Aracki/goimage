package pic

import (
	"image"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func imagingFilter(f string) imaging.ResampleFilter {

	switch f {
	case "nn":
		return imaging.NearestNeighbor
	case "box":
		return imaging.Box
	case "linear":
		return imaging.Linear
	case "mn":
		return imaging.MitchellNetravali
	case "cr":
		return imaging.CatmullRom
	case "gaussian":
		return imaging.Gaussian
	case "lan":
		return imaging.Lanczos
	default:
		return imaging.Box
	}
}

func nfntFilter(f string) resize.InterpolationFunction {

	switch f {
	case "nn":
		return resize.NearestNeighbor
	case "bil":
		return resize.Bilinear
	case "bic":
		return resize.Bicubic
	case "mn":
		return resize.MitchellNetravali
	case "lan2":
		return resize.Lanczos2
	case "lan3":
		return resize.Lanczos3
	default:
		return resize.NearestNeighbor
	}
}

func createFolderAndFile(imgName string, width, height int) (f *os.File, fullImgPath string, err error) {

	basePath := "/tmp/"
	// create proper folder
	folder := basePath + strconv.Itoa(width) + "x" + strconv.Itoa(height) + "/"
	fullImgPath = folder + imgName
	_ = os.MkdirAll(folder, os.ModePerm)

	// create file
	f, err = os.Create(fullImgPath)
	if err != nil {
		return nil, "", err
	}

	return f, fullImgPath, nil
}
