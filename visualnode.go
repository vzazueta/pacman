package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// relacionar numero de estado con sprite

const (
	nodeSize = 50
)

type visualNode struct {
	tex  *sdl.Texture
	x, y float64
}

func newNode(renderer *sdl.Renderer, xcoord float64, ycoord float64) (n visualNode) {

	n.tex = textureFromBMP(renderer, "sprites/node.bmp")
	n.x = xcoord
	n.y = ycoord

	return n
}

func changeTex() {

}

func (n *visualNode) draw(renderer *sdl.Renderer) {
	x := n.x - nodeSize/2.0
	y := n.y - nodeSize/2.0
	renderer.Copy(n.tex,
		&sdl.Rect{X: 0, Y: 0, W: nodeSize, H: nodeSize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: nodeSize, H: nodeSize})
}
