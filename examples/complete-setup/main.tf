terraform {
  required_providers {
    omada = {
      source = "your-org/omada"
    }
  }
}

provider "omada" {
  host     = var.omada_host
  username = var.omada_username
  password = var.omada_password
  site_id  = var.omada_site
  insecure = true
}

# Variables
variable "omada_host" {
  description = "Omada Controller URL"
  type        = string
  default     = "https://192.168.1.1:8043"
}

variable "omada_username" {
  description = "Omada Controller username"
  type        = string
  sensitive   = true
}

variable "omada_password" {
  description = "Omada Controller password"
  type        = string
  sensitive   = true
}

variable "omada_site" {
  description = "Omada site name"
  type        = string
  default     = "Default"
}

# Get site information
data "omada_site" "main" {
  name = var.omada_site
}

# Create networks/VLANs
resource "omada_network" "main" {
  name         = "Main Network"
  vlan_id      = 1
  gateway      = "192.168.1.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.1.100"
  dhcp_end     = "192.168.1.200"
  dns_primary  = "8.8.8.8"
  dns_secondary = "8.8.4.4"
  purpose      = "general"
}

resource "omada_network" "guest" {
  name         = "Guest Network"
  vlan_id      = 100
  gateway      = "192.168.100.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.100.10"
  dhcp_end     = "192.168.100.250"
  dns_primary  = "1.1.1.1"
  dns_secondary = "1.0.0.1"
  purpose      = "guest"
}

resource "omada_network" "iot" {
  name         = "IoT Devices"
  vlan_id      = 50
  gateway      = "192.168.50.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.50.10"
  dhcp_end     = "192.168.50.254"
  purpose      = "iot"
}

# Create WiFi networks (SSIDs)
resource "omada_ssid" "corporate" {
  name            = "Corporate WiFi"
  ssid            = "CompanyCorp"
  security_mode   = "wpa3-personal"
  password        = var.wifi_corporate_password
  vlan_id         = omada_network.main.vlan_id
  enabled         = true
  guest_network   = false
  client_isolation = false
  band_2_4g_enabled = true
  band_5g_enabled   = true
}

resource "omada_ssid" "guest" {
  name             = "Guest WiFi"
  ssid             = "Guest-WiFi"
  security_mode    = "wpa2-personal"
  password         = var.wifi_guest_password
  vlan_id          = omada_network.guest.vlan_id
  enabled          = true
  guest_network    = true
  client_isolation = true
  rate_limit       = true
  downlink_limit   = 10000
  uplink_limit     = 5000
  band_2_4g_enabled = true
  band_5g_enabled   = true
}

resource "omada_ssid" "iot" {
  name              = "IoT Devices"
  ssid              = "IoT-Network"
  security_mode     = "wpa2-personal"
  password          = var.wifi_iot_password
  vlan_id           = omada_network.iot.vlan_id
  enabled           = true
  hide_ssid         = true
  client_isolation  = true
  band_2_4g_enabled = true
  band_5g_enabled   = false
}

variable "wifi_corporate_password" {
  description = "Corporate WiFi password"
  type        = string
  sensitive   = true
}

variable "wifi_guest_password" {
  description = "Guest WiFi password"
  type        = string
  sensitive   = true
}

variable "wifi_iot_password" {
  description = "IoT WiFi password"
  type        = string
  sensitive   = true
}

# DHCP Reservations
resource "omada_dhcp_reservation" "printer" {
  name        = "Office Printer"
  mac_address = "AA:BB:CC:DD:EE:FF"
  ip_address  = "192.168.1.50"
  network_id  = omada_network.main.id
  comment     = "Main office printer"
}

resource "omada_dhcp_reservation" "nas" {
  name        = "NAS Server"
  mac_address = "11:22:33:44:55:66"
  ip_address  = "192.168.1.60"
  network_id  = omada_network.main.id
  comment     = "Network storage"
}

# Get all devices
data "omada_devices" "all" {
  depends_on = [
    omada_network.main,
    omada_network.guest,
    omada_network.iot,
  ]
}

# Configure devices (optional - uncomment to manage)
# resource "omada_device" "main_ap" {
#   mac         = "AA:BB:CC:DD:EE:11"
#   name        = "Main Office AP"
#   led_enabled = false
#   location    = "Main Office"
# }

# Outputs
output "site_info" {
  value = {
    id       = data.omada_site.main.id
    name     = data.omada_site.main.name
    timezone = data.omada_site.main.timezone
  }
}

output "networks" {
  value = {
    main  = omada_network.main.id
    guest = omada_network.guest.id
    iot   = omada_network.iot.id
  }
}

output "ssids" {
  value = {
    corporate = omada_ssid.corporate.id
    guest     = omada_ssid.guest.id
    iot       = omada_ssid.iot.id
  }
}

output "device_count" {
  value = length(data.omada_devices.all.devices)
}

output "device_summary" {
  value = {
    total  = length(data.omada_devices.all.devices)
    aps    = length([for d in data.omada_devices.all.devices : d if d.type == "ap"])
    switches = length([for d in data.omada_devices.all.devices : d if d.type == "switch"])
    gateways = length([for d in data.omada_devices.all.devices : d if d.type == "gateway"])
  }
}
