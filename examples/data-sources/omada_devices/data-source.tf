data "omada_devices" "all" {
  site_id = "Default"
}

output "device_count" {
  value = length(data.omada_devices.all.devices)
}

output "access_points" {
  value = [
    for device in data.omada_devices.all.devices :
    device if device.type == "ap"
  ]
}

output "switches" {
  value = [
    for device in data.omada_devices.all.devices :
    device if device.type == "switch"
  ]
}

output "devices_needing_adoption" {
  value = [
    for device in data.omada_devices.all.devices :
    device.mac if device.need_adopt
  ]
}
