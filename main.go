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

type NamedImage struct {
	img  image.Image
	name string
}

const dimension = 200

func main() {
	dir, out := parseArguments()

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	imgs := readImages(files, dir)

	width, length, remainder := calculateDimensions(len(imgs))

	var padder image.Image
	if remainder > 0 {
		padder, err = readPath("logo.png")
		if err != nil {
			log.Fatal(err)
		}
	}

	outImg := image.NewRGBA(image.Rect(0, 0, width*dimension, length*dimension))
	for index, img := range imgs {
		x := (index % width)
		y := (int(math.Floor(float64(index) / float64(width))))

		if y == (length-1) && remainder == 2 {
			if x == 0 {
				drawImage(padder, outImg, x, y)
			}
			x++
		}
		drawImage(img.img, outImg, x, y)
	}

	if remainder > 0 {
		drawImage(padder, outImg, width-1, length-1)
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

func readImages(files []fs.DirEntry, dir string) []NamedImage {
	var imgs []NamedImage
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	for _, file := range files {
		img, err := readFile(dir, file)
		if err != nil {
			log.Println(file.Name(), err)
		}
		imgs = append(imgs, img)
	}
	return imgs
}

func readFile(dir string, file fs.DirEntry) (NamedImage, error) {
	path := filepath.Join(dir, file.Name())
	img, err := readPath(path)
	return NamedImage{img, file.Name()}, err
}

func readPath(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
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
