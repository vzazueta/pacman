package main

import (
	"math/rand"
	"time"
)

/*
TILES:
0 -> WALKABLE
1 -> WALL

ENTITY:
EMPTY -> 0
PACMAN -> 1
GHOST -> 2

DIRECTION:
LEFT -> 0
UP -> 1
LEFT -> 2
DOWN -> 3
*/

// Node is a tile
type Node struct {
	tile   int
	entity int
	hasDot bool
	next   []Node
}

// Pacman is the player
type Pacman struct {
	currentNode Node
}

// Ghost is the enemy
type Ghost struct {
	currentNode Node
}

// Dimension of the gmae
const Dimension = 5

var layout = [Dimension][Dimension]int{
	{1, 1, 0, 0, 0},
	{0, 0, 0, 1, 0},
	{0, 1, 0, 1, 0},
	{0, 0, 0, 1, 0},
	{1, 1, 0, 0, 0}}

var nodeLayout [Dimension][Dimension]Node
var pacman Pacman
var ghosts []Ghost

func main() {
	rand.Seed(time.Now().UnixNano())

	createNodes()
	createPacman()
	createGhosts(3)
}

func createNodes() {
	for i := range layout {
		for j := range layout[i] {
			nodeLayout[i][j] = Node{tile: layout[i][j], hasDot: true, next: make([]Node, 0)}

			if i > 0 && nodeLayout[i-1][j].tile == 0 {
				nodeLayout[i-1][j].next = append(nodeLayout[i-1][j].next, nodeLayout[i][j])
			}

			if j > 0 && nodeLayout[i][j-1].tile == 0 {
				nodeLayout[i][j-1].next = append(nodeLayout[i][j-1].next, nodeLayout[i][j])
			}
		}
	}
}

func createPacman() {
	pacman = Pacman{currentNode: randomEmptyWalkableTile()}
	pacman.currentNode.entity = 1
}

func createGhosts(n int) {
	ghosts = make([]Ghost, n)

	for _, i := range ghosts {
		i = Ghost{currentNode: randomEmptyWalkableTile()}
		i.currentNode.entity = 2
	}
}

func randomEmptyWalkableTile() Node {
	node := Node{tile: -1}

	for node.tile != 0 || node.entity != 0 {
		node = nodeLayout[rand.Intn(Dimension)][rand.Intn(Dimension)]
	}

	return node
}
