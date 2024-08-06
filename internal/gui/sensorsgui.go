package gui

import (
	"fmt"
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
		content.Content = lmSensorsForm()
	})

	menubar := container.New(layout.NewHBoxLayout(), lmsensorsBtn)

	return container.NewBorder(menubar, nil, nil, nil, wrapper)
}

type selectedSensor struct {
	sensorName sensors.SensorName
	sensorType sensors.SensorType
	deviceName sensors.DeviceName
}
type selectedSensors map[sensors.SensorName]selectedSensor

func lmSensorsForm() *widget.Form {
	devices, readings := sensors.LMSensorsParseReadings()

	var rows []*widget.FormItem
	var selected selectedSensors = make(map[sensors.SensorName]selectedSensor)
	for _, device := range devices {
		// blank line
		rows = append(rows, &widget.FormItem{
			Text:   "",
			Widget: widget.NewLabel(""),
		})
		// device name
		rows = append(rows, &widget.FormItem{
			Text:   string(device),
			Widget: widget.NewLabel(""),
		})
		for _, sensor := range readings[device] {
			if sensor.Source == "lm_sensors" {
				text := fmt.Sprintf("%s: %s", sensor.Name, sensor.Format())
				rows = append(rows, &widget.FormItem{
					Text: text,
					Widget: widget.NewCheck("", func(value bool) {
						if value {
							selected[sensor.Name] = selectedSensor{
								sensorName: sensor.Name,
								sensorType: sensor.Type,
								deviceName: sensor.Device,
							}
						} else {
							delete(selected, sensor.Name)
						}
						fmt.Println("checked " + string(sensor.Name))
					}),
				})

			}
		}
	}
	return &widget.Form{
		Items: rows,
		OnSubmit: func() { // optional, handle form submission
			fmt.Println(selected)
		},
	}

}
