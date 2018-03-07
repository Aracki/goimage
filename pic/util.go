package pic

import (
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
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
