---
page_title: "omada_devices Data Source - terraform-provider-omada"
subcategory: ""
description: |-
  Retrieves a list of all devices in an Omada site.
---

# omada_devices (Data Source)

Retrieves a list of all devices (Access Points, Switches, Gateways) in an Omada site. This data source is useful for discovering devices, filtering by type, and dynamically managing device configurations.

## Example Usage

### Get All Devices

```terraform
data "omada_devices" "all" {}

output "device_count" {
  value = length(data.omada_devices.all.devices)
}
```

### Filter Access Points

```terraform
data "omada_devices" "all" {}

locals {
  access_points = [
    for device in data.omada_devices.all.devices :
    device if device.type == "ap"
  ]
}

output "ap_names" {
  value = [for ap in local.access_points : ap.name]
}
```

### Find Devices Needing Adoption

```terraform
data "omada_devices" "all" {}

output "pending_adoption" {
  value = [
    for device in data.omada_devices.all.devices :
    {
      mac   = device.mac
      model = device.model
    }
    if device.need_adopt
  ]
}
```

### Device Statistics

```terraform
data "omada_devices" "all" {}

output "device_summary" {
  value = {
    total    = length(data.omada_devices.all.devices)
    aps      = length([for d in data.omada_devices.all.devices : d if d.type == "ap"])
    switches = length([for d in data.omada_devices.all.devices : d if d.type == "switch"])
    gateways = length([for d in data.omada_devices.all.devices : d if d.type == "gateway"])
    online   = length([for d in data.omada_devices.all.devices : d if d.status == "connected"])
  }
}
```

### Dynamically Manage Devices

```terraform
data "omada_devices" "all" {}

# Configure all access points
resource "omada_device" "aps" {
  for_each = {
    for device in data.omada_devices.all.devices :
    device.mac => device
    if device.type == "ap"
  }

  mac         = each.value.mac
  name        = "AP - ${each.value.model}"
  led_enabled = false
  location    = "Auto-managed"
}
```

## Schema

### Optional

- `site_id` (String) Site to query devices from. If not specified, uses the provider's configured `site_id`

### Read-Only

- `devices` (List of Object) List of device objects with the following attributes:
  - `mac` (String) MAC address of the device
  - `name` (String) Name of the device
  - `type` (String) Device type: `ap`, `switch`, or `gateway`
  - `model` (String) Device model number
  - `status` (String) Connection status (e.g., `connected`, `disconnected`)
  - `led_enabled` (Boolean) Whether LED is enabled
  - `location` (String) Physical location description
  - `site` (String) Site name where device is located
  - `ip` (String) IP address of the device
  - `uptime` (Number) Device uptime in seconds
  - `firmware_version` (String) Current firmware version
  - `need_adopt` (Boolean) Whether device needs adoption

## Common Patterns

### Filter by Device Type

```terraform
# Access Points only
locals {
  aps = [for d in data.omada_devices.all.devices : d if d.type == "ap"]
}

# Switches only
locals {
  switches = [for d in data.omada_devices.all.devices : d if d.type == "switch"]
}

# Gateways only
locals {
  gateways = [for d in data.omada_devices.all.devices : d if d.type == "gateway"]
}
```

### Filter by Status

```terraform
# Online devices only
locals {
  online_devices = [
    for d in data.omada_devices.all.devices :
    d if d.status == "connected"
  ]
}

# Offline devices
locals {
  offline_devices = [
    for d in data.omada_devices.all.devices :
    d if d.status != "connected"
  ]
}
```

### Filter by Model

```terraform
# Specific model
locals {
  eap650_aps = [
    for d in data.omada_devices.all.devices :
    d if can(regex("EAP650", d.model))
  ]
}
```

### Filter by Location

```terraform
# Devices in specific location
locals {
  lobby_devices = [
    for d in data.omada_devices.all.devices :
    d if can(regex("(?i)lobby", d.location))
  ]
}
```

## Notes

- The data source retrieves all devices in the specified site on every refresh
- Device status is real-time and reflects the current connection state
- Use filters and for-expressions to work with specific subsets of devices
- The `need_adopt` attribute indicates if a device is pending adoption
- Firmware versions may be empty for offline devices
