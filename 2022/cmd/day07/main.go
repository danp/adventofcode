package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/danp/adventofcode/scaffold"
	"golang.org/x/exp/maps"
)

func main() {
	lines := scaffold.Lines()

	type entry struct {
		dir  bool
		size int
	}

	fs := make(map[string]entry)

	pwd := "/"
	for _, l := range lines {
		f := strings.Fields(l)
		if strings.HasPrefix(l, "$") {
			switch f[1] {
			case "cd":
				pwd = filepath.Join(pwd, f[2])
				pwd = filepath.Clean(pwd)
				fs[pwd] = entry{true, 0}
			case "ls":
			}
			continue
		}
		if f[0] == "dir" {
			fs[filepath.Join(pwd, f[1])] = entry{true, 0}
		} else {
			fs[filepath.Join(pwd, f[1])] = entry{false, scaffold.Int(f[0])}
		}
	}

	walk := func(prefix string, fn func(p string, e entry)) {
		keys := maps.Keys(fs)
		sort.Strings(keys)
		for _, p := range keys {
			if !strings.HasPrefix(p, prefix) {
				continue
			}
			fn(p, fs[p])
		}
	}

	sums := make(map[string]int)
	walk("", func(p string, e entry) {
		if e.dir {
			return
		}
		parts := strings.Split(p, "/")
		for i := 1; i < len(parts); i++ {
			part := strings.Join(parts[:i], "/")
			if i == 1 {
				part = "/"
			}
			sums[part] += e.size
		}
	})

	var tot int
	for p, s := range sums {
		if p == "" {
			continue
		}
		if s <= 100_000 {
			tot += s
		}
	}
	fmt.Printf("tot under 100k: %v\n", tot)

	const totalSpace = 70000000
	const needUnused = 30000000

	unused := totalSpace - sums["/"]
	need := needUnused - unused

	var minp string
	var mins int
	for p, s := range sums {
		if s < need {
			continue
		}
		if mins > 0 && s > mins {
			continue
		}
		mins = s
		minp = p
	}

	fmt.Printf("minp: %v mins: %v\n", minp, mins)
}
