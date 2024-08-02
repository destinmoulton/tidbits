package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"tidbits/internal/cliflags"
	"tidbits/internal/db"
	"tidbits/internal/logger"
	"tidbits/internal/sensors"
)

func main() {
	flags := cliflags.ParseFlags()
	log := logger.NewLogger(flags)
	defer log.Close()

	tbdb, err := db.NewTidbitsDB(log)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}
	defer tbdb.Close()

	tbdb.Init()

	a := app.New()
	w := a.NewWindow("Tidbits")

	sensdata, err := sensors.GetSensorReadings()
	if err != nil {

	}
	button := widget.NewButton("Sensors", func() {
		log.Info("clicked button")
	})
	menu := container.New(layout.NewVBoxLayout(), button)
	rawsensor := widget.NewLabel(sensdata.Fulltext)
	//rawsensor.Wrapping = fyne.TextWrapBreak
	vscroll := container.NewVScroll(rawsensor)

	w.SetContent(container.New(layout.NewHBoxLayout(), menu, vscroll))
	w.ShowAndRun()
}
