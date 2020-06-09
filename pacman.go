package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
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
	next   []*Node
}

// Pacman is the player
type Pacman struct {
	direction     int
	nextDirection int
	currentNode   *Node
}

func (p *Pacman) walk() {
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
	//visualNodes[p.currentNode.i][p.currentNode.j].updateTex(renderer, p.currentNode.visual)

	p.currentNode = next
	p.currentNode.entity = 1
	p.currentNode.hasDot = false
	p.currentNode.visual = 3
	//visualNodes[p.currentNode.i][p.currentNode.j].updateTex(renderer, p.currentNode.visual)
}

func (p *Pacman) move() {
	//visualNodes[p.currentNode.i][p.currentNode.j].updateTex(renderer, p.currentNode.visual)
	for {
		//printLayout()
		p.walk()
		time.Sleep(200 * time.Millisecond)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				//fmt.Printf("type:%d\tsym:%d\n", t.Type, t.Keysym.Sym)
				switch t.Keysym.Sym {
				case ArrowDown:
					//fmt.Println("down")
					p.nextDirection = 3
				case ArrowUp:
					//fmt.Println("up")
					p.nextDirection = 1
				case ArrowLeft:
					//fmt.Println("left")
					p.nextDirection = 2
				case ArrowRight:
					//fmt.Println("right")
					p.nextDirection = 0
				}

			case *sdl.QuitEvent:
				fmt.Println("finish")
				done <- 0
				os.Exit(1)
			}
		}
	}

	/*for {

	}*/

	//printLayout()
	//p.walk()
	//time.Sleep(200 * time.Millisecond)

	/*for {
		keys := sdl.GetKeyboardState()

		if keys[sdl.SCANCODE_LEFT] == 1 {
			fmt.Println("left")
			p.nextDirection = 2
		} else if keys[sdl.SCANCODE_RIGHT] == 1 {
			fmt.Println("right")
			p.nextDirection = 0
		} else if keys[sdl.SCANCODE_UP] == 1 {
			fmt.Println("up")
			p.nextDirection = 1
		} else if keys[sdl.SCANCODE_DOWN] == 1 {
			fmt.Println("down")
			p.nextDirection = 3
		}
		continue
	}*/
}

// Ghost is the enemy
type Ghost struct {
	currentNode *Node
	seen        map[*Node]bool
	path        *[]*Node
}

func (g *Ghost) walk(direction int) {

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

	//fmt.Println("at walk: ", g)
	//visualNodes[g.currentNode.i][g.currentNode.j].updateTex(renderer, g.currentNode.visual)
	//g.moveToPacman()
}

func (g *Ghost) moveToPacman() {
	if len(*g.path) < 1 {

		g.seen = make(map[*Node]bool, 0)
		g.seen[g.currentNode] = true

		tmp := make([]*Node, 0)
		g.getNextNode(g.currentNode, &tmp)
		fmt.Println("the length: ", len(*g.path))
		for _, k := range *g.path {
			fmt.Printf("uwu: %d, %d\n", k.i, k.j)
		}
		//g.moveToPacman()
	} else {
		g.walk(getAdjacentDirection(g.currentNode, (*g.path)[0]))
		fmt.Println(len(*g.path))
		tmp := (*g.path)[1:]
		g.path = &tmp
	}

	time.Sleep(200 * time.Millisecond)
	g.moveToPacman()
}

func (g *Ghost) getNextNode(node *Node, path *[]*Node) *[]*Node {
	//fmt.Printf("%d, %d  %d\n", node.i, node.j, len(path))

	if node.entity == 1 {
		//fmt.Printf("%d, %d  %d\n", node.i, node.j, len(*path))
		return path
	}

	for _, i := range node.next {
		if i.tile == 1 {
			if g.seen[i] {
				continue
			}

			g.seen[i] = true
			//fmt.Println("i:", i.i)
			//fmt.Println("yes: ", i.j)
			tmp2 := append(*path, i)
			temp := g.getNextNode(i, &tmp2)
			//fmt.Println(temp.i, temp.j)
			if len(*temp) > 0 && (*temp)[len(*temp)-1].entity == 1 {
				g.path = temp
				return path
			}
		}
	}

	tmp := make([]*Node, 0)
	return &tmp
}

// Dimension of the gmae
const (
	Dimension  = 25
	ArrowDown  = 1073741905
	ArrowUp    = 1073741906
	ArrowLeft  = 1073741904
	ArrowRight = 1073741903
)

var layout = [][]int{
	{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	{0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1},
	{1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1},
	{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}

var nodeLayout [][]Node
var pacman Pacman
var ghosts []*Ghost
var visualNodes [][]visualNode
var renderer *sdl.Renderer

var done chan int

func main() {
	rand.Seed(time.Now().UnixNano())

	createNodes()
	createPacman()
	createGhosts(2)

	for _, i := range ghosts {
		//fmt.Println("at main:", i)
		go i.moveToPacman()
	}

	<-done
}

func gameover() {
	fmt.Println("you dead")
	done <- 0
	os.Exit(1)
}

func getAdjacentDirection(node1, node2 *Node) int {
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
	nodeLayout = make([][]Node, Dimension)
	for i := range layout {
		nodeLayout[i] = make([]Node, Dimension)
		for j := range layout[i] {
			nodeLayout[i][j] = Node{tile: layout[i][j], hasDot: true, i: i, j: j, visual: layout[i][j], next: make([]*Node, 4)}
			temp := Node{}
			for k := range nodeLayout[i][j].next {
				nodeLayout[i][j].next[k] = &temp
			}

			if i > 0 {
				nodeLayout[i-1][j].next[3] = &nodeLayout[i][j]
				nodeLayout[i][j].next[1] = &nodeLayout[i-1][j]
			}

			if j > 0 {
				nodeLayout[i][j-1].next[0] = &nodeLayout[i][j]
				nodeLayout[i][j].next[2] = &nodeLayout[i][j-1]
			}
		}
	}
	go visualSetup()
}

func createPacman() {
	pacman = Pacman{nextDirection: -1, currentNode: randomEmptyWalkableTile()}
	pacman.currentNode.entity = 1
	pacman.currentNode.hasDot = false
	pacman.currentNode.visual = 3
}

func createGhosts(n int) {
	ghosts = make([]*Ghost, n)

	for i := range ghosts {
		tmp := make([]*Node, 0)
		ghosts[i] = &Ghost{currentNode: randomEmptyWalkableTile(), path: &tmp}
		ghosts[i].currentNode.entity = 2
		ghosts[i].currentNode.visual = 4
		//visualNodes[ghosts[i].currentNode.i][ghosts[i].currentNode.j].updateTex(renderer, ghosts[i].currentNode.visual)
	}

	//fmt.Println("lal: ", ghosts)
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

func visualSetup() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err := sdl.CreateWindow(
		"Pacman",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	//fmt.Println("renderer created")
	defer renderer.Destroy()

	initTex(renderer)

	//plr = newPlayer(renderer)

	visualNodes = getVisualNodes(renderer, layout)
	/*for _, i := range ghosts {
		visualNodes[i.currentNode.i][i.currentNode.j].updateTex(renderer, i.currentNode.visual)
	}*/
	go pacman.move()

	for {
		//fmt.Println("uwu")
		/*for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {

		}*/

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		//plr.draw(renderer)
		drawVisualNodes(visualNodes, renderer, nodeLayout)

		renderer.Present()
	}
}
