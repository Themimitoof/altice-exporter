package collectors

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/themimitoof/altice-exporter/pkg/devices"
)

type AlticeCollector struct {
	LinkStatusMetric *prometheus.Desc
	RXRSSIMetric     *prometheus.Desc
	TXRSSIMetric     *prometheus.Desc
	Ctx              devices.Device
}

type DeviceConfiguration struct {
	Device devices.Device
	Model  string
}

func NewAlticeCollector(device DeviceConfiguration) *AlticeCollector {
	labels := prometheus.Labels{
		"host":  device.Device.GetConnectionInfo().Hostname,
		"model": device.Model,
	}

	return &AlticeCollector{
		LinkStatusMetric: prometheus.NewDesc(
			"link_status",
			"Shows if the link status is active or not",
			nil, labels,
		),
		RXRSSIMetric: prometheus.NewDesc(
			"rxrssi_metric",
			"Shows the RX RSSI value from the optical tranceiver",
			nil, labels,
		),
		TXRSSIMetric: prometheus.NewDesc(
			"txrssi_metric",
			"Shows the TX RSSI value from the optical tranceiver",
			nil, labels,
		),
		Ctx: device.Device,
	}
}

func (collector *AlticeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.LinkStatusMetric
	ch <- collector.RXRSSIMetric
	ch <- collector.TXRSSIMetric
}

func (collector *AlticeCollector) Collect(ch chan<- prometheus.Metric) {
	transceiverInfo, err := collector.Ctx.GetTransceiverInfo()

	if err != nil {
		log.Fatal("Unable to collect Transceiver metrics.")
		return
	}

	var linkStatus = 0
	if transceiverInfo.LinkStatus {
		linkStatus = 1
	}

	ch <- prometheus.MustNewConstMetric(collector.LinkStatusMetric, prometheus.GaugeValue, float64(linkStatus))
	ch <- prometheus.MustNewConstMetric(collector.RXRSSIMetric, prometheus.GaugeValue, transceiverInfo.RXRSSI)
	ch <- prometheus.MustNewConstMetric(collector.TXRSSIMetric, prometheus.GaugeValue, transceiverInfo.TXRSSI)
}
