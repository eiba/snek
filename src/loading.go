package main

import "github.com/awesome-gocui/gocui"

const loadingViewName = "loading"

var loadingView *gocui.View

func initLoadingView() error {
	lenX, lenY, err := getLenXY(gameView.name)
	if err != nil {
		return err
	}

	viewPositionX, viewPositionY := (lenX/2)-13, (lenY/2)-2
	viewLenX := 26
	viewLenY := 4

	loadingViewText := "Initiating autopilot..."
	loadingViewProps := viewProperties{
		loadingViewName,
		"Loading",
		loadingViewText,
		position{
			viewPositionX,
			viewPositionY,
			viewPositionX + viewLenX,
			viewPositionY + viewLenY}}
	loadingView, err = createView(loadingViewProps, false)
	return err
}

func loading(loading bool) error {
	if gameFinished && !running {
		return nil
	}
	loadingView.Visible = loading
	if loading {
		if _, err := gui.SetCurrentView(loadingViewName); err != nil {
			return err
		}
		if _, err := gui.SetViewOnTop(loadingViewName); err != nil {
			return err
		}
	}
	return nil
}
