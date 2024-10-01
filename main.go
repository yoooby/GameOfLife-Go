package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1080
	screenHeight = 1080
)

type World struct {
	area   [][]bool
	width  int
	height int
}

var dirs = [8][2]int8{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
	{1, 1},
	{-1, -1},
	{1, -1},
	{-1, 1},
}

func getNeighbors(state [][]bool, y, x int) int {
	count := 0
	for _, dir := range dirs {
		c := x + int(dir[0])
		r := y + int(dir[1])

		if r < 0 || r >= len(state) {
			continue
		}
		if c < 0 || c >= len(state[0]) {
			continue
		}

		if state[r][c] {
			count++
		}
	}
	return count
}

func (g *Game) play() {
	out := make([][]bool, len(g.state.area))
	for i := range out {
		out[i] = make([]bool, len(g.state.area[i]))
	}

	/*
			A live cell dies if it has fewer than two live neighbors.
		    A live cell with two or three live neighbors lives on to the next generation.
		    A live cell with more than three live neighbors dies.
		    A dead cell will be brought back to live if it has exactly three live neighbors.

			< 2 He9 mcha
			 == 3 live because 2 or 3 stays so no need
			 3 < Heee9 mcha no matter what

	*/

	for row := range g.state.area {
		for col := range row {
			count := getNeighbors(g.state.area, row, col)

			if count < 2 {
				out[row][col] = false
			}
			if count > 3 {
				out[row][col] = false
			}
			if count == 3 {
				out[row][col] = true
			}
		}
	}
	g.state.area = out

}

type Game struct {
	state  *World
	pixels []byte
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.state.area[y][y] =  !g.state.area[y][x]
	}
	g.play()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Zeb a zebi!")

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight

}

func main() {
	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
