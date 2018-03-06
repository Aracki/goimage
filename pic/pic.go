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

// Creates folder based on dimension. (200x200, 450x400...)
// Creates file into appropriate folder.
// DstImg is resized original stored in-memory. It is resized with given filter.
// That image is encoded into previously created file.
// FullPath must start with /tmp/ because of lambda write-to-file rule.
func resizeImage(img image.Image, w, h int, imgName string, filter imaging.ResampleFilter) (string, error) {

	// create proper folder
	folder := fmt.Sprintf("/tmp/%sx%s/", strconv.Itoa(w), strconv.Itoa(h))
	fullPath := folder + imgName
	_ = os.MkdirAll(folder, os.ModePerm)

	// create file
	f, err := os.Create(fullPath)
	if err != nil {
		return "", errors.Wrap(err, "Could not create new file")
	}

	dstImg := imaging.Resize(img, w, h, filter)

	// write resized dstImg to file
	if err := jpeg.Encode(f, dstImg, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		return "", errors.Wrap(err, "Could not encode image")
	}

	return fullPath, nil
}

// Resize function take Image interface, image name, and array of dimensions.
// According to that array it will resize each image sequentially.
// Resized images are saved into proper folders.
// Returns list of paths where are saved images.
func Resize(img image.Image, imgName string, dims []api.Dimension) ([]string, error) {

	var paths []string

	for _, d := range dims {

		fullPath, err := resizeImage(img, d.W, d.H, imgName, imaging.Lanczos)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s resized and saved\n", imgName)
		paths = append(paths, fullPath)
	}

	return paths, nil
}
