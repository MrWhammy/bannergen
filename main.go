package main

import (
	"embed"
	"flag"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"log"
	"math"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"golang.org/x/image/draw"
)

type logo string

//go:embed logo.png sponsors/*
var logos embed.FS

func (l *logo) readImage() image.Image {
	f, err := logos.Open(string(*l))
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	dir := "sponsors"

	files, err := logos.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	logos := readLogos(files, dir)

	width, length, remainder := calculateDimensions(len(logos))

	padder := logo("logo.png")

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

	writeImage("out.png", outImg)
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

func readLogos(files []fs.DirEntry, dir string) []logo {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	imgs := make([]logo, 0, len(files))
	for _, file := range files {
		img := readFile(dir, file)
		imgs = append(imgs, img)
	}
	return imgs
}

func readFile(dir string, file fs.DirEntry) logo {
	path := filepath.Join(dir, file.Name())
	return logo(path)
}
