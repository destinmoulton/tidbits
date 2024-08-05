package main

import (
	"embed"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"tidbits/internal/cliflags"
	"tidbits/internal/db"
	"tidbits/internal/logger"
	"tidbits/internal/sensors"
)

// Use go embed to put them
// the db migrationsFS in the binary

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

func main() {
	flags := cliflags.ParseFlags()
	log := logger.NewLogger(flags)
	defer log.Close()

	tbdb, err := db.NewTidbitsDB(log, migrationsFS)
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
