package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	title      = "Conway's Game of Life"
	golSize    = 150
	scale      = 4
	padding    = int(scale / 2)
	windowSize = float64(golSize * scale)
	fillFactor = 0.5
	rng        = rand.New(rand.NewSource(time.Now().UnixNano()))
	paused     = false
)

// GameOfLife struct to hold the simulation state
type GameOfLife struct {
	// world represents the current world
	// nextWorld represents the world after the next update
	worldSize int
	world     [][]bool
	nextWorld [][]bool
}

// Initialize initializes the world with a random state
func (g *GameOfLife) Initialize(rng *rand.Rand, f float64) {
	for x, r := range g.world {
		for y := range r {
			g.world[x][y] = false
			if rng.Float64() < f {
				g.world[x][y] = true
			}
		}
	}
}

// NeighborsAlive returns the number of neighbors alive in the current world
func (g *GameOfLife) NeighborsAlive(x, y int) int {
	nalive := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i == x && j == y {
				continue
			}
			xloc, yloc := i, j
			if xloc < 0 {
				xloc = g.worldSize - 1
			}
			if xloc >= g.worldSize {
				xloc = 0
			}
			if yloc < 0 {
				yloc = g.worldSize - 1
			}
			if yloc >= g.worldSize {
				yloc = 0
			}
			if g.world[xloc][yloc] {
				nalive++
			}
		}
	}
	return nalive
}

// Update updates the current world to the next iteration
func (g *GameOfLife) Update() {
	for x, r := range g.world {
		for y, v := range r {
			na := g.NeighborsAlive(x, y)
			g.nextWorld[x][y] = false
			if v {
				if na == 2 || na == 3 {
					g.nextWorld[x][y] = true
				}
			} else if na == 3 {
				g.nextWorld[x][y] = true
			}
		}
	}
	g.world, g.nextWorld = g.nextWorld, g.world
}

// Draw draws the world onto an IMDraw
func (g *GameOfLife) Draw(imd *imdraw.IMDraw, s int) {
	imd.Clear()
	padding := float64(s / 2)
	for x, r := range g.world {
		for y, v := range r {
			if v {
				// nextWorld is now the previous state
				// Draw newly alive cells in color
				if !g.nextWorld[x][y] {
					imd.Color(colornames.Crimson)
				} else {
					imd.Color(colornames.White)
				}
				xloc := float64(x * s)
				yloc := float64(y * s)
				imd.Push(pixel.V(xloc-padding, yloc-padding))
				imd.Push(pixel.V(xloc+padding, yloc+padding))
				imd.Rectangle(0)
			}
		}
	}
}

// Reinitialize a 11x11 grid centered around the mouse click
func (g *GameOfLife) Reinitialize(rng *rand.Rand, xloc, yloc int) {
	padding := 5
	for x := xloc - padding; x < xloc+padding; x++ {
		for y := yloc - padding; y < yloc+padding; y++ {
			if x < 0 || y < 0 || x >= g.worldSize || y >= g.worldSize {
				continue
			}
			g.world[x][y] = false
			// Wipe out this pixel's history
			g.nextWorld[x][y] = false
			if rng.Float64() < 0.5 {
				g.world[x][y] = true
			}
		}
	}
}

// Clear clears the entire world
func (g *GameOfLife) Clear() {
	for x, r := range g.world {
		for y := range r {
			g.world[x][y] = false
			g.nextWorld[x][y] = false
		}
	}
}

// Set sets a location in the GOL world to alive or dead
func (g *GameOfLife) Set(xloc, yloc int, status bool) {
	g.world[xloc][yloc] = status
	// Delete history of the pixel
	g.nextWorld[xloc][yloc] = false
}

// NewGOL returns a new instance of GameOfLife
func NewGOL(ws int) *GameOfLife {
	gol := &GameOfLife{}
	gol.worldSize = ws
	gol.world = make([][]bool, ws)
	for i := range gol.world {
		gol.world[i] = make([]bool, ws)
	}
	gol.nextWorld = make([][]bool, ws)
	for i := range gol.nextWorld {
		gol.nextWorld[i] = make([]bool, ws)
	}
	return gol
}

// run is the main thread of the Pixel app
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, windowSize, windowSize),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	//imd.Color(colornames.Black)
	g := NewGOL(golSize)
	g.Initialize(rng, fillFactor)
	g.Draw(imd, scale)

	frames := 0
	secondsTicker := time.Tick(time.Second)
	for !win.Closed() {
		win.Clear(colornames.Black)
		if win.Pressed(pixelgl.MouseButtonLeft) {
			// Capture the Mouse Left Button press event and
			// set the appropriate pixels to alive
			if win.MousePosition().X() >= win.Bounds().Min.X() &&
				win.MousePosition().X() < win.Bounds().Max.X() &&
				win.MousePosition().Y() >= win.Bounds().Min.Y() &&
				win.MousePosition().Y() < win.Bounds().Max.Y() {
				xloc := int(win.MousePosition().X() / float64(scale))
				yloc := int(win.MousePosition().Y() / float64(scale))
				if win.Pressed(pixelgl.KeyLeftControl) ||
					win.Pressed(pixelgl.KeyRightControl) {
					g.Reinitialize(rng, xloc, yloc)
				} else {
					if win.Pressed(pixelgl.KeyLeftAlt) ||
						win.Pressed(pixelgl.KeyRightAlt) {
						g.Set(xloc, yloc, false)
					} else {
						g.Set(xloc, yloc, true)
					}
				}
			}
		}
		if win.JustPressed(pixelgl.KeyR) {
			g.Initialize(rng, fillFactor)
		}
		if win.JustReleased(pixelgl.KeySpace) {
			paused = !paused
		}
		if win.JustReleased(pixelgl.KeyC) {
			g.Clear()
		}
		if !paused {
			g.Update()
		}
		g.Draw(imd, scale)
		imd.Draw(win)
		win.Update()
		frames++
		select {
		case <-secondsTicker:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	profileFile, err := os.Create("cpuprofile.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(profileFile)
	defer pprof.StopCPUProfile()
	pixelgl.Run(run)
}
