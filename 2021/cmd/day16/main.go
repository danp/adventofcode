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
		fmt.Printf("%v version sum: %v result: %v\n", l, versionSum(top), eval(top))
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

func eval(p packet) int {
	if p.typ == 4 {
		return p.value.(int)
	}

	subs := p.value.([]packet)

	switch p.typ {
	case 0:
		var n int
		for _, sp := range subs {
			n += eval(sp)
		}
		return n
	case 1:
		n := -1
		for _, sp := range subs {
			v := eval(sp)
			if n == -1 {
				n = v
			} else {
				n *= v
			}
		}
		return n
	case 2:
		n := -1
		for _, sp := range subs {
			m := eval(sp)
			if n == -1 || m < n {
				n = m
			}
		}
		return n
	case 3:
		n := -1
		for _, sp := range subs {
			m := eval(sp)
			if n == -1 || m > n {
				n = m
			}
		}
		return n
	case 5:
		var n int
		a, b := eval(subs[0]), eval(subs[1])
		if a > b {
			n = 1
		}
		return n
	case 6:
		a, b := eval(subs[0]), eval(subs[1])
		var n int
		if a < b {
			n = 1
		}
		return n
	case 7:
		a, b := eval(subs[0]), eval(subs[1])
		var n int
		if a == b {
			n = 1
		}
		return n
	}
	panic("unknown")
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
