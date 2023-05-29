package main

import (
	"log"
	"net/http"
	"shelly_exporter/src/shelly/cloud"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type shellyCollector struct {
	// We want one of these for each different metric we want to add
	RelayStatus *prometheus.Desc
	MeterStatus *prometheus.Desc
}

func newShellyCollector() *shellyCollector {
	// we need one for each different type of
	// metric in the shellyCollector struct
	return &shellyCollector{
		RelayStatus: prometheus.NewDesc("shelly_relay_status",
			"Shows whether or not relay is on",
			[]string{"device_name", "room"},
			nil,
		),
		MeterStatus: prometheus.NewDesc("shelly_meter_power",
			"Shows the current power draw",
			[]string{"device_name", "room"},
			nil,
		),
	}
}

func (collector *shellyCollector) Describe(ch chan<- *prometheus.Desc) {
	// We need one for each different tpye of metric we want in the exporter
	ch <- collector.RelayStatus
	ch <- collector.MeterStatus
}

func (collector *shellyCollector) Collect(ch chan<- prometheus.Metric) {

	// Get the rooms names and device names/device ids
	// to allow us to populate labels later on
	devicesAndRooms := cloud.GetDevicesAndRooms()

	// Get actual metrics from shelly cloud
	deviceMetrics := cloud.GetDeviceMetrics()

	for device := range deviceMetrics.Data.Devices {
		devicedata, _ := devicesAndRooms.Data.Devices[device]
		devicestatus, _ := deviceMetrics.Data.Devices[device]
		name := devicedata.Name
		room_id := strconv.Itoa(devicedata.Room_Id)
		roomdata, _ := devicesAndRooms.Data.Rooms[room_id]
		room_name := roomdata.Name
		metricValue := 0.0
		for index, _ := range devicestatus.Relays {
			relay := devicestatus.Relays[index]
			if relay.Ison {
				metricValue = 1.0
			}
			// Relay status
			relay_metric := prometheus.MustNewConstMetric(collector.RelayStatus, prometheus.GaugeValue, metricValue, name, room_name)
			ch <- relay_metric

			// Power draw
			meter := devicestatus.Meters[index]
			meter_metric := prometheus.MustNewConstMetric(collector.MeterStatus, prometheus.GaugeValue, float64(meter.Power), name, room_name)
			ch <- meter_metric
		}
	}
}

func main() {
	metricsCollector := newShellyCollector()
	prometheus.MustRegister(metricsCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
