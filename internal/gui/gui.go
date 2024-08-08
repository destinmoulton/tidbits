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
	GUIMargin          = 5
	GUIWidthFudge      = 20
	GUIHeightFudge     = 30
	GUIMenuWidth       = 200
	GUIMinHeight       = 200
	GUIMinWidth        = 200
	GUIMessengerHeight = 200
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
	g.content.SetMinSize(fyne.NewSize(GUIMinWidth, GUIMinHeight))
	g.messenger.SetMinSize(fyne.NewSize(size.Width, GUIMessengerHeight))

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
		g.content.Content = g.dashboardView()
	case GViewSensors:
		g.content.Content = g.sensorsView()
	}
	g.content.SetMinSize(fyne.NewSize(g.calcContentWidth(), g.calcContentHeight()))
	g.content.Refresh()
}

func (g *GUI) calcContentWidth() float32 {
	size := g.window.Content().Size()
	return size.Width - GUIMenuWidth - GUIWidthFudge
}

func (g *GUI) calcContentHeight() float32 {
	size := g.window.Content().Size()
	return size.Height - GUIMessengerHeight - GUIHeightFudge
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

func (g *GUI) msg(parts ...interface{}) {
	var builder strings.Builder
	for _, part := range parts {
		builder.WriteString(fmt.Sprintf("%v", part))
	}
	g.messenger.Content = widget.NewLabel(builder.String())
	g.messenger.Refresh()
}
