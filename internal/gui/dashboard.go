package gui

import (
	"context"
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/mem"
	"strings"
	"tidbits/internal/sensors"
	"tidbits/internal/utils"
	"time"
)

func (g *GUI) dashboardView() *fyne.Container {

	sensorsStr := binding.NewString()
	sysStr := binding.NewString()
	go func() {
		for {
			sensorsStr.Set(g.dashSensorsBox())
			sysStr.Set(g.dashSystemInfo())
			time.Sleep(3 * time.Second)
		}
	}()

	sensorBox := widget.NewLabelWithData(sensorsStr)
	sysBox := widget.NewLabelWithData(sysStr)
	content := container.New(layout.NewGridLayout(2), sensorBox, sysBox)
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

func (g *GUI) dashSystemInfo() string {
	v, _ := mem.VirtualMemory()
	var lines []string
	totalMemMiB := int(utils.BytesToMegabytesBinary(int64(v.Total)))
	lines = append(lines, "Total: "+utils.AddCommas(totalMemMiB)+" MiB")

	freeMemMiB := int(utils.BytesToMegabytesBinary(int64(v.Free)))
	lines = append(lines, "Free: "+utils.AddCommas(freeMemMiB)+" MiB")

	cacheMemMiB := int(utils.BytesToMegabytesBinary(int64(v.Cached)))
	lines = append(lines, "Cache: "+utils.AddCommas(cacheMemMiB)+" MiB")

	lines = append(lines, fmt.Sprintf("Used: %d%%", int(v.UsedPercent)))
	return strings.Join(lines, "\n")
}
