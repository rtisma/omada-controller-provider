resource "omada_device" "lobby_ap" {
  mac         = "AA:BB:CC:DD:EE:FF"
  name        = "Lobby Access Point"
  led_enabled = false
  location    = "Main Lobby - 1st Floor"
}

resource "omada_device" "conference_room_ap" {
  mac         = "11:22:33:44:55:66"
  name        = "Conference Room AP"
  led_enabled = true
  location    = "Conference Room A - 3rd Floor"
}

resource "omada_device" "core_switch" {
  mac         = "AA:11:BB:22:CC:33"
  name        = "Core Network Switch"
  led_enabled = true
  location    = "Server Room"
}

# Example: Dynamically manage all access points
data "omada_devices" "all" {}

resource "omada_device" "managed_aps" {
  for_each = {
    for device in data.omada_devices.all.devices :
    device.mac => device
    if device.type == "ap"
  }

  mac         = each.value.mac
  name        = "Managed AP - ${each.value.model}"
  led_enabled = false
  location    = "Auto-managed"
}
