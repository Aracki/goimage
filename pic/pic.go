package pic

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/hexis-hr/goImage/api"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func resizeImage(imgSrc image.Image, width, height int, lib, filter string) (image.Image, error) {

	switch lib {
	case "imaging":
		return imaging.Resize(imgSrc, width, height, imagingFilter(filter)), nil
	case "nfnt":
		return resize.Resize(uint(width), uint(height), imgSrc, nfntFilter(filter)), nil
	}
	return nil, fmt.Errorf("lib not defined")
}

// CreateSpecificDimension creates folder based on dimension. eg. /tmp/200x200/, /tmp/450x400/
// Creates file into appropriate folder.
// DstImg is resized original stored in-memory. It is resized with a given library/filter.
// That image is encoded into previously created file.
// FullImgPath must start with /tmp/ because of lambda write-to-file rule.
func createSpecificDimension(img image.Image, width, height int, imgName, lib, filter string, quality int) (string, error) {

	basePath := "/tmp/"

	// create proper folder
	folder := basePath + strconv.Itoa(width) + "x" + strconv.Itoa(height) + "/"
	fullImgPath := folder + imgName
	_ = os.MkdirAll(folder, os.ModePerm)

	// create file
	f, err := os.Create(fullImgPath)
	if err != nil {
		return "", errors.Wrap(err, "Could not create new file")
	}

	// ----------------
	// do actual resize
	dstImg, err := resizeImage(img, width, height, lib, filter)
	if err != nil {
		return "", errors.Wrap(err, "resize of image failed")
	}
	// ----------------

	// write resized dstImg to file
	if err := jpeg.Encode(f, dstImg, &jpeg.Options{Quality: quality}); err != nil {
		return "", errors.Wrap(err, "could not encode image")
	}

	return fullImgPath, nil
}

// Resize function take Image interface, image name, array of dimensions, specific library we want to use and filter.
// According to that array it will resize each image sequentially.
// Resized images are saved into proper folders.
// Returns list of paths where are saved images.
func regularResize(img image.Image, imgName string, dims []api.Dimension, lib, filter string, quality int) (paths []string, err error) {

	for _, d := range dims {

		fullPath, err := createSpecificDimension(img, d.Width, d.Height, imgName, lib, filter, quality)
		if err != nil {
			return nil, err
		}
		fmt.Printf("regular resize executed; %s saved under %s\n", imgName, fullPath)
		paths = append(paths, fullPath)
	}

	return paths, nil
}

func smartCrop(img image.Image, imgName string, dims []api.Dimension, quality int) (paths []string, err error) {

	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	basePath := "/tmp/"

	for _, d := range dims {
		topCrop, err := analyzer.FindBestCrop(img, d.Width, d.Height)
		if err != nil {
			return nil, err
		}

		// create proper folder
		folder := basePath + strconv.Itoa(d.Width) + "x" + strconv.Itoa(d.Height) + "/"
		fullImgPath := folder + imgName
		_ = os.MkdirAll(folder, os.ModePerm)

		// create file
		f, err := os.Create(fullImgPath)
		if err != nil {
			return nil, err
		}

		// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
		croppedImg := img.(SubImager).SubImage(topCrop)
		if err := jpeg.Encode(f, croppedImg, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}

		log.Printf("smart crop[%+v] executed; %s saved under %s\n", topCrop, imgName, fullImgPath)
		paths = append(paths, fullImgPath)
	}
	return paths, nil
}

func Transform(img image.Image, p api.Params) ([]string, error) {

	switch p.Subtype {
	case "resize":
		return regularResize(img, p.ImgName, p.Dimensions, p.Lib, p.Filter, p.Quality)
	case "smart_crop":
		return smartCrop(img, p.ImgName, p.Dimensions, p.Quality)
	default:
		return nil, fmt.Errorf("%s subtype doesn't exist", p.Subtype)
	}
}
