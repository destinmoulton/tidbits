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

type selectedSensor struct {
	sensorName sensors.SensorName
	sensorType sensors.SensorType
	deviceName sensors.DeviceName
}

type selectedSensors map[string]selectedSensor

func (g *GUI) sensorsView() *fyne.Container {

	placeholder := widget.NewLabel("sensors")
	//content := container.NewVScroll(placeholder)
	wrapper := container.NewVScroll(placeholder)
	chooseLmSensorsBtn := widget.NewButton("Choose lm_sensors", func() {
		g.switchTab(GTabLMSensorsSelectForm)
	})
	editSensorsBtn := widget.NewButton("Edit Sensor Display", func() {
		g.switchTab(GTabLMSensorsUserForm)
	})

	switch g.subtab {
	case GTabDefault:
		wrapper.Content = widget.NewLabel("SENSORS SECTION")
	case GTabLMSensorsSelectForm:
		form := g.buildSelectLmSensorForm()
		formLbl := widget.NewLabel("Select the sensors you want on the dashboard.")
		wrapper.Content = container.New(layout.NewVBoxLayout(), formLbl, form)
		wrapper.Content.Resize(fyne.NewSize(g.calcContentWidth()-GUIWidthFudge*2, g.calcContentHeight()))
	case GTabLMSensorsUserForm:
		form := g.buildEditSensorsForm()
		formLbl := widget.NewLabel("Edit the sensors")
		wrapper.Content = container.New(layout.NewVBoxLayout(), formLbl, form)
		wrapper.Content.Resize(fyne.NewSize(g.calcContentWidth()-GUIWidthFudge*2, g.calcContentHeight()))
	}

	menubar := container.New(layout.NewHBoxLayout(), editSensorsBtn, chooseLmSensorsBtn)

	return container.NewBorder(menubar, nil, nil, nil, wrapper)
}

func (g *GUI) buildSelectLmSensorForm() *widget.Form {
	ctx := context.Background()
	devices, readings := sensors.LMSensorsParseReadings()

	storedSensors, err := g.tbdb.Queries.GetSensorsBySource(ctx, sql.NullString{String: "lm_sensors", Valid: true})
	if err != nil {
		g.msg("ERROR: failed to get the stored sensors")
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
					g.msg("checked event " + string(id))
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

	for _, dbSensor := range storedSettings {
		for id, selSensor := range chkdSensors {
			doesNameMatch := dbSensor.SensorName.String == string(selSensor.sensorName)
			doesDeviceMatch := dbSensor.SensorDevice.String == string(selSensor.deviceName)
			if doesNameMatch && doesDeviceMatch {
				_, ok := sensorsToAdd[id]
				if !ok {
					delete(sensorsToRemove, id)
				}
			}
		}
	}

	// Add the new sensors to the db
	order := 0
	for _, addSensor := range sensorsToAdd {
		name := sql.NullString{String: string(addSensor.sensorName), Valid: true}
		_, err := g.tbdb.Queries.CreateSensor(ctx, sqlc.CreateSensorParams{
			SensorName:   name,
			SensorType:   sql.NullString{String: string(addSensor.sensorType), Valid: true},
			SensorDevice: sql.NullString{String: string(addSensor.deviceName), Valid: true},
			SensorSource: sql.NullString{String: "lm_sensors", Valid: true},
			UserLabel:    name,
			SensorOrder:  sql.NullInt64{Int64: int64(order), Valid: true},
		})
		if err != nil {
			g.msg("error adding sensorDataItem to db", err)
		}
		order += 1
	}

	// Remove the sensors that have been unchecked
	for _, removeSensorId := range sensorsToRemove {
		deleteId, _ := strconv.ParseInt(removeSensorId, 10, 64)
		err := g.tbdb.Queries.DeleteSensor(ctx, deleteId)

		if err != nil {
			g.msg("error removing sensorDataItem from db", err)
		}
	}
}

func (g *GUI) buildEditSensorsForm() *widget.Form {
	ctx := context.Background()

	storedSensors, err := g.tbdb.Queries.GetSensorsBySource(ctx, sql.NullString{String: "lm_sensors", Valid: true})
	if err != nil {
		g.msg("ERROR: failed to get the stored sensors")
		return nil
	}
	var rows []*widget.FormItem
	labels := make(map[string]string)
	orders := make(map[string]string)
	shouldLog := make(map[string]bool)
	for _, sens := range storedSensors {
		id := strconv.FormatInt(sens.ID, 10)
		fmt.Println("id", id)
		labels[id] = sens.UserLabel.String
		orders[id] = strconv.FormatInt(sens.SensorOrder.Int64, 10)
		shouldLog[id] = sens.ShouldLog.Int64 == 1

		// user label for the sensorDataItem
		entryUserLabel := widget.NewEntry()
		entryUserLabel.SetText(labels[id])
		entryUserLabel.OnChanged = func(value string) {
			labels[id] = value
		}

		labelFormItem := &widget.FormItem{
			Text:   sens.SensorName.String + " Label",
			Widget: entryUserLabel,
		}

		rows = append(rows, labelFormItem)

		// order for the sensorDataItem
		orderLabel := widget.NewEntry()
		orderLabel.SetText(orders[id])
		orderLabel.OnChanged = func(value string) {
			orders[id] = value
		}

		orderFormItem := &widget.FormItem{
			Text:   sens.SensorName.String + " Order",
			Widget: orderLabel,
		}

		rows = append(rows, orderFormItem)

		checkbox := widget.NewCheck("", func(value bool) {
			shouldLog[id] = value
		})
		checkItem := &widget.FormItem{
			Text:   sens.SensorName.String + " Log this sensorDataItem?",
			Widget: checkbox,
		}

		rows = append(rows, checkItem)

		// blank row
		rows = append(rows, &widget.FormItem{
			Text:   "",
			Widget: widget.NewLabel(""),
		})
	}
	return &widget.Form{
		Items: rows,
		OnSubmit: func() { // optional, handle form submission
			g.processUserSensor(labels, orders, shouldLog)
		},
	}
}

func (g *GUI) processUserSensor(labels map[string]string, orders map[string]string, shouldLog map[string]bool) {

	ctx := context.Background()
	for id, label := range labels {

		updateId, _ := strconv.ParseInt(id, 10, 64)

		var storeShouldLog int64 = 0
		if shouldLog[id] {
			storeShouldLog = 1
		}
		storeOrder, err := strconv.ParseInt(orders[id], 10, 64)
		if err != nil {
			g.msg("error parsing order to int64", err)
		}
		uerr := g.tbdb.Queries.UpdateSensor(ctx, sqlc.UpdateSensorParams{
			ID:          updateId,
			UserLabel:   sql.NullString{String: label, Valid: true},
			ShouldLog:   sql.NullInt64{Int64: storeShouldLog, Valid: true},
			SensorOrder: sql.NullInt64{Int64: storeOrder, Valid: true},
		})
		if uerr != nil {
			g.msg("error updating sensorDataItem to db", err)
		}
	}
}
