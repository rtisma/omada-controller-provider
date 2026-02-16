# Terraform Provider for Omada Controller - Resources & Data Sources

Complete reference documentation for all resources and data sources in the Omada Terraform provider.

## Table of Contents

- [Provider Configuration](#provider-configuration)
- [Data Sources](#data-sources)
  - [omada_site](#omada_site)
  - [omada_devices](#omada_devices)
- [Resources](#resources)
  - [omada_network](#omada_network)
  - [omada_ssid](#omada_ssid)
  - [omada_dhcp_reservation](#omada_dhcp_reservation)
  - [omada_device](#omada_device)

---

## Provider Configuration

The Omada provider is used to interact with TP-Link Omada Controller for managing network infrastructure.

### Example Usage

```hcl
provider "omada" {
  host     = "https://192.168.1.1:8043"
  username = "admin"
  password = "your-password"
  site_id  = "Default"
  insecure = true
}
```

### Argument Reference

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `host` | string | Yes | - | The URL of your Omada Controller (e.g., `https://192.168.1.1:8043`) |
| `username` | string | Yes | - | Admin username for authentication |
| `password` | string | Yes | - | Admin password for authentication |
| `site_id` | string | No | `"Default"` | The site name to manage |
| `insecure` | bool | No | `false` | Skip TLS certificate verification (useful for self-signed certs) |

### Environment Variables

While not currently implemented, you can prepare for future support by setting:
- `OMADA_HOST`
- `OMADA_USERNAME`
- `OMADA_PASSWORD`
- `OMADA_SITE_ID`
- `OMADA_INSECURE`

---

## Data Sources

### omada_site

Retrieves information about an Omada site.

#### Example Usage

```hcl
# Get the default site
data "omada_site" "default" {
  name = "Default"
}

# Use provider's configured site_id
data "omada_site" "current" {}

output "site_details" {
  value = {
    id       = data.omada_site.default.id
    name     = data.omada_site.default.name
    timezone = data.omada_site.default.timezone
  }
}
```

#### Argument Reference

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | No | Site name to query. If not specified, uses the provider's `site_id` |

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | Unique identifier of the site |
| `name` | string | Name of the site |
| `type` | string | Type of the site |
| `location` | string | Physical location of the site |
| `timezone` | string | Timezone configuration |
| `scenario` | string | Site scenario/template |

---

### omada_devices

Retrieves a list of all devices (Access Points, Switches, Gateways) in a site.

#### Example Usage

```hcl
# Get all devices
data "omada_devices" "all" {}

# Get devices from a specific site
data "omada_devices" "other_site" {
  site_id = "OfficeBranch"
}

# Filter access points
locals {
  access_points = [
    for device in data.omada_devices.all.devices :
    device if device.type == "ap"
  ]
}

# Filter devices needing adoption
locals {
  pending_devices = [
    for device in data.omada_devices.all.devices :
    device if device.need_adopt
  ]
}

output "device_summary" {
  value = {
    total    = length(data.omada_devices.all.devices)
    aps      = length([for d in data.omada_devices.all.devices : d if d.type == "ap"])
    switches = length([for d in data.omada_devices.all.devices : d if d.type == "switch"])
    gateways = length([for d in data.omada_devices.all.devices : d if d.type == "gateway"])
  }
}
```

#### Argument Reference

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `site_id` | string | No | Site to query devices from. If not specified, uses the provider's `site_id` |

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `devices` | list | List of device objects (see below) |

#### Device Object Attributes

| Attribute | Type | Description |
|-----------|------|-------------|
| `mac` | string | MAC address of the device |
| `name` | string | Name of the device |
| `type` | string | Device type (`ap`, `switch`, `gateway`) |
| `model` | string | Device model |
| `status` | string | Connection status |
| `led_enabled` | bool | Whether LED is enabled |
| `location` | string | Physical location description |
| `site` | string | Site name |
| `ip` | string | IP address of the device |
| `uptime` | int | Uptime in seconds |
| `firmware_version` | string | Firmware version |
| `need_adopt` | bool | Whether device needs adoption |

---

## Resources

### omada_network

Manages a network/VLAN in Omada Controller.

#### Example Usage

```hcl
# Basic network with DHCP
resource "omada_network" "main" {
  name         = "Main Network"
  vlan_id      = 1
  gateway      = "192.168.1.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.1.100"
  dhcp_end     = "192.168.1.200"
  lease_time   = 1440  # 24 hours
  dns_primary  = "8.8.8.8"
  dns_secondary = "8.8.4.4"
  purpose      = "general"
}

# Guest network
resource "omada_network" "guest" {
  name         = "Guest Network"
  vlan_id      = 100
  gateway      = "192.168.100.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.100.10"
  dhcp_end     = "192.168.100.250"
  purpose      = "guest"
}

# Static IP network (no DHCP)
resource "omada_network" "servers" {
  name         = "Server Network"
  vlan_id      = 10
  gateway      = "192.168.10.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = false
}

# IoT network with custom DNS
resource "omada_network" "iot" {
  name          = "IoT Devices"
  vlan_id       = 50
  gateway       = "192.168.50.1"
  netmask       = "255.255.255.0"
  dhcp_enabled  = true
  dhcp_start    = "192.168.50.10"
  dhcp_end      = "192.168.50.254"
  dns_primary   = "192.168.1.1"  # Local DNS
  domain_name   = "iot.local"
  purpose       = "iot"
}
```

#### Argument Reference

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `name` | string | Yes | - | Name of the network |
| `vlan_id` | int | Yes | - | VLAN ID (1-4094) |
| `gateway` | string | Yes | - | Gateway IP address |
| `netmask` | string | Yes | - | Subnet mask (e.g., `255.255.255.0`) |
| `site_id` | string | No | Provider's site_id | Site where network belongs |
| `dhcp_enabled` | bool | No | `false` | Enable DHCP server |
| `dhcp_start` | string | No | - | Start of DHCP range |
| `dhcp_end` | string | No | - | End of DHCP range |
| `lease_time` | int | No | `1440` | DHCP lease time in minutes |
| `dns_primary` | string | No | - | Primary DNS server |
| `dns_secondary` | string | No | - | Secondary DNS server |
| `domain_name` | string | No | - | Domain name for DHCP |
| `purpose` | string | No | - | Network purpose (`general`, `guest`, `iot`, etc.) |

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | Unique identifier of the network |

#### Import

Networks can be imported using their ID:

```bash
terraform import omada_network.main <network-id>
```

---

### omada_ssid

Manages a wireless network (SSID) in Omada Controller.

#### Example Usage

```hcl
# Corporate WPA3 network
resource "omada_ssid" "corporate" {
  name            = "Corporate WiFi"
  ssid            = "CorpNet"
  security_mode   = "wpa3-personal"
  password        = var.wifi_password
  vlan_id         = 10
  enabled         = true
  hide_ssid       = false
  guest_network   = false
  client_isolation = false
  band_2_4g_enabled = true
  band_5g_enabled   = true
  band_6g_enabled   = false
}

# Guest network with rate limiting
resource "omada_ssid" "guest" {
  name             = "Guest WiFi"
  ssid             = "Guest-Network"
  security_mode    = "wpa2-personal"
  password         = var.guest_password
  vlan_id          = 100
  enabled          = true
  guest_network    = true
  client_isolation = true
  rate_limit       = true
  downlink_limit   = 10000  # 10 Mbps
  uplink_limit     = 5000   # 5 Mbps
  max_clients      = 50
  band_2_4g_enabled = true
  band_5g_enabled   = true
}

# Hidden IoT network (2.4GHz only)
resource "omada_ssid" "iot" {
  name              = "IoT Devices"
  ssid              = "IoT-Network"
  security_mode     = "wpa2-personal"
  password          = var.iot_password
  vlan_id           = 50
  enabled           = true
  hide_ssid         = true
  client_isolation  = true
  band_2_4g_enabled = true
  band_5g_enabled   = false
}

# Enterprise network with RADIUS
resource "omada_ssid" "enterprise" {
  name           = "Enterprise Network"
  ssid           = "Corp-Secure"
  security_mode  = "wpa2-enterprise"
  radius_profile = "MainRADIUS"
  vlan_id        = 20
  enabled        = true
  band_2_4g_enabled = true
  band_5g_enabled   = true
}

# Open network (not recommended)
resource "omada_ssid" "open" {
  name          = "Public WiFi"
  ssid          = "Free-WiFi"
  security_mode = "open"
  vlan_id       = 200
  enabled       = true
  portal_enabled = true  # Use with captive portal
}
```

#### Argument Reference

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `name` | string | Yes | - | Descriptive name of the SSID |
| `ssid` | string | Yes | - | Network name broadcast to clients |
| `security_mode` | string | Yes | - | Security mode: `wpa2-personal`, `wpa3-personal`, `wpa2-enterprise`, `wpa3-enterprise`, `open` |
| `site_id` | string | No | Provider's site_id | Site where SSID belongs |
| `password` | string | No | - | WiFi password (required for WPA2/WPA3 personal modes) |
| `enabled` | bool | No | `true` | Enable/disable the SSID |
| `hide_ssid` | bool | No | `false` | Hide SSID from broadcast |
| `vlan_id` | int | No | - | VLAN ID for this SSID |
| `guest_network` | bool | No | `false` | Mark as guest network |
| `client_isolation` | bool | No | `false` | Isolate clients from each other |
| `band_2_4g_enabled` | bool | No | `true` | Enable on 2.4GHz band |
| `band_5g_enabled` | bool | No | `true` | Enable on 5GHz band |
| `band_6g_enabled` | bool | No | `false` | Enable on 6GHz band |
| `max_clients` | int | No | - | Maximum number of clients |
| `rate_limit` | bool | No | `false` | Enable bandwidth limiting |
| `downlink_limit` | int | No | - | Download limit in Kbps (requires `rate_limit=true`) |
| `uplink_limit` | int | No | - | Upload limit in Kbps (requires `rate_limit=true`) |
| `schedule_enabled` | bool | No | `false` | Enable scheduled activation |
| `portal_enabled` | bool | No | `false` | Enable captive portal |
| `radius_profile` | string | No | - | RADIUS profile name (for enterprise modes) |

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | Unique identifier of the SSID |

#### Import

SSIDs can be imported using their ID:

```bash
terraform import omada_ssid.corporate <ssid-id>
```

---

### omada_dhcp_reservation

Manages a DHCP reservation (static IP assignment by MAC address).

#### Example Usage

```hcl
# Printer reservation
resource "omada_dhcp_reservation" "printer" {
  name        = "Office Printer HP LaserJet"
  mac_address = "AA:BB:CC:DD:EE:FF"
  ip_address  = "192.168.1.50"
  network_id  = omada_network.main.id
  comment     = "Main office printer on 3rd floor"
}

# Server reservation
resource "omada_dhcp_reservation" "nas" {
  name        = "NAS Storage Server"
  mac_address = "11:22:33:44:55:66"
  ip_address  = "192.168.10.100"
  network_id  = omada_network.servers.id
}

# Security camera
resource "omada_dhcp_reservation" "camera_1" {
  name        = "Security Camera - Entrance"
  mac_address = "AA:11:BB:22:CC:33"
  ip_address  = "192.168.50.10"
  network_id  = omada_network.iot.id
  comment     = "Front entrance camera"
}

# Dynamic reservations from data
locals {
  devices_to_reserve = {
    "printer1" = { mac = "AA:BB:CC:DD:EE:11", ip = "192.168.1.51" }
    "printer2" = { mac = "AA:BB:CC:DD:EE:22", ip = "192.168.1.52" }
  }
}

resource "omada_dhcp_reservation" "printers" {
  for_each = local.devices_to_reserve

  name        = each.key
  mac_address = each.value.mac
  ip_address  = each.value.ip
  network_id  = omada_network.main.id
}
```

#### Argument Reference

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `name` | string | Yes | - | Descriptive name for the device |
| `mac_address` | string | Yes | - | MAC address (format: `AA:BB:CC:DD:EE:FF`) |
| `ip_address` | string | Yes | - | IP address to reserve |
| `network_id` | string | Yes | - | Network ID where reservation applies |
| `site_id` | string | No | Provider's site_id | Site where reservation belongs |
| `comment` | string | No | - | Additional notes |

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | Unique identifier of the reservation |

#### Import

DHCP reservations can be imported using their ID:

```bash
terraform import omada_dhcp_reservation.printer <reservation-id>
```

---

### omada_device

Manages an Omada device (Access Point, Switch, or Gateway) configuration.

**Note:** This resource manages the configuration of existing devices. Devices must already be connected to the controller before they can be managed by Terraform.

#### Example Usage

```hcl
# Configure a single device
resource "omada_device" "lobby_ap" {
  mac         = "AA:BB:CC:DD:EE:FF"
  name        = "Lobby Access Point"
  led_enabled = false
  location    = "Main Lobby - 1st Floor"
}

# Configure multiple devices
resource "omada_device" "conference_aps" {
  for_each = {
    "conf_a" = { mac = "11:22:33:44:55:66", location = "Conference Room A" }
    "conf_b" = { mac = "11:22:33:44:55:77", location = "Conference Room B" }
  }

  mac         = each.value.mac
  name        = "AP - ${each.key}"
  led_enabled = true
  location    = each.value.location
}

# Manage all access points dynamically
data "omada_devices" "all" {}

resource "omada_device" "managed_aps" {
  for_each = {
    for device in data.omada_devices.all.devices :
    device.mac => device
    if device.type == "ap"
  }

  mac         = each.value.mac
  name        = "AP - ${each.value.model}"
  led_enabled = false
  location    = "Auto-managed via Terraform"
}

# Configure a switch
resource "omada_device" "core_switch" {
  mac         = "AA:11:BB:22:CC:33"
  name        = "Core Network Switch"
  led_enabled = true
  location    = "Server Room - Rack 1"
}
```

#### Argument Reference

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `mac` | string | Yes | - | MAC address of the device (format: `AA:BB:CC:DD:EE:FF`) |
| `name` | string | Yes | - | Name/description of the device |
| `site_id` | string | No | Provider's site_id | Site where device belongs |
| `led_enabled` | bool | No | `true` | Enable/disable LED indicator |
| `location` | string | No | - | Physical location description |

#### Attribute Reference

| Attribute | Type | Description |
|-----------|------|-------------|
| `type` | string | Device type (ap, switch, gateway) - Read-only |
| `model` | string | Device model - Read-only |
| `status` | string | Connection status - Read-only |

#### Important Notes

1. **Device Must Exist**: The device must be connected to the controller before it can be managed by Terraform
2. **Adoption**: If the device needs adoption, Terraform will automatically adopt it during creation
3. **Deletion**: Deleting this resource will "forget" the device from the controller, removing it from management

#### Import

Devices can be imported using their MAC address:

```bash
terraform import omada_device.lobby_ap AA:BB:CC:DD:EE:FF
```

---

## Common Patterns

### Complete Network Setup

```hcl
# Create network
resource "omada_network" "main" {
  name         = "Main Network"
  vlan_id      = 1
  gateway      = "192.168.1.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.1.100"
  dhcp_end     = "192.168.1.200"
}

# Create WiFi for that network
resource "omada_ssid" "main_wifi" {
  name          = "Main WiFi"
  ssid          = "MyNetwork"
  security_mode = "wpa3-personal"
  password      = var.wifi_password
  vlan_id       = omada_network.main.vlan_id
}

# Add DHCP reservations
resource "omada_dhcp_reservation" "printer" {
  name        = "Office Printer"
  mac_address = "AA:BB:CC:DD:EE:FF"
  ip_address  = "192.168.1.50"
  network_id  = omada_network.main.id
}

# Configure devices
data "omada_devices" "all" {}

resource "omada_device" "aps" {
  for_each = {
    for device in data.omada_devices.all.devices :
    device.mac => device
    if device.type == "ap"
  }

  mac         = each.value.mac
  name        = "AP - ${each.value.model}"
  led_enabled = false
}
```

### Multi-VLAN Setup

```hcl
locals {
  networks = {
    main = { vlan = 1,   cidr = "192.168.1.0/24",   purpose = "general" }
    guest = { vlan = 100, cidr = "192.168.100.0/24", purpose = "guest" }
    iot = { vlan = 50,  cidr = "192.168.50.0/24",  purpose = "iot" }
  }
}

resource "omada_network" "networks" {
  for_each = local.networks

  name         = title(each.key)
  vlan_id      = each.value.vlan
  gateway      = cidrhost(each.value.cidr, 1)
  netmask      = cidrnetmask(each.value.cidr)
  dhcp_enabled = true
  dhcp_start   = cidrhost(each.value.cidr, 10)
  dhcp_end     = cidrhost(each.value.cidr, 250)
  purpose      = each.value.purpose
}

resource "omada_ssid" "ssids" {
  for_each = omada_network.networks

  name          = "${each.key} WiFi"
  ssid          = title(each.key)
  security_mode = "wpa2-personal"
  password      = var.wifi_passwords[each.key]
  vlan_id       = each.value.vlan_id
  guest_network = each.key == "guest"
}
```

### Conditional Device Management

```hcl
# Only manage devices in specific locations
locals {
  office_aps = [
    for device in data.omada_devices.all.devices :
    device
    if device.type == "ap" && length(regexall("Office", device.name)) > 0
  ]
}

resource "omada_device" "office_aps" {
  for_each = { for ap in local.office_aps : ap.mac => ap }

  mac         = each.value.mac
  name        = each.value.name
  led_enabled = false
  location    = "Managed by Terraform"
}
```

---

## Troubleshooting

### Authentication Issues

If you encounter authentication errors:

```hcl
provider "omada" {
  host     = "https://192.168.1.1:8043"
  username = "admin"
  password = "password"
  insecure = true  # Enable for self-signed certificates
}
```

### Import Existing Resources

To bring existing infrastructure under Terraform management:

```bash
# Find resource IDs in Omada Controller or via API
terraform import omada_network.main <network-id>
terraform import omada_ssid.wifi <ssid-id>
terraform import omada_device.ap1 AA:BB:CC:DD:EE:FF
```

### State Refresh

To detect drift (changes made outside Terraform):

```bash
terraform refresh
terraform plan  # Shows differences
```

---

## Additional Resources

- [Main README](README.md)
- [Contributing Guide](CONTRIBUTING.md)
- [Changelog](CHANGELOG.md)
- [Examples Directory](examples/)
