package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func regenalpha(pages int) {
	rect := image.Rect(0, 0, 1920, 1080)
	oimg := createRandomImage(rect)
	fmt.Println("Transparent")
	for i := 1; i <= pages; i++ {
		fmt.Printf("%d/%d\n", i, pages)
		of := fmt.Sprintf("tmp/F%04d.png", i)
		wf := fmt.Sprintf("tmp/W%04d.png", i)
		bf := fmt.Sprintf("tmp/B%04d.png", i)
		wimg := load(wf)
		bimg := load(bf)
		gen(wimg, bimg, oimg)
		save(of, oimg)
	}
}

func alpha(WR uint8, WG uint8, WB uint8, BR uint8, BG uint8, BB uint8) uint8 {
	var bright int
	bright = int(WR) + int(WG) + int(WB) - int(BR) - int(BG) - int(BB)
	if bright < 0 {
		bright = 0
	}
	if bright > 765 {
		bright = 765
	}
	return (255 - uint8(bright/3))

}

func gen(white *image.NRGBA, black *image.NRGBA, out *image.NRGBA) {
	ptr := 0
	for x := 0; x < 1920; x++ {
		for y := 0; y < 1080; y++ {
			out.Pix[ptr] = black.Pix[ptr]
			out.Pix[ptr+1] = black.Pix[ptr+1]
			out.Pix[ptr+2] = black.Pix[ptr+2]

			out.Pix[ptr+3] = alpha(white.Pix[ptr], white.Pix[ptr+1], white.Pix[ptr+2], black.Pix[ptr], black.Pix[ptr+1], black.Pix[ptr+2])

			ptr += 4
		}
	}
}

func createRandomImage(rect image.Rectangle) (created *image.NRGBA) {
	pix := make([]uint8, rect.Dx()*rect.Dy()*4)
	created = &image.NRGBA{
		Pix:    pix,
		Stride: rect.Dx() * 4,
		Rect:   rect,
	}
	return
}

func save(filePath string, img *image.NRGBA) {
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	png.Encode(imgFile, img.SubImage(img.Rect))
}

func load(filePath string) *image.NRGBA {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot read file:", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Cannot decode file:", err)
	}
	return img.(*image.NRGBA)
}
