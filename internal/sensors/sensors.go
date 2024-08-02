package sensors

import (
	"os/exec"
)

type Reading struct {
	fulltext string
}

func GetSensorReadings() (Reading, error) {
	cmd := exec.Command("sensors")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return Reading{}, err
	}

	return Reading{fulltext: string(output)}, nil
}
