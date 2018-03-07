package pic

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
	case "l":
		return imaging.Lanczos
	default:
		return imaging.Box
	}
}

func resizeImage(img image.Image, width, height int, alg, filter string) image.Image {

	switch alg {
	case "imaging":
		return imaging.Resize(img, width, height, imagingFilter(filter))
	}
	return nil
}

// CreateSpecificDimension creates folder based on dimension. (200x200, 450x400...)
// Creates file into appropriate folder.
// DstImg is resized original stored in-memory. It is resized with a given algorithm and filter.
// That image is encoded into previously created file.
// FullPath must start with /tmp/ because of lambda write-to-file rule.
func createSpecificDimension(img image.Image, width, height int, imgName, alg, filter string) (string, error) {

	// create proper folder
	folder := fmt.Sprintf("/tmp/%sx%s/", strconv.Itoa(width), strconv.Itoa(height))
	fullPath := folder + imgName
	_ = os.MkdirAll(folder, os.ModePerm)

	// create file
	f, err := os.Create(fullPath)
	if err != nil {
		return "", errors.Wrap(err, "Could not create new file")
	}

	// ----------------resize
	// do actual
	dstImg := resizeImage(img, width, height, alg, filter)
	// ----------------

	// write resized dstImg to file
	if err := jpeg.Encode(f, dstImg, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		return "", errors.Wrap(err, "Could not encode image")
	}

	return fullPath, nil
}

// Resize function take Image interface, image name, array of dimensions, specific algorithm we want to use and filter of that algorithm.
// According to that array it will resize each image sequentially.
// Resized images are saved into proper folders.
// Returns list of paths where are saved images.
func Resize(img image.Image, imgName string, dims []api.Dimension, alg string, filter string) ([]string, error) {

	var paths []string

	for _, d := range dims {

		fullPath, err := createSpecificDimension(img, d.Width, d.Height, imgName, alg, filter)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s resized and saved\n", imgName)
		paths = append(paths, fullPath)
	}

	return paths, nil
}
