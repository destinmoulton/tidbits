package sensors

import (
	"os/exec"
)

type sensor struct {
	id    string
	value float64
	units string
}

type sensors []sensor

type Reading struct {
	Fulltext string
}

func GetSensorReadings() (Reading, error) {
	cmd := exec.Command("sensors")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return Reading{}, err
	}

	return Reading{Fulltext: string(output)}, nil
}

//func parseSensorReading(reading string) (Reading, error) {
//	renum := regexp.MustCompile(`[0-9]`)
//	lines := strings.Split(reading, "\n")
//	prevline := ""
//	for _, line := range lines {
//		if strings.HasPrefix(line, "Adapter:") {
//			// previous line is the device identifier
//
//		} else {
//			if strings.Contains(line, ":") && renum.MatchString(line) {
//
//			}
//		}
//		// The previous line might have the hardware id
//		prevline = line
//	}
//}
