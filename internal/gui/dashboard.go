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

type dashSensorBinding struct {
	deviceName string
	sensorName string
	rowTitle   string
	boundValue binding.String
}
type dashSensors []dashSensorBinding

func (g *GUI) dashboardView() *fyne.Container {

	sysStr := binding.NewString()
	var boundSensors dashSensors

	g.initSensorBindings(&boundSensors)
	fmt.Println("count boundSensors " + string(len(boundSensors)))
	sensorsTable := g.buildDashSensorsTable(boundSensors)
	go func() {
		for {
			g.updateSensorData(&boundSensors)

			sysStr.Set(g.dashSystemInfo())
			time.Sleep(3 * time.Second)
		}
	}()

	sysBox := widget.NewLabelWithData(sysStr)
	content := container.New(layout.NewGridLayout(2), sensorsTable, sysBox)
	return content
}

func (g *GUI) initSensorBindings(boundSensors *dashSensors) {
	ctx := context.Background()

	_, readings := sensors.LMSensorsParseReadings()
	storedSensors, err := g.tbdb.Queries.GetSensorsBySource(ctx, sql.NullString{String: "lm_sensors", Valid: true})
	if err != nil {
		return
	}

	for _, dbSensor := range storedSensors {

		deviceName := dbSensor.SensorDevice.String
		for _, reading := range readings[sensors.DeviceName(deviceName)] {
			sensorName := string(reading.Name)
			if sensorName == dbSensor.SensorName.String {
				*boundSensors = append(*boundSensors, dashSensorBinding{
					deviceName: deviceName,
					sensorName: sensorName,
					rowTitle:   dbSensor.UserLabel.String,
					boundValue: binding.NewString(),
				})
			}
		}
	}
}

func (g *GUI) updateSensorData(boundSensors *dashSensors) {
	_, readings := sensors.LMSensorsParseReadings()

	for _, bnd := range *boundSensors {
		for _, reading := range readings[sensors.DeviceName(bnd.deviceName)] {
			if string(reading.Name) == bnd.sensorName {
				bnd.boundValue.Set(reading.Format())
			}
		}
	}
}

func (g *GUI) buildDashSensorsTable(boundSensors dashSensors) *widget.Table {
	fmt.Println(boundSensors)
	length := func() (int, int) {
		return len(boundSensors), 2
	}
	createCell := func() fyne.CanvasObject {
		return widget.NewLabel("Sensors")
	}
	updateCell := func(i widget.TableCellID, o fyne.CanvasObject) {
		row, col := i.Row, i.Col
		switch col {
		case 0:
			o.(*widget.Label).SetText(boundSensors[row].rowTitle)
		case 1:
			item := boundSensors[row].boundValue
			o.(*widget.Label).Bind(item)
		}
	}

	table := widget.NewTable(length, createCell, updateCell)

	return table
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
