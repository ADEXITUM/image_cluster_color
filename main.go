package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"sort"

	"golang.org/x/image/draw"
)

type Pixel [3]float64

func distance(a, b Pixel) float64 {
	var d float64
	for i := 0; i < 3; i++ {
		d += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(d)
}

func averagePixels(pixels []Pixel) Pixel {
	var avg Pixel
	for _, p := range pixels {
		for i := 0; i < 3; i++ {
			avg[i] += p[i]
		}
	}
	n := float64(len(pixels))
	for i := 0; i < 3; i++ {
		avg[i] /= n
	}
	return avg
}

func kmeans(pixels []Pixel, k int, iterations int) []Pixel {
	centroids := make([]Pixel, k)
	copy(centroids, pixels[:k])

	for it := 0; it < iterations; it++ {
		clusters := make([][]Pixel, k)
		for _, p := range pixels {
			minDist := math.MaxFloat64
			minIdx := 0
			for i, c := range centroids {
				d := distance(p, c)
				if d < minDist {
					minDist = d
					minIdx = i
				}
			}
			clusters[minIdx] = append(clusters[minIdx], p)
		}

		for i := 0; i < k; i++ {
			if len(clusters[i]) > 0 {
				centroids[i] = averagePixels(clusters[i])
			}
		}
	}
	return centroids
}

func printColorBlock(c Pixel) {
	r, g, b := int(c[0]), int(c[1]), int(c[2])
	fmt.Printf("\x1b[48;2;%d;%d;%dm  \x1b[0m #%02x%02x%02x\n", r, g, b, r, g, b)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image_path>")
		return
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	pixels := make([]Pixel, 0, 100*100)
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			r, g, b, _ := dst.At(x, y).RGBA()
			p := Pixel{float64(r >> 8), float64(g >> 8), float64(b >> 8)}
			pixels = append(pixels, p)
		}
	}

	centroids := kmeans(pixels, 5, 10)

	sort.Slice(centroids, func(i, j int) bool {
		sumI := centroids[i][0] + centroids[i][1] + centroids[i][2]
		sumJ := centroids[j][0] + centroids[j][1] + centroids[j][2]
		return sumI < sumJ
	})

	for i, c := range centroids {
		fmt.Printf("Color %d: ", i+1)
		printColorBlock(c)
	}
}
