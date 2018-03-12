package pic

import (
	"image"
	"image/jpeg"
	"log"

	"github.com/hexis-hr/goImage/api"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

func smartCrop(img image.Image, imgName string, dims []api.Dimension, quality int) (paths []string, err error) {

	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())

	for _, d := range dims {
		topCrop, err := analyzer.FindBestCrop(img, d.Width, d.Height)
		if err != nil {
			return nil, err
		}

		f, imgFullPath, err := createFolderAndFile(imgName, d.Width, d.Height)
		if err != nil {
			return nil, err
		}

		// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
		croppedImg := img.(SubImager).SubImage(topCrop)
		if err := jpeg.Encode(f, croppedImg, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}

		log.Printf("smart crop[%+v]; quality=%d; %s saved under %s\n", topCrop, quality, imgName, imgFullPath)
		paths = append(paths, imgFullPath)

		f.Close()
	}
	return paths, nil
}
