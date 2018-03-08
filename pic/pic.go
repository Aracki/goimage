package pic

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/hexis-hr/goImage/api"
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

// CreateSpecificDimension creates folder based on dimension. (200x200, 450x400...)
// Creates file into appropriate folder.
// DstImg is resized original stored in-memory. It is resized with a given library/filter.
// That image is encoded into previously created file.
// FullPath must start with /tmp/ because of lambda write-to-file rule.
func createSpecificDimension(img image.Image, width, height int, imgName, lib, filter string) (string, error) {

	// create proper folder
	folder := fmt.Sprintf("/tmp/%sx%s/", strconv.Itoa(width), strconv.Itoa(height))
	fullPath := folder + imgName
	_ = os.MkdirAll(folder, os.ModePerm)

	// create file
	f, err := os.Create(fullPath)
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
	if err := jpeg.Encode(f, dstImg, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		return "", errors.Wrap(err, "could not encode image")
	}

	return fullPath, nil
}

// Resize function take Image interface, image name, array of dimensions, specific library we want to use and filter.
// According to that array it will resize each image sequentially.
// Resized images are saved into proper folders.
// Returns list of paths where are saved images.
func Resize(img image.Image, imgName string, dims []api.Dimension, lib, filter string) ([]string, error) {

	var paths []string

	for _, d := range dims {

		fullPath, err := createSpecificDimension(img, d.Width, d.Height, imgName, lib, filter)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s resized and saved\n", imgName)
		paths = append(paths, fullPath)
	}

	return paths, nil
}
