resource "omada_network" "guest_vlan" {
  name         = "Guest Network"
  vlan_id      = 100
  gateway      = "192.168.100.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.100.100"
  dhcp_end     = "192.168.100.200"
  lease_time   = 1440 # 24 hours in minutes
  dns_primary  = "8.8.8.8"
  dns_secondary = "8.8.4.4"
  purpose      = "guest"
}

resource "omada_network" "iot_vlan" {
  name         = "IoT Devices"
  vlan_id      = 50
  gateway      = "192.168.50.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.50.10"
  dhcp_end     = "192.168.50.254"
  purpose      = "iot"
}

resource "omada_network" "management" {
  name         = "Management Network"
  vlan_id      = 10
  gateway      = "192.168.10.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = false
}
