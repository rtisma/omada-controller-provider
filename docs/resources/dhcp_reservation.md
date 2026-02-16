---
page_title: "omada_dhcp_reservation Resource - terraform-provider-omada"
subcategory: ""
description: |-
  Manages a DHCP reservation (static IP assignment by MAC address).
---

# omada_dhcp_reservation (Resource)

Manages a DHCP reservation in Omada Controller. DHCP reservations assign a fixed IP address to a specific device based on its MAC address.

## Example Usage

### Basic Reservation

```terraform
resource "omada_dhcp_reservation" "printer" {
  name        = "Office Printer"
  mac_address = "AA:BB:CC:DD:EE:FF"
  ip_address  = "192.168.1.50"
  network_id  = omada_network.main.id
}
```

### Reservation with Comment

```terraform
resource "omada_dhcp_reservation" "camera" {
  name        = "Security Camera - Entrance"
  mac_address = "11:22:33:44:55:66"
  ip_address  = "192.168.50.10"
  network_id  = omada_network.iot.id
  comment     = "Front entrance camera - POE port 8"
}
```

### Multiple Reservations

```terraform
locals {
  servers = {
    "nas"  = { mac = "AA:11:BB:22:CC:33", ip = "192.168.10.10" }
    "mail" = { mac = "AA:11:BB:22:CC:44", ip = "192.168.10.11" }
    "web"  = { mac = "AA:11:BB:22:CC:55", ip = "192.168.10.12" }
  }
}

resource "omada_dhcp_reservation" "servers" {
  for_each = local.servers

  name        = "${each.key} server"
  mac_address = each.value.mac
  ip_address  = each.value.ip
  network_id  = omada_network.servers.id
}
```

## Schema

### Required

- `name` (String) Descriptive name for the device
- `mac_address` (String) MAC address in the format `AA:BB:CC:DD:EE:FF`
- `ip_address` (String) IP address to reserve for this device
- `network_id` (String) ID of the network where this reservation applies

### Optional

- `site_id` (String) Site where the reservation belongs. Defaults to the provider's `site_id`
- `comment` (String) Additional notes about this reservation

### Read-Only

- `id` (String) Unique identifier of the reservation

## Import

DHCP reservations can be imported using their ID:

```shell
terraform import omada_dhcp_reservation.printer <reservation-id>
```

## Notes

- The IP address must be within the network's subnet but outside the DHCP range to avoid conflicts
- MAC addresses should be in uppercase with colons as separators
- Changes to `mac_address` or `network_id` will force resource replacement
- The device will receive the reserved IP on its next DHCP renewal
