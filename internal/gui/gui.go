package gui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tidbits/internal/db"
	"tidbits/internal/logger"
)

type GUI struct {
	log  *logger.Logger
	tbdb *db.TidbitsDB
}

func NewGUI(log *logger.Logger, tbdb *db.TidbitsDB) *GUI {
	return &GUI{
		log:  log,
		tbdb: tbdb,
	}
}

func (g *GUI) Run() {

	a := app.New()
	w := a.NewWindow("Tidbits")

	placeholder := widget.NewLabel("placeholder")
	temp2 := widget.NewLabel("temp2")

	mainContent := container.NewVScroll(placeholder)
	mainContent.Content = temp2

	dashBtn := widget.NewButton("Dashboard", func() {
		mainContent.Content = BuildDashboardBox()
	})
	sensorsBtn := widget.NewButton("Sensors", func() {
		mainContent.Content = BuildSensorsBox()
	})
	quitBtn := widget.NewButton("Quit", func() {
		temp2.SetText("NUEVO")
		a.Quit()
	})

	top := canvas.NewText("top bar", color.White)
	left := container.New(layout.NewVBoxLayout(), dashBtn, sensorsBtn, quitBtn)
	content := container.NewBorder(top, nil, left, nil, mainContent)
	//rawsensor.Wrapping = fyne.TextWrapBreak

	w.SetContent(container.New(layout.NewHBoxLayout(), content))
	w.ShowAndRun()
}
