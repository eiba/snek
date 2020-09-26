package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
	"math/rand"
	"time"
)

const delta = 1

type snekBodyPart struct {
	currentDirection  direction
	previousDirection direction
	viewName string
}

type direction int
type movementDirections struct {
	up    direction
	right direction
	down  direction
	left  direction
}

var (
	directions        = movementDirections{0, 1, 2, 3}
	snekBodyParts     = []snekBodyPart{{directions.up, directions.up, "s0"}}
	headDirection     = directions.up
	gameView, boxView = "game", "box"
	running           = true
	tickInterval      = 100 * time.Millisecond
	r                 = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func main() {
	run()
}

func run() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Highlight = true
	g.SelFgColor = gocui.ColorRed
	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}
}

func getOppositeDirection(direction direction) direction  {
	return (direction + 2) % 4
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("help", maxX-25, 0, maxX-1, 6, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "Keybindings"
		fmt.Fprintln(v, "Space: Restart")
		fmt.Fprintln(v, "← ↑ → ↓: Move thing")
		fmt.Fprintln(v, "Ctrl+W: Speed up")
		fmt.Fprintln(v, "Ctrl+S: Slow down")
		fmt.Fprintln(v, "^C: Exit")
	}

	if v, err := g.SetView(gameView, 0, 0, maxX-26, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		if _, err := g.SetViewOnBottom(gameView); err != nil {
			return err
		}
		if err := setViewAtRandom(g, snekBodyParts[0].viewName, true); err != nil {
			log.Panicln(err)
		}
		go updateMovement(g)
		if err := setViewAtRandom(g, boxView, false); err != nil {
			log.Panicln(err)
		}
		v.Title = "Snek"
	}

	return nil
}

func updateMovement(g *gocui.Gui) error {
	for {
		if !running {
			continue
		}

		snekBodyParts[0].previousDirection = snekBodyParts[0].currentDirection
		snekBodyParts[0].currentDirection = headDirection
		err := moveViewInDirection(g, snekBodyParts[0].viewName, snekBodyParts[0].currentDirection, true); if err != nil { return err}

		for i := 1; i < len(snekBodyParts); i++ {
			currentsSnekBodyPart := snekBodyParts[i]
			previousSnekBodyPartPreviousDirection := snekBodyParts[i-1].previousDirection
			err := moveViewInDirection(g, currentsSnekBodyPart.viewName, previousSnekBodyPartPreviousDirection, false)
			if err != nil {
				return err
			}
			snekBodyParts[i].previousDirection = snekBodyParts[i].currentDirection
			snekBodyParts[i].currentDirection = previousSnekBodyPartPreviousDirection
		}

		time.Sleep(tickInterval)
	}
}

func moveViewInDirection(g *gocui.Gui, viewName string, direction direction, headView bool) error {
	g.Update(func(g *gocui.Gui) error {
		var err error
		switch direction {
			case directions.up: //up
				err = moveView(g, viewName, 0, -delta, headView)
			case directions.right: //right
				err = moveView(g, viewName, delta+1, 0, headView)
			case directions.down: //down
				err = moveView(g, viewName, 0, delta, headView)
			case directions.left: //left
				err = moveView(g, viewName, -delta-1, 0, headView)
		}
		return err
	})
	return nil
}

func reset(g *gocui.Gui) error {
	headDirection = 0
	running = true
	tickInterval = 100 * time.Millisecond
	for i := 1; i < len(snekBodyParts); i++ {
		if err := g.DeleteView(snekBodyParts[i].viewName); err != nil && !gocui.IsUnknownView(err) {
			return err
		}
	}
	snekBodyParts = []snekBodyPart{{0, 0, "s0"}}

	if err := setViewAtRandom(g, snekBodyParts[0].viewName, true); err != nil {
		return err
	}
	if err := setViewAtRandom(g, boxView, false); err != nil {
		return err
	}
	if err := g.DeleteView("gameOver"); err != nil && !gocui.IsUnknownView(err) {
		return err
	}

	return nil
}

func gameOver(g *gocui.Gui) error {
	running = false
	x0, y0, x1, y1, err := g.ViewPosition(gameView)
	if err != nil {
		return err
	}
	maxX, maxY := x1-x0, y1-y0

	positionX, positionY := (maxX/2)-5, (maxY/2)-2

	lenX := 12
	lenY := 4
	name := "gameOver"
	if v, err := g.SetView(name, positionX, positionY, positionX+lenX, positionY+lenY, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = "Game over"
		fmt.Fprintln(v, "\n  u lose")

		if _, err := g.SetCurrentView(name); err != nil {
			return err
		}
		if _, err := g.SetViewOnTop(name); err != nil {
			return err
		}
	}
	return nil
}

