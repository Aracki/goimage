package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/hexis-hr/goImage/api"
	"github.com/hexis-hr/goImage/pic"
	"github.com/spf13/pflag"
)

func main() {

	var imgName string

	pflag.StringVarP(
		&imgName,
		"name",
		"n",
		"",
		"name of image",
	)
	pflag.Parse()

	f, err := os.Open(imgName)
	if err != nil {
		log.Fatalln(err)
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}

	p := api.Params{
		ImgName:    imgName,
		Dimensions: []api.Dimension{{Width: 300, Height: 300}},
		Quality:    1,
		Subtype:    "crop",
	}

	_, err = pic.Transform(img, p)
	if err != nil {
		log.Fatalln(err)
	}
}
