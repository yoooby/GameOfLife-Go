package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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

	for row := range g.state.area {
		for col := range g.state.area[row] {
			count := getNeighbors(g.state.area, row, col)

			if g.state.area[row][col] {
				out[row][col] = count == 2 || count == 3
			} else {
				out[row][col] = count == 3
			}
		}
	}
	g.state.area = out
}

type Game struct {
	state   *World
	pixels  []byte
	running bool
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cellX := x / (screenWidth / g.state.width)
		cellY := y / (screenHeight / g.state.height)
		if cellX >= 0 && cellX < g.state.width && cellY >= 0 && cellY < g.state.height {
			g.state.area[cellY][cellX] = !g.state.area[cellY][cellX]
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.running = !g.running
	}

	if g.running {
		g.play()
	}
	return nil
}

func (g *Game) Layout(outHeight, outWidth int) (int, int) {
	return screenHeight, screenWidth
}
func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < g.state.height; y++ {
		for x := 0; x < g.state.width; x++ {
			if g.state.area[y][x] {
				for dy := 0; dy < screenHeight/g.state.height; dy++ {
					for dx := 0; dx < screenWidth/g.state.width; dx++ {
						screen.Set(x*(screenWidth/g.state.width)+dx, y*(screenHeight/g.state.height)+dy, color.White)
					}
				}
			}
		}
	}
}

func main() {
	world := &World{area: make([][]bool, screenHeight), width: screenWidth / 10, height: screenHeight / 10}
	for i := range world.area {
		world.area[i] = make([]bool, world.width)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(&Game{state: world}); err != nil {
		log.Fatal(err)
	}
}
