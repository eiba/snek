package main

import (
	"github.com/awesome-gocui/gocui"
	"github.com/eiba/snake/autopilot"
	"github.com/eiba/snake/game"
	"github.com/eiba/snake/game/view"
	"github.com/eiba/snake/hamiltonian-cycle"
	"log"
	"math/rand"
	"time"
)

var (
	gui              *gocui.Gui
	r                = rand.New(rand.NewSource(time.Now().UnixNano()))
	Running          = true
	GameFinished     = false
	AutoPilotEnabled = false
	tickInterval     = 50 * time.Millisecond
	gameView         = view.Properties{"game", "snake", "", game.Position{}}
	positionMatrix   [][]game.Position
)

func main() {
	gui = initGUI()
	defer gui.Close()

	if err := game.initKeybindings(); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}
}

func initGUI() *gocui.Gui {
	gui, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	gui.Highlight = true
	gui.SelFgColor = gocui.ColorRed
	gui.SetManagerFunc(manageGame)
	return gui
}

func initGameView(maxX int, maxY int) (game.position, error) {
	gameViewPosition := calculateGameViewPosition(maxX, maxY)
	if v, err := gui.SetView(gameView.name, gameViewPosition.x0, gameViewPosition.x0, gameViewPosition.x1, gameViewPosition.y1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return gameViewPosition, err
		}
		v.Title = "snake"
		if _, err := gui.SetViewOnBottom(gameView.name); err != nil {
			return gameViewPosition, err
		}
		game.initPositionMatrix(gameViewPosition)
		return gameViewPosition, initGame()
	}
	return gameViewPosition, nil
}

func calculateGameViewPosition(maxX int, maxY int) game.position {
	defaultPosition := game.position{0, 0, maxX - 25, maxY - 1}

	if defaultPosition.x1%2 != 0 {
		defaultPosition.x1--
	}
	if (defaultPosition.x1/game.deltaX)%2 != 0 {
		defaultPosition.x1 = defaultPosition.x1 - game.deltaX
	}

	if defaultPosition.y1%2 != 0 {
		defaultPosition.y1--
	}
	if (defaultPosition.y1/game.deltaY)%2 != 0 {
		defaultPosition.y1 = defaultPosition.y1 - game.deltaY
	}
	return defaultPosition
}

func initGame() error {
	var err error
	game.snakeHead.position, err = view.setViewAtRandomPosition(game.snakeHead.viewName, positionMatrix, true)
	if err != nil {
		return err
	}
	game.foodView.position, err = view.setViewAtRandomPosition(game.foodView.name, positionMatrix, false)
	if err != nil {
		return err
	}
	go updateMovement()
	return nil
}

func manageGame(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	var err error
	gameView.position, err = initGameView(maxX, maxY)
	if err != nil {
		log.Panicln(err)
	}

	if err := game.initKeybindingsView(); err != nil {
		log.Panicln(err)
	}

	if err := view.initStatsView(); err != nil {
		log.Panicln(err)
	}

	if err := view.initLoadingView(); err != nil {
		log.Panicln(err)
	}

	if err := view.initPauseView(); err != nil {
		log.Panicln(err)
	}

	if err := view.initGameOverView(); err != nil {
		log.Panicln(err)
	}
	return nil
}

func updateMovement() {
	for {
		time.Sleep(tickInterval)
		if !Running {
			continue
		}
		gui.Update(func(gui *gocui.Gui) error {
			game.initPositionMatrix(gameView.position)
			if err := hamiltonian_cycle.initHamiltonianCycle(gameView.position); err != nil {
				log.Panicln(err)
			}
			if AutoPilotEnabled {
				err := autopilot.autopilot()
				if err != nil {
					log.Panicln(err)
				}
			}
			if err := game.movesnakeHead(); err != nil {
				log.Panicln(err)
			}
			if err := game.movesnakeBodyParts(); err != nil {
				log.Panicln(err)
			}
			return nil
		})
	}
}
