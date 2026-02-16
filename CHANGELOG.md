# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-15

### Added
- Initial release of the Terraform Provider for Omada Controller
- Provider configuration with host, username, password, site_id, and insecure options
- Authentication with Omada Controller 5.x API using session-based auth
- Data sources:
  - `omada_site` - Retrieve site information
  - `omada_devices` - List all devices (APs, switches, gateways)
- Resources:
  - `omada_network` - Manage VLANs and networks with DHCP configuration
  - `omada_ssid` - Manage WiFi networks (SSIDs) with security settings
  - `omada_dhcp_reservation` - Manage DHCP reservations (static IP assignments)
  - `omada_device` - Manage device configuration (name, LED, location)
- Comprehensive examples for all resources and data sources
- Complete README with installation and usage instructions

### Known Limitations
- API documentation from TP-Link is limited; some endpoints discovered through reverse engineering
- Import functionality needs to be tested with real Omada Controller
- Some advanced features (inter-VLAN routing, advanced RADIUS configuration) may have limited support

## [Unreleased]

### Planned
- Enhanced error handling and retry logic
- Support for portal configuration
- ACL (Access Control List) management
- Port forwarding rules
- Advanced device settings (radio power, channel selection)
- Comprehensive acceptance tests
- Auto-generated documentation with tfplugindocs
