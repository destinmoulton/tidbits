package sensors

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func lmParseSensorReadings() Readings {

	cmd := exec.Command("sensors")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	var readings Readings
	reVoltage := regexp.MustCompile(`.*[0-9].*V.*`)
	reRPM := regexp.MustCompile(`.*[0-9].*RPM.*`)
	reTemp := regexp.MustCompile(`.*[0-9].*°C.*`)
	reNumber := regexp.MustCompile(`\d+(\.\d+)?`)
	lines := strings.Split(string(output), "\n")
	prevline := ""
	var device DeviceName
	for _, line := range lines {
		isSensor := false
		sensorType := ""
		if strings.HasPrefix(line, "Adapter:") {
			// previous line is the device identifier
			device = DeviceName(prevline)
			_, ok := readings[device]
			if !ok {
				readings[device] = nil
			}
		} else {
			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				data := strings.TrimSpace(parts[1])
				if reVoltage.MatchString(data) {
					isSensor = true
					sensorType = "voltage"
				} else if reRPM.MatchString(data) {
					isSensor = true
					sensorType = "fan"
				} else if reTemp.MatchString(data) {
					isSensor = true
					sensorType = "temperature"
				}

				if isSensor {
					// Get the number from the string
					numberMatch := reNumber.FindString(data)
					tmp, err := strconv.ParseFloat(numberMatch, 64)
					if err == nil {
						sensorName := SensorName(parts[0])
						sensors := readings[device]
						sensors[sensorName] = Sensor{
							Name:    sensorName,
							Type:    sensorType,
							Device:  device,
							Reading: tmp,
						}
						readings[device] = sensors
					}
				}
			}
		}

		// The previous line might have the hardware id
		prevline = line
	}
	return readings
}
