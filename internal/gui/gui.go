package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"tidbits/internal/db"
	"tidbits/internal/logger"
)

type GUIViewID int

const (
	GViewDashoard GUIViewID = iota
	GViewSensors
)

type GUITabID int

const (
	GTabDefault GUITabID = iota
	GTabLMSensorsSelectForm
	GTabLMSensorsUserForm
)

type GUI struct {
	log     *logger.Logger
	tbdb    *db.TidbitsDB
	app     fyne.App
	window  fyne.Window
	content fyne.CanvasObject
	view    GUIViewID
	subtab  GUITabID
}

func NewGUI(log *logger.Logger, tbdb *db.TidbitsDB) *GUI {
	tmp := widget.NewLabel("")
	return &GUI{
		log:     log,
		tbdb:    tbdb,
		content: tmp,
		view:    GViewDashoard,
		subtab:  GTabDefault,
	}
}

func (g *GUI) Run() {

	g.app = app.New()
	g.window = g.app.NewWindow("Tidbits")
	g.render()
	g.window.ShowAndRun()
}

func (g *GUI) switchTab(tab GUITabID) {
	g.subtab = tab
	g.render()
}

func (g *GUI) switchView(view GUIViewID) {
	g.view = view
	g.subtab = GTabDefault
	g.render()
}

func (g *GUI) render() {

	placeholder := widget.NewLabel("placeholder")
	temp2 := widget.NewLabel("temp2")

	mainContent := container.NewVScroll(placeholder)
	mainContent.Content = g.content

	dashBtn := widget.NewButton("Dashboard", func() {
		g.switchView(GViewDashoard)
	})
	sensorsBtn := widget.NewButton("Sensors", func() {
		g.switchView(GViewSensors)
	})
	quitBtn := widget.NewButton("Quit", func() {
		temp2.SetText("NUEVO")
		g.app.Quit()
	})

	switch g.view {
	case GViewDashoard:
		mainContent.Content = BuildDashboardBox()
	case GViewSensors:
		mainContent.Content = g.sensorsView()
	}

	top := canvas.NewText("top bar", color.White)
	left := container.New(layout.NewVBoxLayout(), dashBtn, sensorsBtn, quitBtn)
	content := container.NewBorder(top, nil, left, nil, mainContent)

	g.window.SetContent(container.New(layout.NewHBoxLayout(), content))
}
