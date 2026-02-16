# Documentation Index

Complete index of all documentation for the Terraform Provider for Omada Controller.

## 📚 Main Documentation

### Getting Started
- **[README.md](README.md)** - Main project overview and quick introduction
- **[QUICK_START.md](QUICK_START.md)** - 5-minute guide to get started
- **[RESOURCES.md](RESOURCES.md)** - Complete reference for all resources and data sources

### Project Information
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - How to contribute to the project
- **[LICENSE](LICENSE)** - MPL 2.0 license text

## 📖 Provider Documentation

### Provider Configuration
- **[docs/index.md](docs/index.md)** - Provider overview, configuration, and authentication

## 🔧 Resource Documentation

Individual resource documentation with schemas, examples, and usage patterns:

### Networks & VLANs
- **[docs/resources/network.md](docs/resources/network.md)**
  - Create and manage networks/VLANs
  - DHCP configuration
  - DNS settings
  - Network purposes (general, guest, iot)

### Wireless Networks
- **[docs/resources/ssid.md](docs/resources/ssid.md)**
  - WiFi SSID configuration
  - Security modes (WPA2/WPA3, Personal/Enterprise, Open)
  - VLAN assignment
  - Rate limiting and client restrictions
  - Band selection (2.4GHz, 5GHz, 6GHz)

### IP Management
- **[docs/resources/dhcp_reservation.md](docs/resources/dhcp_reservation.md)**
  - Static IP assignments by MAC address
  - DHCP reservations
  - Network association

### Device Management
- **[docs/resources/device.md](docs/resources/device.md)**
  - Configure Access Points, Switches, and Gateways
  - Device adoption
  - LED control
  - Location tracking

## 📊 Data Source Documentation

### Site Information
- **[docs/data-sources/site.md](docs/data-sources/site.md)**
  - Retrieve site details
  - Timezone information
  - Site type and location

### Device Discovery
- **[docs/data-sources/devices.md](docs/data-sources/devices.md)**
  - List all devices in a site
  - Filter by device type (AP, switch, gateway)
  - Check adoption status
  - Device statistics and firmware versions

## 💡 Examples

### Provider Configuration
- **[examples/provider/provider.tf](examples/provider/provider.tf)**
  - Basic provider setup
  - Authentication configuration
  - Site selection

### Data Source Examples
- **[examples/data-sources/omada_site/](examples/data-sources/omada_site/)**
  - Site information retrieval
  - Using site data in configurations

- **[examples/data-sources/omada_devices/](examples/data-sources/omada_devices/)**
  - Listing all devices
  - Filtering devices by type
  - Device statistics

### Resource Examples

#### Networks
- **[examples/resources/omada_network/resource.tf](examples/resources/omada_network/resource.tf)**
  - Guest network with DHCP
  - IoT network configuration
  - Management network (static IPs)
  - Multiple VLAN examples

#### Wireless Networks
- **[examples/resources/omada_ssid/resource.tf](examples/resources/omada_ssid/resource.tf)**
  - Corporate WiFi (WPA3)
  - Guest WiFi with rate limiting
  - IoT WiFi (2.4GHz only, hidden)
  - Enterprise WiFi with RADIUS

#### DHCP Reservations
- **[examples/resources/omada_dhcp_reservation/resource.tf](examples/resources/omada_dhcp_reservation/resource.tf)**
  - Printer reservations
  - Server static IPs
  - Security camera IPs
  - Bulk reservations with for_each

#### Device Management
- **[examples/resources/omada_device/resource.tf](examples/resources/omada_device/resource.tf)**
  - Single device configuration
  - Multiple devices with for_each
  - Dynamic device management from data source

### Complete Setup
- **[examples/complete-setup/main.tf](examples/complete-setup/main.tf)**
  - Full working example
  - Multiple networks (main, guest, IoT)
  - Multiple SSIDs with different security levels
  - DHCP reservations
  - Device discovery
  - Comprehensive outputs

- **[examples/complete-setup/terraform.tfvars.example](examples/complete-setup/terraform.tfvars.example)**
  - Example variable values
  - Template for your own configuration

## 📋 Documentation Features

### In Each Resource Doc
✅ Complete schema reference with all attributes
✅ Required vs optional arguments clearly marked
✅ Multiple real-world usage examples
✅ Import instructions
✅ Common patterns and best practices
✅ Notes about limitations and behaviors

### In RESOURCES.md
✅ Provider configuration reference
✅ All resources with complete schemas
✅ All data sources with complete schemas
✅ Extensive examples for each resource
✅ Common patterns section
✅ Troubleshooting guide
✅ Import instructions

### In Individual Docs (docs/)
✅ Terraform Registry-compatible format
✅ YAML frontmatter for metadata
✅ Organized by resource/data-source
✅ Ready for tfplugindocs generation

## 🎯 Quick Reference by Use Case

### "I want to create a new VLAN"
→ [docs/resources/network.md](docs/resources/network.md)
→ [examples/resources/omada_network/](examples/resources/omada_network/)

### "I want to set up WiFi"
→ [docs/resources/ssid.md](docs/resources/ssid.md)
→ [examples/resources/omada_ssid/](examples/resources/omada_ssid/)

### "I want to assign static IPs"
→ [docs/resources/dhcp_reservation.md](docs/resources/dhcp_reservation.md)
→ [examples/resources/omada_dhcp_reservation/](examples/resources/omada_dhcp_reservation/)

### "I want to manage my access points"
→ [docs/resources/device.md](docs/resources/device.md)
→ [docs/data-sources/devices.md](docs/data-sources/devices.md)
→ [examples/resources/omada_device/](examples/resources/omada_device/)

### "I want to see what devices I have"
→ [docs/data-sources/devices.md](docs/data-sources/devices.md)
→ [examples/data-sources/omada_devices/](examples/data-sources/omada_devices/)

### "I want a complete working example"
→ [examples/complete-setup/](examples/complete-setup/)
→ [QUICK_START.md](QUICK_START.md)

### "I'm getting started for the first time"
→ [QUICK_START.md](QUICK_START.md)
→ [README.md](README.md)

### "I want detailed API documentation"
→ [RESOURCES.md](RESOURCES.md)

## 📝 Documentation Statistics

- **7** Markdown documentation files in docs/
- **1** Comprehensive RESOURCES.md reference (15,000+ words)
- **1** Quick start guide
- **1** Contributing guide
- **10+** Terraform example files
- **4** Resources documented
- **2** Data sources documented
- **30+** Code examples across all docs

## 🔄 Keeping Documentation Updated

When adding new resources or features:

1. Update [RESOURCES.md](RESOURCES.md) with complete reference
2. Create individual doc in `docs/resources/` or `docs/data-sources/`
3. Add example in `examples/resources/` or `examples/data-sources/`
4. Update [README.md](README.md) if it's a major feature
5. Add entry to [CHANGELOG.md](CHANGELOG.md)
6. Update this index if needed

## 🎨 Documentation Style Guide

- Use clear, concise language
- Provide real-world examples
- Include both simple and complex use cases
- Document limitations and gotchas
- Use code blocks with syntax highlighting
- Include imports instructions where applicable
- Reference related resources
- Keep examples self-contained and runnable

## 🔗 External References

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Terraform Registry Documentation Standards](https://developer.hashicorp.com/terraform/registry/providers/docs)
- [Omada API Examples](https://gist.github.com/mbentley/03c198077c81d52cb029b825e9a6dc18)

---

**Last Updated:** 2026-02-15
**Provider Version:** 0.1.0
