package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Tidbits")

	w.SetContent(widget.NewLabel("Hello World2"))
	w.ShowAndRun()
}
