package main

import (
	"errors"
	"fmt"
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

	mzl, err := findFewestZeroes(layers)
	if err != nil {
		panic(err)
	}

	fmt.Println(mzl)

	up := countUniquePixels(layers[mzl])
	fmt.Println(up)

	fmt.Println(up[1] * up[2])

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
		for _, row := range l {
			for _, x := range row {
				if x == 0 {
					zcs[i]++
				}
			}
		}
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
