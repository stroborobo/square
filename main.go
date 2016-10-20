package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"
)

func usage() {
	fmt.Fprintf(os.Stderr, "square is a tool to make pictures square by adding transparent bars.\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [-o | --override] file [file...]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var override bool
	flag.BoolVarP(&override, "override", "o", false, "Override original file on save")
	flag.Usage = usage
	flag.Parse()

	for _, file := range flag.Args() {
		if err := processFile(file, override); err != nil {
			log.Fatalln(err)
		}
	}
}

func processFile(file string, override bool) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	im, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		return err
	}

	size := im.Bounds()
	length, offsetX, offsetY := 0, 0, 0
	if size.Max.X == size.Max.Y {
		fmt.Fprintf(os.Stderr, "Already square: %s\n", file)
		return nil
	}
	if size.Max.X > size.Max.Y {
		length = size.Max.X
		offsetY = (size.Max.X - size.Max.Y) / 2
	} else {
		length = size.Max.Y
		offsetX = (size.Max.Y - size.Max.X) / 2
	}

	// write rect to new image
	dest := image.NewRGBA(image.Rect(0, 0, length, length))
	destBounds := dest.Bounds()
	destBounds.Min.X += offsetX
	destBounds.Min.Y += offsetY
	draw.Draw(dest, destBounds, im, image.Pt(0, 0), draw.Src)

	// add ".cut" to extenstions
	if !override {
		ext := filepath.Ext(file)
		file = file[:len(file)-len(ext)] + ".square" + ext
	}

	f, err = os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, dest)
}
