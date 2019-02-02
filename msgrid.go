package main

import (
	"math/rand"
	"time"
)

var rng *rand.Rand

func init() {
	// initialize RNG so our boards aren't always the same on each run
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// MSGrid represents a Minesweeper grid.
// This simple grid of integers is largely exactly what you would see in an actual Minesweeper board.
// Cells set to 9 signify mines.
type MSGrid [][]int

// NewMSGrid instantiates a new MSGrid object.
func NewMSGrid(w, h int) MSGrid {
	grid := make([][]int, h)
	for i := range grid {
		grid[i] = make([]int, w)
	}
	return grid
}

// Populate populates an MSGrid with mines.
func (g MSGrid) Populate(numberOfMines int) {
	h := len(g)
	w := len(g[0])
	for i := 0; i < numberOfMines; i++ {
		// get a random point and set it to 9
		x := rand.Int() % w
		y := rand.Int() % h
		// skip if the point is already 9
		if g[y][x] == 9 {
			i -= 1
			continue
		}
		g[y][x] = 9
	}
}

// updateMineCount is a function that writes the mine proximity numbers.
func (g MSGrid) updateMineCount() {
	// for each grid point, update the cell
	for y := range g {
		for x := range g[y] {
			g.updateCell(x, y)
		}
	}
}

// updateCell tries to update a single cell's number by looking at nearby cells.
func (g MSGrid) updateCell(x, y int) {
	// if it's not a mine
	if g[y][x] != 9 {
		// instantiate a mine counter
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
