package devices

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/themimitoof/altice-exporter/pkg/ssh_"
)

type GR241AG struct {
	ConnectionInfo ConnectionInfo
}

func (device *GR241AG) GetDevice() Device {
	return device
}

func (device *GR241AG) GetConnectionInfo() ConnectionInfo {
	return device.ConnectionInfo
}

func (device *GR241AG) GetTransceiverInfo() (*DeviceOpticTransceiver, error) {
	client, err := ssh_.ConnectToHost(
		device.ConnectionInfo.Hostname,
		device.ConnectionInfo.Username,
		device.ConnectionInfo.Password,
		nil,
	)

	if err != nil {
		log.Errorf("Unable to connect to %v via SSH. Unable to get new metrics.\n", device.ConnectionInfo.Hostname)
		return nil, err
	}

	defer client.Close()

	result, err := ssh_.RunCommand(client, "statistics/optical/show --option=all", 1)

	if err != nil {
		log.Errorf("Unable to run the command to retrieve optical tranceiver info on %v via SSH. Unable to get new metrics.\n", device.ConnectionInfo.Hostname)
		return nil, err
	}

	var (
		output_lines []string = strings.Split(result, "\n")
		linkStatus   bool     = false
		RXRSSI       float64
		TXRSSI       float64
	)

	for i := 0; i < len(output_lines); i++ {
		if strings.Contains(output_lines[i], "Link status") {
			if strings.Contains(output_lines[i], "Up") {
				linkStatus = true
			}
		}
		if strings.Contains(output_lines[i], "Received Optical Level") {
			re := regexp.MustCompile(`-\d+.\d+`)
			RXRSSI, err = strconv.ParseFloat(re.FindString(output_lines[i]), 64)

			if err != nil {
				panic(err)
			}
		}

		if strings.Contains(output_lines[i], "Transmited Optical Level") {
			re := regexp.MustCompile(`\d+.\d+`)
			TXRSSI, err = strconv.ParseFloat(re.FindString(output_lines[i]), 64)

			if err != nil {
				panic(err)
			}
		}
	}

	return &DeviceOpticTransceiver{
		LinkStatus: linkStatus,
		RXRSSI:     RXRSSI,
		TXRSSI:     TXRSSI,
	}, nil
}
