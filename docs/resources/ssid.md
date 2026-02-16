---
page_title: "omada_ssid Resource - terraform-provider-omada"
subcategory: ""
description: |-
  Manages a wireless network (SSID) in Omada Controller.
---

# omada_ssid (Resource)

Manages a wireless network (SSID) in Omada Controller. SSIDs can be configured with various security modes, VLAN assignments, bandwidth limits, and client restrictions.

## Example Usage

### WPA3 Personal Network

```terraform
resource "omada_ssid" "corporate" {
  name            = "Corporate WiFi"
  ssid            = "CorpNet"
  security_mode   = "wpa3-personal"
  password        = var.wifi_password
  vlan_id         = 10
  enabled         = true
  band_2_4g_enabled = true
  band_5g_enabled   = true
}
```

### Guest Network with Rate Limiting

```terraform
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
}
```

### Hidden IoT Network (2.4GHz Only)

```terraform
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
```

### Enterprise Network with RADIUS

```terraform
resource "omada_ssid" "enterprise" {
  name           = "Enterprise Network"
  ssid           = "Corp-Secure"
  security_mode  = "wpa2-enterprise"
  radius_profile = "MainRADIUS"
  vlan_id        = 20
  enabled        = true
}
```

## Schema

### Required

- `name` (String) Descriptive name of the SSID
- `ssid` (String) Network name broadcast to clients
- `security_mode` (String) Security mode. Must be one of: `wpa2-personal`, `wpa3-personal`, `wpa2-enterprise`, `wpa3-enterprise`, `open`

### Optional

- `site_id` (String) Site where the SSID belongs. Defaults to the provider's `site_id`
- `password` (String, Sensitive) WiFi password. Required for WPA2/WPA3 personal modes
- `enabled` (Boolean) Enable or disable the SSID. Defaults to `true`
- `hide_ssid` (Boolean) Hide SSID from broadcast. Defaults to `false`
- `vlan_id` (Number) VLAN ID for this SSID
- `guest_network` (Boolean) Mark as guest network. Defaults to `false`
- `client_isolation` (Boolean) Isolate clients from each other. Defaults to `false`
- `band_2_4g_enabled` (Boolean) Enable on 2.4GHz band. Defaults to `true`
- `band_5g_enabled` (Boolean) Enable on 5GHz band. Defaults to `true`
- `band_6g_enabled` (Boolean) Enable on 6GHz band. Defaults to `false`
- `max_clients` (Number) Maximum number of clients allowed
- `rate_limit` (Boolean) Enable bandwidth limiting. Defaults to `false`
- `downlink_limit` (Number) Download limit in Kbps. Requires `rate_limit = true`
- `uplink_limit` (Number) Upload limit in Kbps. Requires `rate_limit = true`
- `schedule_enabled` (Boolean) Enable scheduled activation. Defaults to `false`
- `portal_enabled` (Boolean) Enable captive portal. Defaults to `false`
- `radius_profile` (String) RADIUS profile name. Required for enterprise security modes

### Read-Only

- `id` (String) Unique identifier of the SSID

## Import

SSIDs can be imported using their ID:

```shell
terraform import omada_ssid.corporate <ssid-id>
```

## Notes

- The `password` attribute is sensitive and will not be displayed in plan output
- API responses do not return the password, so Terraform uses the value from state
- Band availability depends on your access point hardware capabilities
- Guest networks typically have restricted access to local network resources
