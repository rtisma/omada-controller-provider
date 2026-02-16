resource "omada_dhcp_reservation" "office_printer" {
  name        = "Office Printer HP LaserJet"
  mac_address = "AA:BB:CC:DD:EE:FF"
  ip_address  = "192.168.1.50"
  network_id  = omada_network.main.id
  comment     = "Main office printer on 3rd floor"
}

resource "omada_dhcp_reservation" "conference_tv" {
  name        = "Conference Room TV"
  mac_address = "11:22:33:44:55:66"
  ip_address  = "192.168.1.51"
  network_id  = omada_network.main.id
}

resource "omada_dhcp_reservation" "security_camera_1" {
  name        = "Security Camera - Entrance"
  mac_address = "AA:11:BB:22:CC:33"
  ip_address  = "192.168.50.10"
  network_id  = omada_network.iot_vlan.id
  comment     = "Front entrance camera"
}

resource "omada_dhcp_reservation" "nas_server" {
  name        = "NAS Storage Server"
  mac_address = "DD:EE:FF:00:11:22"
  ip_address  = "192.168.10.100"
  network_id  = omada_network.management.id
  comment     = "Main file storage for IT department"
}
