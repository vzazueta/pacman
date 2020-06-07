package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/eiannone/keyboard"
)

/*
TILES:
0 -> WALL
1 -> WALKABLE

ENTITY:
EMPTY -> 0
PACMAN -> 1
GHOST -> 2

VISUAL:
0 -> WALL
1 -> DOT WALKABLE
2 -> EMPTY WALKABLE
3 -> PACMAN
4 -> GHOST

DIRECTION:
NONE -> -1
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
	i      int
	j      int
	visual int
	next   []Node
}

// Pacman is the player
type Pacman struct {
	direction     int
	nextDirection int
	currentNode   *Node
}

func (p Pacman) walk() {
	temp := Node{}
	next := &temp

	switch p.nextDirection {
	case 0:
		if p.currentNode.j < Dimension-1 {
			next = &nodeLayout[p.currentNode.i][p.currentNode.j+1]
		}
	case 1:
		if p.currentNode.i > 0 {
			next = &nodeLayout[p.currentNode.i-1][p.currentNode.j]
		}
	case 2:
		if p.currentNode.j > 0 {
			next = &nodeLayout[p.currentNode.i][p.currentNode.j-1]
		}
	case 3:
		if p.currentNode.i < Dimension-1 {
			next = &nodeLayout[p.currentNode.i+1][p.currentNode.j]
		}
	}

	if next.tile != 1 {
		switch p.direction {
		case 0:
			if p.currentNode.j < Dimension-1 {
				next = &nodeLayout[p.currentNode.i][p.currentNode.j+1]
			}
		case 1:
			if p.currentNode.i > 0 {
				next = &nodeLayout[p.currentNode.i-1][p.currentNode.j]
			}
		case 2:
			if p.currentNode.j > 0 {
				next = &nodeLayout[p.currentNode.i][p.currentNode.j-1]
			}
		case 3:
			if p.currentNode.i < Dimension-1 {
				next = &nodeLayout[p.currentNode.i+1][p.currentNode.j]
			}
		}
	} else {
		p.direction = p.nextDirection
	}

	p.nextDirection = -1

	if next.tile != 1 {
		return
	}

	if next.entity == 2 {
		gameover()
	}

	p.currentNode.entity = 0
	p.currentNode.visual = 2

	p.currentNode = next
	p.currentNode.entity = 1
	p.currentNode.hasDot = false
	p.currentNode.visual = 3
}

func (p Pacman) move() {
	keyEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		printLayout()
		p.walk()
		time.Sleep(1000 * time.Millisecond)

		for len(keyEvents) != 0 {
			event := <-keyEvents

			switch event.Key {
			case keyboard.KeyArrowRight:
				p.nextDirection = 0
			case keyboard.KeyArrowUp:
				p.nextDirection = 1
			case keyboard.KeyArrowLeft:
				p.nextDirection = 2
			case keyboard.KeyArrowDown:
				p.nextDirection = 3
			default:
				p.nextDirection = -1
			}

			for len(keyEvents) != 0 {
				<-keyEvents
			}

			break
		}
	}
}

// Ghost is the enemy
type Ghost struct {
	currentNode *Node
}

func (g Ghost) walk(direction int) {
	temp := Node{}
	next := &temp

	switch direction {
	case 0:
		if g.currentNode.j < Dimension-1 {
			next = &nodeLayout[g.currentNode.i][g.currentNode.j+1]
		}
	case 1:
		if g.currentNode.i > 0 {
			next = &nodeLayout[g.currentNode.i-1][g.currentNode.j]
		}
	case 2:
		if g.currentNode.j > 0 {
			next = &nodeLayout[g.currentNode.i][g.currentNode.j-1]
		}
	case 3:
		if g.currentNode.i < Dimension-1 {
			next = &nodeLayout[g.currentNode.i+1][g.currentNode.j]
		}
	}

	if next.tile != 1 {
		return
	}

	if next.entity == 1 {
		gameover()
	}

	g.currentNode.entity = 0
	if g.currentNode.hasDot {
		g.currentNode.visual = 1
	} else {
		g.currentNode.visual = 2
	}

	g.currentNode = next
	g.currentNode.entity = 2
	g.currentNode.visual = 4

	g.moveToPacman()
}

func (g Ghost) moveToPacman() {
	g.walk(0)
}

// Dimension of the gmae
const Dimension = 5

var layout = [Dimension][Dimension]int{
	{0, 0, 1, 1, 1},
	{1, 1, 1, 0, 1},
	{1, 0, 1, 0, 1},
	{1, 1, 1, 0, 1},
	{0, 0, 1, 1, 1}}

var nodeLayout [Dimension][Dimension]Node
var pacman Pacman
var ghosts []Ghost

var done chan int

func main() {
	rand.Seed(time.Now().UnixNano())

	createNodes()
	createPacman()
	createGhosts(0)

	go pacman.move()

	for _, i := range ghosts {
		go i.moveToPacman()
	}

	<-done
}

func gameover() {
	done <- 0
	os.Exit(1)
}

func getAdjacentDirection(node1, node2 Node) int {
	if node1.j < node2.j {
		return 0
	} else if node1.i > node2.i {
		return 1
	} else if node1.j > node2.j {
		return 2
	} else {
		return 3
	}
}

func createNodes() {
	for i := range layout {
		for j := range layout[i] {
			nodeLayout[i][j] = Node{tile: layout[i][j], hasDot: true, i: i, j: j, visual: layout[i][j], next: make([]Node, 4)}

			if i > 0 {
				nodeLayout[i-1][j].next[3] = nodeLayout[i][j]
				nodeLayout[i][j].next[1] = nodeLayout[i-1][j]
			}

			if j > 0 {
				nodeLayout[i][j-1].next[0] = nodeLayout[i][j]
				nodeLayout[i][j].next[2] = nodeLayout[i][j-1]
			}
		}
	}
}

func createPacman() {
	pacman = Pacman{nextDirection: -1, currentNode: randomEmptyWalkableTile()}
	pacman.currentNode.entity = 1
	pacman.currentNode.hasDot = false
	pacman.currentNode.visual = 3
}

func createGhosts(n int) {
	ghosts = make([]Ghost, n)

	for _, i := range ghosts {
		i = Ghost{currentNode: randomEmptyWalkableTile()}
		i.currentNode.entity = 2
		i.currentNode.visual = 4
	}
}

func randomEmptyWalkableTile() *Node {
	temp := Node{}
	node := &temp

	for node.tile != 1 || node.entity != 0 {
		node = &nodeLayout[rand.Intn(Dimension)][rand.Intn(Dimension)]
	}

	return node
}

func printLayout() {
	for _, i := range nodeLayout {
		for _, j := range i {
			fmt.Printf("%d \t", j.visual)
		}
		fmt.Println()
	}
	fmt.Println()
}
