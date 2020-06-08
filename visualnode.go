package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	nodeSize = 50
)

type visualNode struct {
	tex  *sdl.Texture
	x, y float64
}

func (n *visualNode) getCoin(renderer *sdl.Renderer){
	n.tex = textureFromBMP(renderer, "sprites/empty.bmp")
}

func newNode(renderer *sdl.Renderer, xcoord float64, ycoord float64) (n visualNode) {

	n.tex = textureFromBMP(renderer, "sprites/node.bmp")
	n.x = xcoord
	n.y = ycoord

	return n
}

func (n *visualNode) draw(renderer *sdl.Renderer) {
	x := n.x - nodeSize/2.0
	y := n.y - nodeSize/2.0
	renderer.Copy(n.tex,
		&sdl.Rect{X: 0, Y: 0, W: nodeSize, H: nodeSize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: nodeSize, H: nodeSize})
}
