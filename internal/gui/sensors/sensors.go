package sensors

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"tidbits/internal/sensors"
)

func BuildSensorsBox() *fyne.Container {

	placeholder := widget.NewLabel("sensors")
	content := container.NewVScroll(placeholder)
	wrapper := container.NewVScroll(content)
	lmsensorsBtn := widget.NewButton("Choose lm_sensors", func() {
		sensdata, err := sensors.GetSensorReadings()
		if err != nil {

		}
		rawtext := widget.NewLabel(sensdata.Fulltext)
		content.Content = rawtext
	})

	menubar := container.New(layout.NewHBoxLayout(), lmsensorsBtn)

	return container.NewBorder(menubar, nil, nil, nil, wrapper)
}
