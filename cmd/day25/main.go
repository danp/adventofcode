package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/danp/adventofcode2019/intcode"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	program, err := intcode.Parse(string(b))
	if err != nil {
		panic(err)
	}

	mem := make([]int, 8192)

	if len(os.Args) > 1 && os.Args[1] == "auto" {
		control := newAutoControl()
		if err := intcode.Run(program, mem, control.input, control.output); err != nil {
			panic(err)
		}

		fmt.Println(control.outb.String())
	} else {
		var inb bytes.Buffer
		sc := bufio.NewScanner(os.Stdin)

		input := func() (int, error) {
			if inb.Len() == 0 {
				if sc.Scan() {
					inb.Write(sc.Bytes())
					inb.WriteString("\n")
				}

				if sc.Err() != nil {
					return 0, sc.Err()
				}
			}

			c, err := inb.ReadByte()
			return int(c), err
		}

		output := func(x int) error {
			if x < 0 || x > 255 {
				return fmt.Errorf("non-ascii output %d", x)
			}
			fmt.Print(string(x))
			return nil
		}

		if err := intcode.Run(program, mem, input, output); err != nil {
			panic(err)
		}
	}
}

type autoControl struct {
	inb  bytes.Buffer
	outb bytes.Buffer

	toTake []string

	rroot *room
	cr    *room
}

type room struct {
	name string
	dirs map[string]*room
}

func newRoom() *room {
	return &room{name: "unknown", dirs: make(map[string]*room)}
}

func newAutoControl() *autoControl {
	rr := newRoom()
	ac := &autoControl{
		rroot: rr,
		cr:    rr,
	}
	return ac
}

func (c *autoControl) input() (int, error) {
	if c.inb.Len() == 0 {
		// there must be output waiting for us to check
		outs := strings.TrimSpace(c.outb.String())
		c.outb.Reset()

		for _, section := range strings.Split(outs, "\n\n") {
			slines := strings.Split(section, "\n")
			if strings.HasPrefix(slines[0], "== ") {
				c.cr.name = slines[0][3 : len(slines[0])-3]
			} else if strings.HasPrefix(slines[0], "Doors here lead:") {
				for _, sl := range slines[1:] {
					dir := sl[2:]
					if _, ok := c.cr.dirs[dir]; !ok {
						nr := newRoom()
						nr.dirs[invdir(dir)] = c.cr
						c.cr.dirs[dir] = nr
					}
				}
			} else if strings.HasPrefix(slines[0], "Items here:") {
				for _, sl := range slines[1:] {
					item := sl[2:]
					if item == "infinite loop" || item == "escape pod" || item == "molten lava" || item == "giant electromagnet" || item == "photons" {
						continue
					}
					c.toTake = append(c.toTake, item)
				}
			} else if strings.HasPrefix(slines[0], "Command?") {
			} else if strings.HasPrefix(slines[0], "You take the") {
			} else {
				fmt.Println("unknown section", section)
			}
		}

		if len(c.toTake) > 0 {
			fmt.Println("taking", c.toTake[0])
			fmt.Fprintln(&c.inb, "take", c.toTake[0])
			c.toTake = c.toTake[1:]
		} else {
			nr, nd := c.nextUnknownRoomMove()
			if nr == nil {
				// we've discovered everything, move to the checkpoint
				nr, nd = c.nextNamedRoomMove("Security Checkpoint")
			}
			if nr == nil {
				// we have discovered everything and we are at the checkpoint
				for _, i := range []string{"mouse", "monolith", "manifold", "space law space brochure"} {
					fmt.Fprintln(&c.inb, "drop", i)
				}
				fmt.Fprintln(&c.inb, "west")
			} else {
				fmt.Printf("moving %q from %q to %q\n", nd, c.cr.name, nr.name)
				fmt.Fprintln(&c.inb, nd)
				c.cr = nr
			}
		}
	}

	b, err := c.inb.ReadByte()
	return int(b), err
}

func (c *autoControl) output(x int) error {
	return c.outb.WriteByte(byte(x))
}

func (c *autoControl) nextUnknownRoomMove() (*room, string) {
	rp := c.roomPath(func(r *room) bool {
		return r.name == "unknown" && (r.dirs["east"] == nil || r.dirs["east"].name != "Security Checkpoint")
	})

	if len(rp) < 2 {
		return nil, ""
	}

	for dir, r := range c.cr.dirs {
		if r == rp[1] {
			return r, dir
		}
	}

	panic("got here")
}

func (c *autoControl) nextNamedRoomMove(name string) (*room, string) {
	rp := c.roomPath(func(r *room) bool {
		return r.name == name
	})

	if len(rp) < 2 {
		return nil, ""
	}

	for dir, r := range c.cr.dirs {
		if r == rp[1] {
			return r, dir
		}
	}

	panic("got here")
}

func (c *autoControl) roomPath(check func(r *room) bool) []*room {
	sr := make(map[*room]bool)
	q := [][]*room{{c.cr}}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		pr := p[len(p)-1]
		if check(pr) {
			return p
		}

		for _, r := range pr.dirs {
			if sr[r] {
				continue
			}
			sr[r] = true
			newp := make([]*room, len(p)+1)
			copy(newp, p)
			newp[len(newp)-1] = r
			q = append(q, newp)
		}
	}

	return nil
}

func invdir(dir string) string {
	return map[string]string{
		"north": "south",
		"south": "north",
		"east":  "west",
		"west":  "east",
	}[dir]
}
