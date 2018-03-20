package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/hexis-hr/goImage/api"
	"github.com/hexis-hr/goImage/pic"
	"github.com/spf13/pflag"
)

func testTransform() {

	var imgName, subtype string

	pflag.StringVarP(
		&imgName,
		"name",
		"n",
		"",
		"name of image",
	)
	pflag.StringVarP(
		&subtype,
		"subtype",
		"s",
		"",
		"subtype",
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
		Dimensions: []api.Dimension{{Width: 500, Height: 500}},
		Quality:    100,
		Subtype:    subtype,
		Lib:        "imaging",
		Filter:     "l",
	}

	_, err = pic.Transform(img, p)
	if err != nil {
		log.Fatalln(err)
	}
}

func testWatermark() {

	imgb, _ := os.Open("x2.jpg")
	img, _ := jpeg.Decode(imgb)
	defer imgb.Close()

	wmb, _ := os.Open("watermark.png")
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()

	offset := image.Pt(20, 20)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	imgw, _ := os.Create("watermarked1.jpg")
	jpeg.Encode(imgw, m, &jpeg.Options{jpeg.DefaultQuality})
	defer imgw.Close()
}

func testBackground() {

	pngImgFile, err := os.Open("watermark.png")

	if err != nil {
		fmt.Println("PNG-file.png file not found!")
		os.Exit(1)
	}

	defer pngImgFile.Close()

	// create image from PNG file
	imgSrc, err := png.Decode(pngImgFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create a new Image with the same dimension of PNG image
	newImg := image.NewRGBA(imgSrc.Bounds())

	// we will use white background to replace PNG's transparent background
	// you can change it to whichever color you want with
	// a new color.RGBA{} and use image.NewUniform(color.RGBA{<fill in color>}) function

	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// paste PNG image OVER to newImage
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

	// create new out JPEG file
	jpgImgFile, err := os.Create("JPEG-file.jpg")

	if err != nil {
		fmt.Println("Cannot create JPEG-file.jpg !")
		fmt.Println(err)
		os.Exit(1)
	}

	defer jpgImgFile.Close()

	var opt jpeg.Options
	opt.Quality = 80

	// convert newImage to JPEG encoded byte and save to jpgImgFile
	// with quality = 80
	err = jpeg.Encode(jpgImgFile, newImg, &opt)

	//err = jpeg.Encode(jpgImgFile, newImg, nil) -- use nil if ignore quality options

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Converted PNG file to JPEG file")
}

func testBackground2() {

	//foreGroundColor := image.NewUniform(color.Black)
	backGroundColor := image.Transparent
	backgroundWidth := 700
	backgroundHeight := 50
	background := image.NewRGBA(image.Rect(0, 0, backgroundWidth, backgroundHeight))

	draw.Draw(background, background.Bounds(), backGroundColor, image.ZP, draw.Src)
	// draw something with foreGroundColor.....

	// Save that background image to PNG.
	imgFile, err := os.Create("background.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer imgFile.Close()
	buff := bufio.NewWriter(imgFile)
	err = png.Encode(buff, background)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = buff.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Save to background.png")
}

func main() {

	testTransform()
}
