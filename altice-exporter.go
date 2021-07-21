package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/themimitoof/altice-exporter/pkg/collectors"
	"github.com/themimitoof/altice-exporter/pkg/devices"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	deviceType     = kingpin.Flag("device.type", "Device type to collect").Required().Enum("GR241AG", "GS0100GH")
	deviceHostname = kingpin.Flag("device.hostname", "IP Address:port of the device to collect").Required().String()
	deviceUsername = kingpin.Flag("device.username", "Username of the device to collect").Required().String()
	devicePassword = kingpin.Flag("device.password", "Password of the device to collect").Required().String()

	// Prometheus client settings
	listenAddress = kingpin.Flag("web.listen-address", "Update the bind address:port for the exporter").Default(":9876").String()
	metricsRoute  = kingpin.Flag("web.route-path", "Update the route where the metrics will be exposed").Default("/metrics").String()
)

func main() {
	kingpin.Parse()

	var connInfo = devices.ConnectionInfo{
		Hostname: *deviceHostname,
		Username: *deviceUsername,
		Password: *devicePassword,
	}

	if *deviceType == "GR241AG" {
		deviceInterface := devices.GR241AG{ConnectionInfo: connInfo}
		collector := collectors.NewAlticeCollector(collectors.DeviceConfiguration{
			Device: deviceInterface.GetDevice(),
			Model:  *deviceType,
		})

		prometheus.MustRegister(collector)
	}

	if *deviceType == "GS0100GH" {
		deviceInterface := devices.GS0100GH{ConnectionInfo: connInfo}
		collector := collectors.NewAlticeCollector(collectors.DeviceConfiguration{
			Device: deviceInterface.GetDevice(),
			Model:  *deviceType,
		})

		prometheus.MustRegister(collector)
	}

	http.Handle(*metricsRoute, promhttp.Handler())
	log.Infof("Beginning to serve on port %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
