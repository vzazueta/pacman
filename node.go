package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	nodeSize = 50
)

type node struct {
	tex  *sdl.Texture
	x, y float64
}

func newNode(renderer *sdl.Renderer) (n node) {

	n.tex = textureFromBMP(renderer, "sprites/node.bmp")
	n.x = screenWidth / 2.0
	n.y = screenHeight / 2.0

	return n
}

func (n *node) draw(renderer *sdl.Renderer) {
	x := n.x - nodeSize/2.0
	y := n.y - nodeSize/2.0
	renderer.Copy(n.tex,
		&sdl.Rect{X: 0, Y: 0, W: nodeSize, H: nodeSize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: nodeSize, H: nodeSize})
}
