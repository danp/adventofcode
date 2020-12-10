package main

import (
	"fmt"
	"sort"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	// initial connection to outlet
	adapters := append([]int{0}, scaffold.Ints(scaffold.Lines())...)
	sort.Ints(adapters)
	// final connection to device
	adapters = append(adapters, adapters[len(adapters)-1]+3)

	diffs := make(map[int]int)
	for i, a := range adapters {
		if i == 0 {
			continue
		}
		diffs[a-adapters[i-1]]++
	}

	fmt.Println(diffs, diffs[1]*diffs[3])

	graph := make(map[int][]int)
	for i, a := range adapters {
		graph[a] = []int{}
		for j := i - 1; j >= 0 && a-adapters[j] <= 3; j-- {
			aj := adapters[j]
			graph[aj] = append(graph[aj], a)
		}
	}

	for _, a := range adapters {
		fmt.Println(a, graph[a])
	}

	fmt.Println()

	paths := 1
	i := 0
	for {
		a := adapters[i]
		edges := graph[a]
		if len(edges) == 0 {
			fmt.Println(a, "end", paths)
			break
		}

		conv := findConvergence(graph, a)
		fmt.Println(a, edges, "converges at", conv)
		paths *= countPaths(graph, a, conv)
		for j := i + 1; j < len(adapters); j++ {
			if adapters[j] == conv {
				i = j
				break
			}
		}
	}

	fmt.Println(paths)
}

// convergence is the minimum node where all edges of
// the given node meet
func findConvergence(graph map[int][]int, node int) int {
	edges := graph[node]

	if len(edges) == 1 {
		return edges[len(edges)-1]
	}

	var conv int
	for _, e := range edges {
		econv := findFirstOneEdgeNode(graph, e)
		if econv > conv {
			conv = econv
		}
	}
	return conv
}

func findFirstOneEdgeNode(graph map[int][]int, start int) int {
	q := []int{start}
	for len(q) > 0 {
		k := q[0]
		q = q[1:]

		fmt.Println(start, k, graph[k])

		if len(graph[k]) == 1 {
			return graph[k][0]
		}

		for _, o := range graph[k] {
			q = append(q, o)
		}
	}

	panic("got here")
}

func countPaths(graph map[int][]int, start, end int) int {
	pq := [][]int{{start}}
	var paths int
	for len(pq) > 0 {
		path := pq[0]
		k := path[len(path)-1]
		pq = pq[1:]

		if k == end {
			fmt.Println(path)
			paths++
			continue
		}

		for _, o := range graph[k] {
			newp := make([]int, len(path))
			copy(newp, path)
			newp = append(newp, o)
			pq = append(pq, newp)
		}
	}
	return paths
}

func contains(set []int, i int) bool {
	for _, j := range set {
		if j == i {
			return true
		}
	}
	return false
}
