package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type node struct {
	tex *sdl.Texture
}

func newNode(renderer *sdl.Renderer) (n node, err error) {
	img, err := sdl.LoadBMP("sprites/node.bmp")
	if err != nil {
		return node{}, fmt.Errorf("loading node sprite: %v", err)
	}
	defer img.Free()

	n.tex, err = renderer.CreateTextureFromSurface(img)

	if err != nil {
		return node{}, fmt.Errorf("creating node texture: %v", err)
	}

	return n, nil

}

func (n *node) draw(renderer *sdl.Renderer) {
	renderer.Copy(n.tex,
		&sdl.Rect{X: 0, Y: 0, W: 105, H: 105},
		&sdl.Rect{X: 150, Y: 20, W: 105, H: 105})
}
