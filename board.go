package main

import (
	"fmt"
	"strings"
)

// Create a width x height board
func createBoard(width int, height int) *Board {
	b := new(Board)
	b.Width = width
	b.Height = height
	grid := make([][]entity, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]entity, width)
		for j := 0; j < width; j++ {
			grid[i][j] = empty()
		}
	}

	b.Grid = grid
	return b
}

func (b Board) copy() *Board {
	copy := createBoard(b.Width, b.Height)
	for i := 0; i < copy.Height; i++ {
		for j := 0; j < copy.Width; j++ {
			copy.Grid[i][j] = getEntity(b.Grid[i][j].EntityType)
		}
	}

	return copy
}

// Return if an x coordinate is in bounds
func (b Board) xInBounds(xpos int) bool {
	return (0 <= xpos && xpos < b.Width)
}

// Return if an x coordinate is in bounds
func (b Board) yInBounds(ypos int) bool {
	return (0 <= ypos && ypos < b.Height)
}

// Return what is in the grid at point p
func (b Board) getTile(p Coord) entity {
	if b.tileInBounds(p) {
		return b.Grid[p.Y][p.X]
	}
	return invalid()
}

// Return if a point p is in bounds
func (b Board) tileInBounds(p Coord) bool {
	return (b.xInBounds(p.X) && b.yInBounds(p.Y))
}

func (b Board) insert(p Coord, e entity) {
	if b.xInBounds(p.X) && b.yInBounds(p.Y) {
		b.Grid[p.Y][p.X] = e
	}
}

func (b Board) getValidTiles(p Coord) []Coord {
	validTiles := make([]Coord, 0)
	potential := []Coord{
		Coord{p.X - 1, p.Y},
		Coord{p.X + 1, p.Y},
		Coord{p.X, p.Y - 1},
		Coord{p.X, p.Y + 1},
	}
	for i, p := range potential {
		if b.tileInBounds(p) && b.getTile(p).EntityType != OBSTACLE {
			validTiles = append(validTiles, potential[i])
		}
	}
	return validTiles
}

func (b Board) show() {
	rowDivider := strings.Repeat(" ---", b.Width)
	println(rowDivider)
	for i := 0; i < b.Height; i++ {
		print("| ")
		for j := 0; j < b.Width; j++ {
			fmt.Printf("%s |", b.Grid[i][j].Display)
			print(" ")
		}
		println("\n" + rowDivider)
	}
}
