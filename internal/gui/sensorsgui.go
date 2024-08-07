package gui

import (
	"context"
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"tidbits/internal/db/sqlc"
	"tidbits/internal/sensors"
)

func (g *GUI) sensorsView() *fyne.Container {

	placeholder := widget.NewLabel("sensors")
	content := container.NewVScroll(placeholder)
	wrapper := container.NewVScroll(content)
	lmsensorsBtn := widget.NewButton("Choose lm_sensors", func() {
		g.switchTab(GTabLMSensorsSelectForm)
	})

	switch g.subtab {
	case GTabDefault:
		content.Content = widget.NewLabel("SENSORS SECTION")
	case GTabLMSensorsSelectForm:
		content.Content = g.lmRenderSensorsForm()
	}

	menubar := container.New(layout.NewHBoxLayout(), lmsensorsBtn)

	return container.NewBorder(menubar, nil, nil, nil, wrapper)
}

type selectedSensor struct {
	sensorName sensors.SensorName
	sensorType sensors.SensorType
	deviceName sensors.DeviceName
}

type selectedSensors map[string]selectedSensor

func (g *GUI) lmRenderSensorsForm() *widget.Form {
	ctx := context.Background()
	devices, readings := sensors.LMSensorsParseReadings()

	storedSensors, err := g.tbdb.Queries.GetSensorsBySource(ctx, sql.NullString{String: "lm_sensors", Valid: true})
	fmt.Println(storedSensors)
	if err != nil {
		return nil
	}

	var rows []*widget.FormItem
	var selected selectedSensors = make(map[string]selectedSensor)
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
				id := string(sensor.Name) + ":" + string(sensor.Device)
				data := selectedSensor{
					sensorName: sensor.Name,
					sensorType: sensor.Type,
					deviceName: sensor.Device,
				}
				checkbox := widget.NewCheck("", func(value bool) {
					if value {
						selected[id] = data
					} else {
						delete(selected, id)
					}
					fmt.Println("checked event " + string(id))
				})

				// Find already selected sensors
				for _, storedSensor := range storedSensors {
					doesNameMatch := storedSensor.SensorName.String == string(sensor.Name)
					doesDeviceMatch := storedSensor.SensorDevice.String == string(sensor.Device)
					if doesNameMatch && doesDeviceMatch {
						selected[id] = data
						checkbox.Checked = true
					}
				}

				formItem := &widget.FormItem{
					Text:   text,
					Widget: checkbox,
				}

				rows = append(rows, formItem)
			}
		}
	}

	return &widget.Form{
		Items: rows,
		OnSubmit: func() { // optional, handle form submission
			g.saveLmSensors(storedSensors, selected)
			g.switchTab(GTabDefault)
		},
	}

}

func (g *GUI) saveLmSensors(storedSettings []sqlc.Sensor, chkdSensors selectedSensors) {
	ctx := context.Background()

	// start assuming all the form sensors need to be added
	sensorsToAdd := make(map[string]selectedSensor)
	for id, sensor := range chkdSensors {
		sensorsToAdd[id] = sensor
	}

	fmt.Println("chkdSensors pre first loop")
	fmt.Println(chkdSensors)
	// don't add any sensors that are already in the database
	for id, selSensor := range chkdSensors {
		for _, dbSensor := range storedSettings {
			doesNameMatch := dbSensor.SensorName.String == string(selSensor.sensorName)
			doesDeviceMatch := dbSensor.SensorDevice.String == string(selSensor.deviceName)

			if doesNameMatch && doesDeviceMatch {
				delete(sensorsToAdd, id)
			}
		}
	}

	sensorsToRemove := make(map[string]string)
	for _, dbSensor := range storedSettings {
		id := dbSensor.SensorName.String + ":" + dbSensor.SensorDevice.String
		sensorsToRemove[id] = string(dbSensor.ID)
	}

	fmt.Println("chkdSensors")
	fmt.Println(chkdSensors)
	fmt.Println("storedSettings")
	fmt.Println(storedSettings)
	for _, dbSensor := range storedSettings {
		for id, selSensor := range chkdSensors {
			fmt.Println("checking removal: " + id)
			doesNameMatch := dbSensor.SensorName.String == string(selSensor.sensorName)
			doesDeviceMatch := dbSensor.SensorDevice.String == string(selSensor.deviceName)
			if doesNameMatch && doesDeviceMatch {
				_, ok := sensorsToAdd[id]
				fmt.Println("checking ", id, ok)
				if !ok {
					delete(sensorsToRemove, id)
				}
			}
		}
	}

	fmt.Println("SENSORS TO ADD1")
	fmt.Println(sensorsToAdd)

	fmt.Println("SENSORS TO REMOVE")
	fmt.Println(sensorsToRemove)

	// Add the new sensors to the db
	for _, addSensor := range sensorsToAdd {
		_, err := g.tbdb.Queries.CreateSensor(ctx, sqlc.CreateSensorParams{
			SensorName:   sql.NullString{String: string(addSensor.sensorName), Valid: true},
			SensorType:   sql.NullString{String: string(addSensor.sensorType), Valid: true},
			SensorDevice: sql.NullString{String: string(addSensor.deviceName), Valid: true},
			SensorSource: sql.NullString{String: "lm_sensors", Valid: true},
		})
		if err != nil {

		}
	}

	// Remove the sensors that have been unchecked
	for _, removeSensorId := range sensorsToRemove {
		deleteId, _ := strconv.ParseInt(removeSensorId, 10, 64)
		err := g.tbdb.Queries.DeleteSensor(ctx, deleteId)

		if err != nil {

		}
	}
}
