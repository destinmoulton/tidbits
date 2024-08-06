package sensors

import "fmt"

type DeviceName string
type SensorType string
type SensorName string
type SensorReading struct {
	Name    SensorName
	Type    SensorType
	Device  DeviceName
	Reading float64
	Source  string
}

type Sensors []SensorReading

type SensorReadings map[DeviceName]Sensors
type SensorDevices []DeviceName

const (
	STV           SensorType = "voltage"
	STmV          SensorType = "millivoltage"
	STSpeed       SensorType = "speed"
	STTemperature SensorType = "temperature"
)

func (s *SensorReading) Format() string {
	switch s.Type {
	case STTemperature:
		return fmt.Sprintf("%.2f %s", s.Reading, "°C")
	case STV:
		return fmt.Sprintf("%.2f %s", s.Reading, "V")
	case STmV:
		return fmt.Sprintf("%.2f %s", s.Reading, "mV")
	case STSpeed:
		return fmt.Sprintf("%.2f %s", s.Reading, "RPM")
	}
	return ""
}
