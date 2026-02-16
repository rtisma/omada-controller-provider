---
page_title: "omada_network Resource - terraform-provider-omada"
subcategory: ""
description: |-
  Manages an Omada network/VLAN with DHCP configuration.
---

# omada_network (Resource)

Manages a network/VLAN in Omada Controller. Networks are Layer 3 constructs that provide IP addressing, DHCP services, and network isolation via VLANs.

## Example Usage

### Basic Network with DHCP

```terraform
resource "omada_network" "main" {
  name         = "Main Network"
  vlan_id      = 1
  gateway      = "192.168.1.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.1.100"
  dhcp_end     = "192.168.1.200"
  lease_time   = 1440
  dns_primary  = "8.8.8.8"
  dns_secondary = "8.8.4.4"
}
```

### Guest Network

```terraform
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
```

### Static IP Network (No DHCP)

```terraform
resource "omada_network" "servers" {
  name         = "Server Network"
  vlan_id      = 10
  gateway      = "192.168.10.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = false
}
```

## Schema

### Required

- `name` (String) Name of the network
- `vlan_id` (Number) VLAN ID (1-4094)
- `gateway` (String) Gateway IP address
- `netmask` (String) Subnet mask (e.g., `255.255.255.0`)

### Optional

- `site_id` (String) Site where the network belongs. Defaults to the provider's `site_id`
- `dhcp_enabled` (Boolean) Enable DHCP server. Defaults to `false`
- `dhcp_start` (String) Start of DHCP range. Required if `dhcp_enabled` is true
- `dhcp_end` (String) End of DHCP range. Required if `dhcp_enabled` is true
- `lease_time` (Number) DHCP lease time in minutes. Defaults to `1440` (24 hours)
- `dns_primary` (String) Primary DNS server
- `dns_secondary` (String) Secondary DNS server
- `domain_name` (String) Domain name for DHCP
- `purpose` (String) Network purpose/category (e.g., `general`, `guest`, `iot`)

### Read-Only

- `id` (String) Unique identifier of the network

## Import

Networks can be imported using their ID:

```shell
terraform import omada_network.main <network-id>
```

To find the network ID, you can inspect the Omada Controller API or use browser developer tools.
