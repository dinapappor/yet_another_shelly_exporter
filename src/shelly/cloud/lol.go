package cloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Device struct {
	Id          string
	Name        string
	Image       string
	Category    string
	Type        string
	Room_Id     int
	Channel     int
	Gen         int
	Ip          string
	Relay_usage string
}

type DeviceStatusRelay struct {
	Ison            bool
	Has_timer       bool
	Timer_started   int
	Timer_duration  int
	Timer_remaining int
	Source          string
}

type DeviceStatusMeter struct {
	Power    int
	Is_valid bool
}

type DeviceStatus struct {
	Temperature     float32
	Overtemperature bool
	Uptime          int
	Relays          []DeviceStatusRelay `json:"relays"`
	Meters          []DeviceStatusMeter `json:"meters"`
}

type DeviceStatusResultData struct {
	Devices map[string]DeviceStatus `json:"devices_status"`
}

type DeviceStatusResult struct {
	IsOK bool                   `json:"isok"`
	Data DeviceStatusResultData `json:"data"`
}

type Room struct {
	Id             int
	Image          string
	Modified       int
	Name           string
	Overview_style bool
	Position       int
}

type ResultData struct {
	Devices map[string]Device `json:"devices"`
	Rooms   map[string]Room   `json:"rooms"`
}

type Result struct {
	IsOK bool       `json:"isok"`
	Data ResultData `json:"data"`
}

func GetDevicesAndRooms() Result {
	host := os.Getenv("SHELLY_HOST")
	auth_key := os.Getenv("SHELLY_AUTH_KEY")

	DeviceListUrl := fmt.Sprintf("https://%s/interface/device/get_all_lists?show_info=true&no_shared=true&auth_key=%s", host, auth_key)
	resp, err := http.Get(DeviceListUrl)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result Result

	json.Unmarshal([]byte(body), &result)
	return result
}

func GetDeviceMetrics() DeviceStatusResult {
	host := os.Getenv("SHELLY_HOST")
	auth_key := os.Getenv("SHELLY_AUTH_KEY")

	DeviceStatusUrl := fmt.Sprintf("https://%s/device/all_status?show_info=true&no_shared=true&auth_key=%s", host, auth_key)

	devicestatusresp, err := http.Get(DeviceStatusUrl)
	if err != nil {
		log.Fatalln(err)
	}

	devicestatusbody, err := ioutil.ReadAll(devicestatusresp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var devicestatusresult DeviceStatusResult
	json.Unmarshal([]byte(devicestatusbody), &devicestatusresult)

	return devicestatusresult
}
