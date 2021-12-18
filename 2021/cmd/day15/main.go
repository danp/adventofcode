package main

import (
	"container/heap"
	"fmt"
	"image"
	"math"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	graph := make(map[image.Point]int)
	for y, l := range lines {
		for x, c := range l {
			pt := image.Pt(x, y)
			graph[pt] = int(c) - '0'
		}
	}

	var risk int
	gr := rect(graph)
	for _, pt := range dijkstra(graph, image.Point{}, gr.Max) {
		risk += graph[pt]
	}
	fmt.Println("part 1 risk:", risk)

	graph2 := makeGraph2(graph, 5)
	gr = rect(graph2)
	risk = 0
	for _, pt := range dijkstra(graph2, image.Point{}, gr.Max) {
		risk += graph2[pt]
	}
	fmt.Println("part 2 risk:", risk)
}

func makeGraph2(graph map[image.Point]int, n int) map[image.Point]int {
	newGraph := make(map[image.Point]int)
	rect := rect(graph)
	for pt, r := range graph {
		for yn := 0; yn < n; yn++ {
			for xn := 0; xn < n; xn++ {
				npt := image.Pt(pt.X+(xn*(rect.Max.X+1)), pt.Y+(yn*(rect.Max.Y+1)))
				nr := r + yn + xn
				if nr > 9 {
					nr -= 9
				}
				newGraph[npt] = nr
			}
		}
	}
	return newGraph
}

func dijkstra(graph map[image.Point]int, start, end image.Point) []image.Point {
	dists := map[image.Point]int{start: 0}
	prev := make(map[image.Point]image.Point)
	q := make(map[image.Point]*Item)
	var pq PriorityQueue
	var i int
	for pt := range graph {
		priority := math.MaxInt
		if pt == start {
			priority = 0
		}
		item := &Item{value: pt, priority: priority, index: i}
		pq = append(pq, item)
		q[pt] = item
		i++
	}
	heap.Init(&pq)

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		pt := item.value
		delete(q, pt)

		if pt == end {
			break
		}

		for _, npt := range neighbs(q, pt) {
			alt := dists[pt] + graph[npt]
			if nd, ok := dists[npt]; !ok || alt < nd {
				dists[npt] = alt
				pq.update(q[npt], npt, alt)
				prev[npt] = pt
			}
		}
	}

	var path []image.Point
	u := end
	for u != start {
		path = append(path, u)
		u = prev[u]
	}
	return path
}

func min(q map[image.Point]struct{}, dists map[image.Point]int) image.Point {
	mind := math.MaxInt
	var minp image.Point
	for pt := range q {
		if d, ok := dists[pt]; ok && d < mind {
			mind = d
			minp = pt
		}
	}
	if mind == math.MaxInt {
		for pt := range q {
			return pt
		}
	}
	return minp
}

var dirs = []image.Point{
	image.Pt(0, -1),
	image.Pt(1, 0),
	image.Pt(0, 1),
	image.Pt(-1, 0),
}

func neighbs(graph map[image.Point]*Item, pt image.Point) []image.Point {
	var out []image.Point
	for _, d := range dirs {
		npt := pt.Add(d)
		if _, ok := graph[npt]; ok {
			out = append(out, npt)
		}
	}
	return out
}

func rect(grid map[image.Point]int) image.Rectangle {
	var r image.Rectangle
	for pt := range grid {
		r.Max.X = max(r.Max.X, pt.X)
		r.Max.Y = max(r.Max.Y, pt.Y)
	}
	return r
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func show(graph map[image.Point]int) {
	r := rect(graph)

	for y := 0; y <= r.Max.Y; y++ {
		for x := 0; x <= r.Max.X; x++ {
			fmt.Print(graph[image.Pt(x, y)])
			if x > 0 && x%10 == 0 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		if y > 0 && y%10 == 0 {
			fmt.Println()
		}
	}
}

type Item struct {
	value    image.Point // The value of the item; arbitrary.
	priority int         // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value image.Point, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
