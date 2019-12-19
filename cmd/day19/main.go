package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"

	"github.com/danp/adventofcode2019/intcode"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	program, err := intcode.Parse(string(b))
	if err != nil {
		panic(err)
	}

	// part 1 missing, lost to hacking on part 2

	bc := bc{p: program, cache: make(map[image.Point]bool)}

	var (
		pos   = image.Pt(0, 100)
		bminX int
	)

	for {
		in := bc.inBeam(pos)
		if in && bminX == 0 {
			bminX = pos.X
		}
		if !in && bminX > 0 {
			trc := pos.Sub(image.Pt(1, 0))
			tlc := trc.Sub(image.Pt(99, 0))
			blc := tlc.Add(image.Pt(0, 99))
			brc := blc.Add(image.Pt(99, 0))

			// we know top right corner is in the beam, check the rest
			if bc.inBeam(tlc) && bc.inBeam(blc) && bc.inBeam(brc) {
				fmt.Println(pos, tlc, trc, blc, brc, tlc.X*10000+tlc.Y)
				display(&bc, image.Pt(bminX-20, pos.Y-10), 120, trc, tlc, blc, brc)
				break
			}

			pos.Y++
			pos.X = bminX - 5
			bminX = 0
			continue
		}

		pos = pos.Add(image.Pt(1, 0))
	}

	fmt.Println(bc.hits)
}

type bc struct {
	p     []int
	cache map[image.Point]bool
	hits  int
}

func (b *bc) inBeam(pt image.Point) bool {
	if c, ok := b.cache[pt]; ok {
		b.hits++
		return c
	}

	var ist int
	input := func() (int, error) {
		if ist == 0 {
			ist++
			return pt.X, nil
		}
		ist = 0
		return pt.Y, nil
	}

	var in bool
	output := func(x int) error {
		in = x == 1
		return nil
	}

	mem := make([]int, 4096)
	if err := intcode.Run(b.p, mem, input, output); err != nil {
		panic(err)
	}

	b.cache[pt] = in
	return in
}

func display(bc *bc, start image.Point, rows int, hl ...image.Point) {
	for y := start.Y; y <= start.Y+rows; y++ {
		out := fmt.Sprintf("%-8d", y)

		for x := start.X; x < start.X+400; x++ {
			pt := image.Pt(x, y)

			int := containtsPt(pt, hl)
			if bc.inBeam(pt) {
				if int {
					out += "H"
				} else {
					out += "#"
				}
			} else {
				if int {
					out += "x"
				} else {
					out += "."
				}
			}
		}

		fmt.Println(out)
	}
}

func containtsPt(pt image.Point, pts []image.Point) bool {
	for _, cpt := range pts {
		if cpt.Eq(pt) {
			return true
		}
	}
	return false
}
