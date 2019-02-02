package main

import (
	"math/rand"
	"time"
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type MSGrid [][]int

func NewMSGrid(w, h int) MSGrid {
	grid := make([][]int, h)
	for i := range grid {
		grid[i] = make([]int, w)
	}
	return grid
}

func (g MSGrid) Populate(numberOfMines int) {
	h := len(g)
	w := len(g[0])
	for i := 0; i < numberOfMines; i++ {
		x := rand.Int() % w
		y := rand.Int() % h
		if g[y][x] == 9 {
			i -= 1
			continue
		}
		g[y][x] = 9
	}
}

func (g MSGrid) updateMineCount() {
	// tricky looping bullshit time:
	// for each grid point
	for y := range g {
		for x := range g[y] {
			g.updateCell(x, y)
		}
	}
}

func (g MSGrid) updateCell(x, y int) {
	// if it's not a mine
	if g[y][x] != 9 {
		mineCount := 0
		// read each spot around it
		for i := y - 1; i <= y + 1; i++ {
			for j := x - 1; j <= x + 1; j++ {
				// if it's out of bounds, discard it
				if i < 0 || j < 0 || i >= len(g) || j >= len(g[y]) {
					continue
				}
				// if that spot is a mine, up the counter of nearby mines
				if g[i][j] == 9 {
					mineCount += 1
				}
			}
		}
		// set cell to the count of nearby mines
		g[y][x] = mineCount
	}
}
