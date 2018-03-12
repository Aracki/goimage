package pic

import (
	"image"
	"image/jpeg"
	"log"

	"github.com/disintegration/imaging"
	"github.com/hexis-hr/goImage/api"
)

func crop(img image.Image, imgName string, dims []api.Dimension, quality int) (paths []string, err error) {

	for _, d := range dims {

		f, imgFullPath, err := createFolderAndFile(imgName, d.Width, d.Height)
		if err != nil {
			return nil, err
		}

		dstImg := imaging.Fill(img, d.Width, d.Height, imaging.Center, imaging.Lanczos)
		if err := jpeg.Encode(f, dstImg, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}

		log.Printf("crop; quality=%s; %s saved under %s\n", imgName, quality, imgFullPath)
		paths = append(paths, imgFullPath)

		f.Close()
	}

	return paths, nil
}
