package main

import (
	"fmt"
	"math"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	for _, l := range lines {
		bits := hexToBits(l)
		br := &bitReader{bits}
		top := parse(br)
		fmt.Printf("%v version sum: %v\n", l, versionSum(top))
	}
}

func versionSum(p packet) int {
	n := p.version
	if subs, ok := p.value.([]packet); ok {
		for _, sp := range subs {
			n += versionSum(sp)
		}
	}
	return n
}

type bitReader struct {
	bits []int
}

func (b *bitReader) read(n int) []int {
	if len(b.bits) < n {
		panic(fmt.Sprintf("don't have %v bits, have %v", n, len(b.bits)))
	}
	p := b.bits[:n]
	b.bits = b.bits[n:]
	return p
}

type packet struct {
	version int
	typ     int
	value   interface{}
}

func parse(br *bitReader) packet {
	v := bitsToInt(br.read(3))
	typ := bitsToInt(br.read(3))

	p := packet{version: v, typ: typ}

	switch typ {
	case 4:
		var nbits []int
		for {
			chunk := br.read(5)
			nbits = append(nbits, chunk[1:]...)
			if chunk[0] == 0 {
				break
			}
		}
		p.value = bitsToInt(nbits)
	default:
		var subs []packet
		lt := br.read(1)
		switch lt[0] {
		case 0:
			sl := bitsToInt(br.read(15))
			sb := br.read(sl)
			sr := &bitReader{bits: sb}
			for len(sr.bits) > 0 {
				p := parse(sr)
				subs = append(subs, p)
			}
		case 1:
			sp := bitsToInt(br.read(11))
			for i := 0; i < sp; i++ {
				subs = append(subs, parse(br))
			}
		}
		p.value = subs
	}

	return p
}

func bitsToInt(bits []int) int {
	var n int
	for i := 0; i < len(bits); i++ {
		b := bits[len(bits)-1-i]
		if b == 1 {
			n += int(math.Pow(2, float64(i)))
		}
	}
	return n
}

func hexToBits(s string) []int {
	var out []int
	for _, c := range s {
		ci := int(c)
		var n int
		if ci >= '0' && ci <= '9' {
			n = ci - '0'
		} else if ci >= 'A' && ci <= 'F' {
			n = 10 + ci - 'A'
		}
		for _, c := range fmt.Sprintf("%04b", n) {
			out = append(out, int(c)-'0')
		}
	}
	return out
}