func setViewAtRandom(g *gocui.Gui, name string, setCurrent bool) error {
	x0, y0, x1, y1, err := g.ViewPosition(gameView)
	if err != nil {
		return err
	}

	maxX, maxY := x1-x0-3, y1-y0-2

	positionX, positionY := r.Intn(maxX)+1, r.Intn(maxY)+1

	lenX := 2
	lenY := 1
	_, err = g.SetView(name, positionX, positionY, positionX+lenX, positionY+lenY, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
	}

	if setCurrent {
		if _, err := g.SetCurrentView(name); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func addView(g *gocui.Gui, viewName string, direction direction) error {

	x0, y0, x1, y1, err := g.ViewPosition(viewName)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("s%v", len(snekBodyParts))

	offsetX := 0
	offsetY := 1
	switch direction {
	case directions.right: //right
		offsetX = -2
		offsetY = 0
	case directions.down: //down
		offsetX = 0
		offsetY = -1
	case directions.left: //left
		offsetX = 2
		offsetY = 0
	}

	_, err = g.SetView(name, x0+offsetX, y0+offsetY, x1+offsetX, y1+offsetY, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
	}
	currentLastSnekBodyPart := snekBodyParts[len(snekBodyParts)-1]
	snekBodyParts = append(snekBodyParts, snekBodyPart{currentLastSnekBodyPart.currentDirection, currentLastSnekBodyPart.currentDirection,name})

	return nil
}

//Checks collision between view1 and view2, returning true for collision and false otherwise.
func checkCollision(g *gocui.Gui, view1 string, view2 string) (bool, error) {
	x10, y10, x11, y11, err := g.ViewPosition(view1)
	if err != nil {
		return false, err
	}

	x20, y20, x21, y21, err := g.ViewPosition(view2)
	if err != nil {
		return false, err
	}

	Ax, Ay, Aw, Ah := x10, y10, x11-x10, y11-y10
	Bx, By, Bw, Bh := x20, y20, x21-x20, y21-y20

	if Bx+Bw > Ax &&
		By+Bh > Ay &&
		Ax+Aw > Bx &&
		Ay+Ah > By {
		return true, nil

	}

	return false, nil
}

func moveView(g *gocui.Gui, viewName string, dx, dy int, headView bool) error {
	if headView {
		return moveHead(g, viewName, dx, dy)
	}
	return moveBodyPart(g, viewName, dx,dy)
}

func checkBoxCollision(g *gocui.Gui, viewName string)  {

}

func moveHead(g *gocui.Gui, viewName string, dx, dy int) error  {
	xg0, yg0, xg1, yg1, err := g.ViewPosition(gameView); if err != nil {return err}
	
	maxX, maxY, minX, minY := xg1-xg0, yg1-yg0, 0, 0
	newX0, newY0, newX1, newY1, err:= getNewViewCoordinates(g, viewName,dx,dy); if err != nil { return err}

	if newX0 >= minX && newY0 >= minY && newX1 <= maxX && newY1 <= maxY {
		if _, err := g.SetView(viewName, newX0, newY0, newX1, newY1, 0); err != nil {
			return err
		}

		collision, err := checkCollision(g, snekBodyParts[0].viewName, boxView)
		if err != nil {
			return err
		}

		if collision {
			if err = addView(g, snekBodyParts[len(snekBodyParts)-1].viewName, snekBodyParts[len(snekBodyParts)-1].currentDirection); err != nil {
				return err
			}
			return setViewAtRandom(g, boxView, false)
		}
	} else {
		return gameOver(g)
	}

	return nil
}

func moveBodyPart(g *gocui.Gui, viewName string, dx, dy int) error {
	newX0, newY0, newX1, newY1, err:= getNewViewCoordinates(g, viewName,dx,dy); if err != nil { return err}
	_, err = g.SetView(viewName, newX0, newY0, newX1, newY1, 0); if err != nil {return err}
	return nil
}

func getNewViewCoordinates(g *gocui.Gui, viewName string, dx, dy int) (int, int, int, int, error)  {
	x0, y0, x1, y1, err := g.ViewPosition(viewName); if err != nil {
		return 0, 0, 0, 0, err
	}
	newX0, newY0, newX1, newY1 := x0+dx, y0+dy, x1+dx, y1+dy
	return newX0, newY0, newX1, newY1, nil
}