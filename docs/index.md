---
page_title: "Omada Provider"
subcategory: ""
description: |-
  Terraform provider for managing TP-Link Omada Controller infrastructure.
---

# Omada Provider

The Omada provider enables you to manage your TP-Link Omada Controller infrastructure as code. This provider supports Omada Controller 5.x and allows you to manage networks, WiFi configurations, DHCP reservations, and devices.

## Features

- **Network Management**: Create and manage VLANs and networks with DHCP configuration
- **WiFi Configuration**: Configure SSIDs with various security modes and VLAN assignments
- **DHCP Management**: Static IP reservations by MAC address
- **Device Management**: Configure and adopt Omada devices (APs, switches, gateways)

## Example Usage

```terraform
terraform {
  required_providers {
    omada = {
      source = "your-org/omada"
    }
  }
}

provider "omada" {
  host     = "https://192.168.1.1:8043"
  username = "admin"
  password = var.omada_password
  site_id  = "Default"
  insecure = true  # For self-signed certificates
}

# Create a network
resource "omada_network" "main" {
  name         = "Main Network"
  vlan_id      = 1
  gateway      = "192.168.1.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.1.100"
  dhcp_end     = "192.168.1.200"
}

# Create a WiFi network
resource "omada_ssid" "corporate" {
  name          = "Corporate WiFi"
  ssid          = "CorpNet"
  security_mode = "wpa3-personal"
  password      = var.wifi_password
  vlan_id       = omada_network.main.vlan_id
}
```

## Authentication

The provider uses username and password authentication with automatic session management. Sessions are maintained for approximately 14 days and are automatically refreshed when they expire.

## Schema

### Required

- `host` (String) The URL of your Omada Controller (e.g., `https://192.168.1.1:8043`)
- `username` (String) Admin username for authentication
- `password` (String, Sensitive) Admin password for authentication

### Optional

- `site_id` (String) The site name to manage. Defaults to `"Default"`
- `insecure` (Boolean) Skip TLS certificate verification. Useful for self-signed certificates. Defaults to `false`

## API Compatibility

This provider is designed for:
- TP-Link Omada Controller 5.x
- API version 2 (`/api/v2/` endpoints)

Tested with:
- Omada Controller 5.11+
- Omada Controller 5.12+

## Security Considerations

- Never commit credentials to source control
- Use Terraform variables or environment variables for sensitive data
- Consider using Terraform Cloud/Enterprise for secure state storage
- Enable `insecure = true` only for development/testing with self-signed certificates

## Limitations

- Some advanced features may have limited API support
- API documentation from TP-Link is limited; some endpoints discovered through reverse engineering
- Inter-VLAN routing configuration may require manual setup in the controller UI

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/your-org/terraform-provider-omada/issues
- Documentation: See the [RESOURCES.md](../RESOURCES.md) file for detailed reference
