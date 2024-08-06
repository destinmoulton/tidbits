package sensors

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func LMSensorsParseReadings() (SensorDevices, SensorReadings) {

	cmd := exec.Command("sensors")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil
	}

	var readings SensorReadings = make(SensorReadings)
	var devices SensorDevices
	reVoltage := regexp.MustCompile(`.*[0-9].*V.*`)
	reMilliVoltage := regexp.MustCompile(`.*[0-9].*mV.*`)
	reRPM := regexp.MustCompile(`.*[0-9].*RPM.*`)
	reTemp := regexp.MustCompile(`.*[0-9].*°C.*`)
	reNumber := regexp.MustCompile(`\d+(\.\d+)?`)
	lines := strings.Split(string(output), "\n")
	prevline := ""
	var device DeviceName
	for _, line := range lines {
		isSensor := false
		var sensorType SensorType
		if strings.HasPrefix(line, "Adapter:") {
			// previous line is the device identifier
			device = DeviceName(prevline)
			_, ok := readings[device]
			if !ok {
				devices = append(devices, device)
				readings[device] = nil
			}
		} else {
			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				data := strings.TrimSpace(parts[1])
				if reVoltage.MatchString(data) {
					isSensor = true
					sensorType = STV
					if reMilliVoltage.MatchString(data) {
						sensorType = STmV
					}
				} else if reRPM.MatchString(data) {
					isSensor = true
					sensorType = STSpeed
				} else if reTemp.MatchString(data) {
					isSensor = true
					sensorType = STTemperature
				}

				if isSensor {
					// Get the number from the string
					numberMatch := reNumber.FindString(data)
					tmp, err := strconv.ParseFloat(numberMatch, 64)
					if err == nil {
						sensorName := SensorName(parts[0])
						sensors := readings[device]
						sensors = append(sensors, SensorReading{
							Name:    sensorName,
							Type:    sensorType,
							Device:  device,
							Reading: tmp,
							Source:  "lm_sensors",
						})
						readings[device] = sensors
					}
				}
			}
		}

		// The previous line might have the hardware id
		prevline = line
	}
	return devices, readings
}
