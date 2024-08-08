package gui

import (
	"context"
	"database/sql"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"
	"tidbits/internal/sensors"
	"time"
)

func (g *GUI) dashboardView() *fyne.Container {

	sensorsStr := binding.NewString()
	go func() {
		for {
			sensorsStr.Set(g.dashSensorsBox())
			time.Sleep(3 * time.Second)
		}
	}()

	sensorBox := widget.NewLabelWithData(sensorsStr)
	content := container.New(layout.NewGridLayout(2), sensorBox)
	return content
}

func (g *GUI) dashSensorsBox() string {
	ctx := context.Background()
	_, readings := sensors.LMSensorsParseReadings()
	storedSensors, err := g.tbdb.Queries.GetSensorsBySource(ctx, sql.NullString{String: "lm_sensors", Valid: true})
	if err != nil {
		return ""
	}

	var lines []string
	for _, dbSensor := range storedSensors {
		deviceName := sensors.DeviceName(dbSensor.SensorDevice.String)
		for _, reading := range readings[deviceName] {
			if reading.Name == sensors.SensorName(dbSensor.SensorName.String) {
				line := dbSensor.UserLabel.String + ": " + reading.Format()
				lines = append(lines, line)
			}
		}
	}
	return strings.Join(lines, "\n")
}
