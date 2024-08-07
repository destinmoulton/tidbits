package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
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

const (
	GUIMenuWidth = 200
	GUIMinHeight
	GUIWidthOffset = 20
	GUIMargin      = 5
)

type GUI struct {
	log       *logger.Logger
	tbdb      *db.TidbitsDB
	app       fyne.App
	window    fyne.Window
	view      GUIViewID
	subtab    GUITabID
	messages  []string
	messenger *container.Scroll
	content   *container.Scroll
}

func NewGUI(log *logger.Logger, tbdb *db.TidbitsDB) *GUI {
	return &GUI{
		log:    log,
		tbdb:   tbdb,
		view:   GViewDashoard,
		subtab: GTabDefault,
	}
}

func (g *GUI) Run() {

	g.app = app.New()
	g.window = g.app.NewWindow("Tidbits")

	placeholder := widget.NewLabel("")
	blank := widget.NewLabel("")

	g.content = container.NewVScroll(placeholder)
	g.content.Content = blank

	top := canvas.NewText("top bar", color.White)
	leftMenu := g.buildLeftMenu()
	g.messenger = container.NewVScroll(widget.NewLabel("messages"))
	content := container.NewBorder(top, g.messenger, leftMenu, nil, g.content)
	g.window.SetContent(container.New(layout.NewHBoxLayout(), content))

	size := g.window.Content().Size()
	leftMenu.SetMinSize(fyne.NewSize(GUIMenuWidth, GUIMinHeight))
	g.content.SetMinSize(fyne.NewSize(size.Width-GUIMenuWidth, GUIMinHeight))
	g.messenger.SetMinSize(fyne.NewSize(size.Width, GUIMinHeight))

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
	switch g.view {
	case GViewDashoard:
		g.content.Content = BuildDashboardBox()
	case GViewSensors:
		g.content.Content = g.sensorsView()
	}
	size := g.window.Content().Size()
	g.content.SetMinSize(fyne.NewSize(size.Width-GUIMenuWidth-GUIWidthOffset, GUIMinHeight))
	g.content.Refresh()
}

func (g *GUI) buildLeftMenu() *container.Scroll {

	dashBtn := widget.NewButton("Dashboard", func() {
		g.switchView(GViewDashoard)
	})
	sensorsBtn := widget.NewButton("Sensors", func() {
		g.switchView(GViewSensors)
	})
	quitBtn := widget.NewButton("Quit", func() {
		g.app.Quit()
	})
	cont := container.New(layout.NewVBoxLayout(), dashBtn, sensorsBtn, quitBtn)
	return container.NewVScroll(cont)
}

func (g *GUI) msg(msg string, parts ...string) {
	if len(parts) > 0 {
		tmp := strings.Join(parts, "")
		g.messages = append(g.messages, msg+tmp)
	} else {
		g.messages = append(g.messages, msg)
	}

	output := ""
	for i := len(g.messages) - 1; i >= 0; i-- {
		output += fmt.Sprintf("%s\n", g.messages[i])
	}

	g.messenger.Content = widget.NewLabel(output)
	g.messenger.Refresh()
}
