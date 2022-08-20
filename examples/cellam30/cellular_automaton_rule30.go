package main

import "fmt"

const ON, OFF = '\u2b1b', '\u2b1c'
const CELL_COUNT = 51
const STEP_COUNT = CELL_COUNT

type Cells [CELL_COUNT]bool

func (c Cells) String() string {
	cellChars := make([]rune, 0, CELL_COUNT)
	for i := 0; i < CELL_COUNT; i++ {
		if c[i] {
			cellChars = append(cellChars, ON)
		} else {
			cellChars = append(cellChars, OFF)
		}
	}
	return string(cellChars)
}

func main() {
	var c Cells
	c[CELL_COUNT/2] = true
	for i := 0; i < STEP_COUNT; i++ {
		fmt.Println(c)
		c.Transform()
	}
}

func (c *Cells) Transform() {
	var next Cells
	for i := range c {
		next[i] = c.nextOf(i)
	}
	for i := range c {
		c[i] = next[i]
	}
}

func (c Cells) nextOf(i int) bool {
	self := i % CELL_COUNT
	left := (i - 1 + CELL_COUNT) % CELL_COUNT
	right := (i + 1) % CELL_COUNT
	return rule30(c[left], c[self], c[right])
}

func rule30(left, self, right bool) bool {
	if left {
		return !self && !right
	}
	return self || right
}
