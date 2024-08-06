package sensors

type DeviceName string
type SensorName string
type Sensor struct {
	Name    SensorName
	Type    string
	Device  DeviceName
	Reading float64
}

type Sensors map[SensorName]Sensor

type Readings map[DeviceName]Sensors

func GetSensorReadings() Readings {
	return lmParseSensorReadings()
}
