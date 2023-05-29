# yet_another_shelly_exporter
Because it wasn't invented here


This exporter is created as a learning experience coding in go. It's quite rudimentary and only exposes two metrics:

```
# HELP shelly_meter_power Shows the current power draw
# TYPE shelly_meter_power gauge
shelly_meter_power{device_name="shelly-device-name",room="room_name_comes_here "} 0
# HELP shelly_relay_status Shows whether or not relay is on
# TYPE shelly_relay_status gauge
shelly_relay_status{device_name="shelly-device-name",room="room_name_comes_here "} 0
```
