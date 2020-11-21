package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	d, err := parse(string(b))
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "orbits":
		fmt.Println(orbitCount(d))
	case "transfers":
		fmt.Println(transferCount(d, "YOU", "SAN"))
	}
}

type node struct {
	name     string
	parent   *node
	children map[string]*node
}

func parse(input string) (map[string]*node, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	data := make(map[string]*node)

	for _, l := range lines {
		parts := strings.Split(l, ")")
		k := parts[0]
		v := parts[1]

		if data[k] == nil {
			data[k] = &node{name: k}
		}
		if data[k].children == nil {
			data[k].children = make(map[string]*node)
		}
		if data[v] == nil {
			data[v] = &node{name: v}
		}
		if data[v].parent == nil {
			data[v].parent = data[k]
		}
		children := data[k].children
		children[v] = data[v]
		data[k].children = children
	}

	return data, nil
}

func orbitCount(d map[string]*node) int {
	var c int
	seen := make(map[string]bool)
	q := [][]string{{"COM"}}

	for len(q) > 0 {
		p := q[0]
		k := p[len(p)-1]
		q = q[1:]

		c += len(p) - 1

		for j, _ := range d[k].children {
			if seen[j] {
				continue
			}
			seen[j] = true

			newp := make([]string, len(p))
			copy(newp, p)
			newp = append(newp, j)
			q = append(q, newp)
		}
	}

	return c
}

func transferCount(d map[string]*node, source, target string) int {
	var c int
	seen := make(map[string]bool)
	q := [][]string{{source}}

	for len(q) > 0 {
		p := q[0]
		k := p[len(p)-1]
		q = q[1:]

		if d[k].parent != nil && d[k].parent.name == target {
			panic("need to figure this out")
		}

		enqueue := func(k string) {
			newp := make([]string, len(p))
			copy(newp, p)
			newp = append(newp, k)
			q = append(q, newp)
		}

		for j, _ := range d[k].children {
			if j == target {
				fmt.Println(p)
				return len(p) - 2
			}

			if seen[j] {
				continue
			}
			seen[j] = true

			enqueue(j)
		}

		if d[k].parent != nil && !seen[d[k].parent.name] {
			enqueue(d[k].parent.name)
		}
	}

	return c
}
