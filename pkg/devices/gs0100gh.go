package devices

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/themimitoof/altice-exporter/pkg/ssh_"
)

type GS0100GH struct {
	ConnectionInfo ConnectionInfo
}

func (device *GS0100GH) GetDevice() Device {
	return device
}

func (device *GS0100GH) GetConnectionInfo() ConnectionInfo {
	return device.ConnectionInfo
}

func (device *GS0100GH) GetTransceiverInfo() (*DeviceOpticTransceiver, error) {
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

	result, err := ssh_.RunCommand(client, "show gpon rssi\nshow gpon status\n", 1)

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
		if strings.Contains(output_lines[i], "LOSS_OF_SIGNAL") {
			if strings.Contains(output_lines[i], "INACTIVE") {
				linkStatus = true
			}
		}
		if strings.Contains(output_lines[i], "receive RSSI") {
			re := regexp.MustCompile(`-\d+.\d+`)
			RXRSSI, err = strconv.ParseFloat(re.FindString(output_lines[i]), 64)

			if err != nil {
				panic(err)
			}
		}

		if strings.Contains(output_lines[i], "transmit RSSI") {
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
