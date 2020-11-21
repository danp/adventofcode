package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	layers, err := extractLayers(strings.TrimSpace(string(b)), 25, 6)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "checksum":
		mzl, err := findFewestZeroes(layers)
		if err != nil {
			panic(err)
		}

		fmt.Println(mzl)

		up := countUniquePixels(layers[mzl])
		fmt.Println(up)

		fmt.Println(up[1] * up[2])
	case "decode-ascii", "decode-png":
		layer, err := decodeLayers(layers)
		if err != nil {
			panic(err)
		}

		switch os.Args[1] {
		case "decode-ascii":
			for _, r := range layer {
				for _, p := range r {
					if p == 1 {
						fmt.Print("*")
					} else {
						fmt.Print(" ")
					}
				}
				fmt.Println()
			}
		case "decode-png":
			img := layerImage(layer)
			if err := png.Encode(os.Stdout, img); err != nil {
				panic(err)
			}
		}
	}

}

func extractLayers(data string, width, height int) ([][][]int, error) {
	layerSize := width * height

	if len(data)%layerSize != 0 {
		return nil, errors.New("data has incorrect length for width and height")
	}

	layerCount := len(data) / layerSize
	out := make([][][]int, layerCount)

	for i := 0; i < layerCount; i++ {
		out[i] = make([][]int, height)
		for j := 0; j < height; j++ {
			out[i][j] = make([]int, width)
			for k := 0; k < width; k++ {
				pos := (i * width * height) + (j * width) + k
				p, err := strconv.Atoi(string(data[pos]))
				if err != nil {
					return nil, err
				}
				out[i][j][k] = p
			}
		}
	}

	return out, nil
}

func findFewestZeroes(layers [][][]int) (int, error) {
	zcs := make(map[int]int)

	for i, l := range layers {
		up := countUniquePixels(l)
		zcs[i] = up[0]
	}

	var fewestLayer, fewestZeroes int = -1, 0
	for l, c := range zcs {
		if fewestLayer == -1 {
			fewestLayer = l
			fewestZeroes = c
			continue
		}
		if c < fewestZeroes {
			fewestLayer = l
			fewestZeroes = c
		}
	}

	return fewestLayer, nil
}

func countUniquePixels(layer [][]int) map[int]int {
	up := make(map[int]int)

	for _, row := range layer {
		for _, x := range row {
			up[x]++
		}
	}

	return up
}

func decodeLayers(layers [][][]int) ([][]int, error) {
	if len(layers) == 0 {
		return nil, nil
	}

	if len(layers[0]) == 0 {
		return nil, nil
	}

	height := len(layers[0])
	width := len(layers[0][0])

	out := make([][]int, height)

	for i := 0; i < height; i++ {
		out[i] = make([]int, width)
		for j := 0; j < width; j++ {
			for k := 0; k < len(layers); k++ {
				p := layers[k][i][j]
				if p < 2 {
					out[i][j] = p
					break
				}
			}
		}
	}

	return out, nil
}

type layerImage [][]int

func (l layerImage) ColorModel() color.Model {
	return color.GrayModel
}

func (l layerImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(l[0]), len(l))
}

func (l layerImage) At(x, y int) color.Color {
	return color.Gray{Y: uint8(l[y][x]) * 255}
}
