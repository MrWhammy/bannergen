package main

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"
)

type Logo struct {
	path string
}

func (l *Logo) readImage() image.Image {
	f, err := os.Open(l.path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

const dimension = 200

func main() {
	dir, out := parseArguments()

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	logos := readLogos(files, dir)

	width, length, remainder := calculateDimensions(len(logos))

	padder := Logo{"logo.png"}

	outImg := image.NewRGBA(image.Rect(0, 0, width*dimension, length*dimension))
	for index, logo := range logos {
		x := (index % width)
		y := (int(math.Floor(float64(index) / float64(width))))

		bottomRow := y == (length - 1)
		if bottomRow && remainder == 2 { // bottom row requires left padding
			if x == 0 {
				drawImage(padder.readImage(), outImg, x, y)
			}
			x++
		}
		drawImage(logo.readImage(), outImg, x, y)
	}

	if remainder > 0 { // padding required on the right
		drawImage(padder.readImage(), outImg, width-1, length-1)
	}

	writeImage(out, outImg)
}

func drawImage(img image.Image, outImg draw.Image, x int, y int) {
	draw.BiLinear.Scale(outImg, image.Rect(x*dimension, y*dimension, (x+1)*dimension, (y+1)*dimension), img, img.Bounds(), draw.Src, nil)
}

var widths = []int{
	3, 3, 3, 6, 6, 6, 7, 4, 5, 5,
	6, 6, 7, 7, 5, 6, 6, 6, 7, 7,
	7, 6, 6, 6, 5, 7, 7, 7, 6, 6,
}

func calculateDimensions(imgNumber int) (int, int, int) {
	width := widths[imgNumber-1]
	length := (imgNumber / width)
	if imgNumber%width > 0 {
		length++
	}
	remainder := (width * length) - imgNumber
	return width, length, remainder
}

func writeImage(out string, outImg image.Image) {
	myfile, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer myfile.Close()
	png.Encode(myfile, outImg)
}

func readLogos(files []fs.DirEntry, dir string) []Logo {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	imgs := make([]Logo, 0, len(files))
	for _, file := range files {
		img := readFile(dir, file)
		imgs = append(imgs, img)
	}
	return imgs
}

func readFile(dir string, file fs.DirEntry) Logo {
	path := filepath.Join(dir, file.Name())
	return Logo{path}
}

func parseArguments() (string, string) {
	if len(os.Args) > 2 {
		return os.Args[1], os.Args[2]
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		sponsorDir := filepath.Join(dir, "sponsors")
		if len(os.Args) > 1 {
			return sponsorDir, os.Args[1]
		} else {
			return sponsorDir, "out.png"
		}
	}
}
