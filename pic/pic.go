package pic

import (
	"fmt"
	"image"

	"github.com/hexis-hr/goImage/api"
)

func Transform(img image.Image, p api.Params) ([]string, error) {

	switch p.Subtype {
	case "resize":
		return regularResize(img, p.ImgName, p.Dimensions, p.Lib, p.Filter, p.Quality)
	case "smart_crop":
		return smartCrop(img, p.ImgName, p.Dimensions, p.Quality)
	case "crop":
		return crop(img, p.ImgName, p.Dimensions, p.Quality)
	default:
		return nil, fmt.Errorf("%s subtype doesn't exist", p.Subtype)
	}
}
