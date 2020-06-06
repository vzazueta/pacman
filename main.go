package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 600.0
	screenHeight = 800.0
)

func textureFromBMP(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading %v: %v", filename, err))
	}
	defer img.Free()
	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture from %v: %v", filename, err))
	}

	return tex
}

func getVisualNodes(renderer *sdl.Renderer) [][]visualNode {
	cellSize := screenWidth / Dimension

	fmt.Println(cellSize)
	output := make([][]visualNode, Dimension)

	yPos := cellSize / 2.0

	for i := 0; i < Dimension; i++ {
		xPos := cellSize / 2.0
		row := make([]visualNode, Dimension)
		for j := 0; j < Dimension; j++ {
			row[j] = newNode(renderer, xPos, yPos)
			xPos += cellSize
		}
		output[i] = row
		yPos += cellSize

	}
	return output
}

func drawVisualNodes(visualNodes [][]visualNode, renderer *sdl.Renderer) {
	for _, nodeRow := range visualNodes {
		for _, node := range nodeRow {
			node.draw(renderer)
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"Hentai Pacman",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	defer renderer.Destroy()

	plr := newPlayer(renderer)

	visualNodes := getVisualNodes(renderer)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		plr.draw(renderer)
		drawVisualNodes(visualNodes, renderer)
		//nde.draw(renderer)

		renderer.Present()
	}
}
