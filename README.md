# Terraform Provider for Omada Controller

This Terraform provider enables you to manage your TP-Link Omada Controller infrastructure as code, supporting Omada Controller 5.x.

## Features

- **Network Management**: Create and manage VLANs and networks
- **WiFi Configuration**: Configure SSIDs with security settings and VLAN assignments
- **DHCP Management**: Static IP reservations and DHCP configuration
- **Device Management**: Adopt and configure Omada devices (APs, switches, gateways)
- **Site Information**: Retrieve site details and configuration

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for building from source)
- TP-Link Omada Controller 5.x

## Installation

### Using Terraform Registry (Coming Soon)

```hcl
terraform {
  required_providers {
    omada = {
      source  = "your-org/omada"
      version = "~> 0.1"
    }
  }
}
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/your-org/terraform-provider-omada.git
cd terraform-provider-omada
```

2. Build the provider:
```bash
go build -o terraform-provider-omada
```

3. Install the provider locally:
```bash
# Create the plugins directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/your-org/omada/0.1.0/darwin_arm64

# Copy the binary (adjust path for your OS/architecture)
cp terraform-provider-omada ~/.terraform.d/plugins/registry.terraform.io/your-org/omada/0.1.0/darwin_arm64/
```

For Linux (amd64):
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/your-org/omada/0.1.0/linux_amd64
cp terraform-provider-omada ~/.terraform.d/plugins/registry.terraform.io/your-org/omada/0.1.0/linux_amd64/
```

## Usage

### Provider Configuration

```hcl
provider "omada" {
  host     = "https://192.168.1.1:8043"
  username = "admin"
  password = "password"
  site_id  = "Default"
  insecure = true  # Set to true for self-signed certificates
}
```

### Configuration Options

| Name | Description | Required | Default |
|------|-------------|----------|---------|
| `host` | URL of the Omada Controller | Yes | - |
| `username` | Admin username | Yes | - |
| `password` | Admin password | Yes | - |
| `site_id` | Site name to manage | No | "Default" |
| `insecure` | Skip TLS verification | No | false |

### Example: Get Site Information

```hcl
data "omada_site" "default" {
  name = "Default"
}

output "site_id" {
  value = data.omada_site.default.id
}
```

## Development Status

This provider is under active development. The following features are implemented:

- [x] Phase 1: Foundation & Authentication
  - [x] API client with session management
  - [x] Provider configuration
  - [x] Site data source

Upcoming features:

- [ ] Phase 2: Network/VLAN resources
- [ ] Phase 3: WiFi/SSID resources
- [ ] Phase 4: DHCP reservation resources
- [ ] Phase 5: Device management
- [ ] Phase 6: Tests and documentation

## Authentication

The provider uses username/password authentication and maintains a session with the Omada Controller. Sessions are automatically refreshed when they expire.

**Security Best Practices:**
- Never commit credentials to source control
- Use environment variables or Terraform Cloud/Enterprise for sensitive data
- Consider using Terraform's built-in encryption for state files

## Documentation

### Quick Links

- **[Quick Start Guide](QUICK_START.md)** - Get up and running in 5 minutes
- **[Complete Resource Reference](RESOURCES.md)** - Comprehensive documentation for all resources and data sources
- **[Provider Documentation](docs/index.md)** - Provider configuration and overview
- **[Examples Directory](examples/)** - Working examples for all features

### Resource Documentation

Individual resource documentation with detailed examples:

- **Resources:**
  - [omada_network](docs/resources/network.md) - Manage VLANs and networks
  - [omada_ssid](docs/resources/ssid.md) - Manage wireless networks
  - [omada_dhcp_reservation](docs/resources/dhcp_reservation.md) - Static IP assignments
  - [omada_device](docs/resources/device.md) - Configure devices (APs, switches, gateways)

- **Data Sources:**
  - [omada_site](docs/data-sources/site.md) - Retrieve site information
  - [omada_devices](docs/data-sources/devices.md) - List and filter devices

## API Compatibility

This provider is designed for TP-Link Omada Controller 5.x and uses the `/api/v2/` endpoints. It has been tested with:
- Omada Controller 5.11+
- Omada Controller 5.12+

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Build: `go build`
4. Run tests: `go test ./...`

### Testing

To test the provider:

1. Set up a test Omada Controller instance
2. Create a `terraform.tfvars` file with your credentials:
```hcl
omada_host     = "https://your-controller:8043"
omada_username = "admin"
omada_password = "your-password"
```
3. Run `terraform init` and `terraform plan`

## License

This project is licensed under the MPL 2.0 License - see the LICENSE file for details.

## Acknowledgments

- HashiCorp for the Terraform Plugin Framework
- TP-Link for the Omada Controller API
- [go-omada](https://github.com/dougbw/go-omada) for API reference

## Support

For issues, questions, or contributions, please open an issue on [GitHub](https://github.com/your-org/terraform-provider-omada/issues).
