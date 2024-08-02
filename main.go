package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"tidbits/internal/sensors"
)

func main() {
	a := app.New()
	w := a.NewWindow("Tidbits")

	sensdata, err := sensors.GetSensorReadings()
	if err != nil {

	}
	button := widget.NewButton("Sensors", func() {
		log.Println("clicked")
	})
	menu := container.New(layout.NewVBoxLayout(), button)
	rawsensor := widget.NewLabel(sensdata.Fulltext)
	//rawsensor.Wrapping = fyne.TextWrapBreak
	vscroll := container.NewVScroll(rawsensor)

	w.SetContent(container.New(layout.NewHBoxLayout(), menu, vscroll))
	w.ShowAndRun()
}
