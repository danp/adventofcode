package main

import (
	"container/heap"
	"fmt"
	"image"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	indexes := make(map[string]int)
	grid := make(map[image.Point]string)
	for y, l := range lines {
		for x, c := range l {
			pt := image.Pt(x, y)
			switch c {
			case ' ', '#':
			case '.':
				grid[pt] = ""
			default:
				grid[pt] = string(c) + fmt.Sprint(indexes[string(c)]+1)
				indexes[string(c)]++
			}
		}
	}

	en := dijkstra(grid)

	fmt.Printf("en: %v\n", en)
}

type gridstate struct {
	a1, a2 image.Point
	b1, b2 image.Point
	c1, c2 image.Point
	d1, d2 image.Point
}

var rooms = []image.Rectangle{
	image.Rect(3, 2, 3, 3),
	image.Rect(5, 2, 5, 3),
	image.Rect(7, 2, 7, 3),
	image.Rect(9, 2, 9, 3),
}

// 3/5/7/9,2-3
func (g gridstate) isDone() bool {
	return image.Rect(g.a1.X, g.a2.X, g.a1.Y, g.a2.Y).Eq(rooms[0]) &&
		image.Rect(g.b1.X, g.b2.X, g.b1.Y, g.b2.Y).Eq(rooms[1]) &&
		image.Rect(g.c1.X, g.c2.X, g.c1.Y, g.c2.Y).Eq(rooms[2]) &&
		image.Rect(g.d1.X, g.d2.X, g.d1.Y, g.d2.Y).Eq(rooms[3])
}

func gs(grid map[image.Point]string) gridstate {
	pts := make(map[string]image.Point)
	for pt, v := range grid {
		if v == "" {
			continue
		}
		pts[v] = pt
	}
	return gridstate{
		a1: pts["A1"],
		a2: pts["A2"],
		b1: pts["B1"],
		b2: pts["B2"],
		c1: pts["C1"],
		c2: pts["C2"],
		d1: pts["D1"],
		d2: pts["D2"],
	}
}

func dijkstra(grid map[image.Point]string) int {
	start := gs(grid)
	dists := map[gridstate]int{start: 0}
	prev := make(map[gridstate]gridstate)
	var pq PriorityQueue
	item := &Item{value: start, priority: 0, index: 0}
	pq = append(pq, item)
	heap.Init(&pq)
	visited := map[gridstate]struct{}{
		start: {},
	}

	var endstate gridstate
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		state := item.value

		if state.isDone() {
			endstate = state
			break
		}

		for _, ns := range next(visited, state) {
			alt := dists[state] + ns.cost
			if nd, ok := dists[ns.state]; !ok || alt < nd {
				dists[ns.state] = alt
				pq.upsert(ns.state, alt)
				prev[ns.state] = state
			}
		}
	}

	u := endstate
	var cost int
	for u != start {
		cost += dists[u]
		u = prev[u]
	}
	return cost
}

type nextstate struct {
	cost  int
	state gridstate
}

func next(visited map[gridstate]struct{}, cur gridstate) []nextstate {
	// can't stop in a space outside a room
	// can't move from hallway into a room unless:
	//   that room is their dest AND
	//   room contains no amphipods which do not also have that room as their own destination
	// Once an amphipod stops moving in the hallway, it will stay in that spot until it can move into a room
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

type Item struct {
	value    gridstate // The value of the item; arbitrary.
	priority int       // The priority of the item in the queue.
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

func (pq *PriorityQueue) upsert(st gridstate, priority int) {
	for _, i := range *pq {
		if i.value == st {
			pq.update(i, priority)
			return
		}
	}
	pq.push(&Item{value: st, priority: priority})
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) push(item *Item) {
	heap.Push(pq, item)
}
