# Quick Start Guide

Get started with the Terraform Provider for Omada Controller in 5 minutes.

## Prerequisites

- Terraform >= 1.0
- TP-Link Omada Controller 5.x
- Admin credentials for your controller

## Step 1: Install the Provider

### Local Installation (Development)

```bash
# Clone and build
git clone https://github.com/your-org/terraform-provider-omada.git
cd terraform-provider-omada
make install
```

### From Terraform Registry (Future)

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

## Step 2: Configure the Provider

Create a `main.tf` file:

```hcl
provider "omada" {
  host     = "https://192.168.1.1:8043"
  username = "admin"
  password = var.omada_password
  site_id  = "Default"
  insecure = true  # For self-signed certificates
}

variable "omada_password" {
  type      = string
  sensitive = true
}
```

Create a `terraform.tfvars` file (don't commit this!):

```hcl
omada_password = "your-actual-password"
```

## Step 3: Your First Resource

Add a network to your `main.tf`:

```hcl
resource "omada_network" "test" {
  name         = "Test Network"
  vlan_id      = 99
  gateway      = "192.168.99.1"
  netmask      = "255.255.255.0"
  dhcp_enabled = true
  dhcp_start   = "192.168.99.100"
  dhcp_end     = "192.168.99.200"
}
```

## Step 4: Deploy

```bash
# Initialize Terraform
terraform init

# Preview changes
terraform plan

# Apply changes
terraform apply
```

## Step 5: Verify

Log into your Omada Controller web interface and verify that the "Test Network" was created in Settings → Networks.

## What's Next?

### Create a WiFi Network

```hcl
resource "omada_ssid" "test_wifi" {
  name          = "Test WiFi"
  ssid          = "TestNetwork"
  security_mode = "wpa2-personal"
  password      = "SecurePass123!"
  vlan_id       = omada_network.test.vlan_id
}
```

### Add a DHCP Reservation

```hcl
resource "omada_dhcp_reservation" "printer" {
  name        = "Office Printer"
  mac_address = "AA:BB:CC:DD:EE:FF"
  ip_address  = "192.168.99.50"
  network_id  = omada_network.test.id
}
```

### List Your Devices

```hcl
data "omada_devices" "all" {}

output "devices" {
  value = {
    total = length(data.omada_devices.all.devices)
    aps   = length([
      for d in data.omada_devices.all.devices : d if d.type == "ap"
    ])
  }
}
```

## Common Commands

```bash
# Initialize/update provider
terraform init

# Format your configuration
terraform fmt

# Validate configuration
terraform validate

# Preview changes
terraform plan

# Apply changes
terraform apply

# Destroy resources
terraform destroy

# Show current state
terraform show

# Import existing resource
terraform import omada_network.main <network-id>
```

## Troubleshooting

### Certificate Errors

If you see TLS/certificate errors:
```hcl
provider "omada" {
  # ... other settings ...
  insecure = true  # Skip certificate verification
}
```

### Authentication Errors

Verify your credentials work by logging into the web interface first.

### Resources Not Found

Make sure you're using the correct site ID (default is "Default").

## Complete Examples

See the `examples/` directory for:
- `examples/provider/` - Provider configuration examples
- `examples/resources/` - Individual resource examples
- `examples/data-sources/` - Data source examples
- `examples/complete-setup/` - Full network setup example

## Documentation

- **[Complete Resource Reference](RESOURCES.md)** - Detailed documentation for all resources
- **[Individual Docs](docs/)** - Per-resource documentation files
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Changelog](CHANGELOG.md)** - Version history

## Getting Help

- 📖 Read the [RESOURCES.md](RESOURCES.md) for detailed documentation
- 💬 Open an [issue](https://github.com/your-org/terraform-provider-omada/issues) for bugs or questions
- 🔍 Check the [examples](examples/) directory for working configurations

## Next Steps

1. ✅ You've created your first resource!
2. 📚 Read the [complete documentation](RESOURCES.md)
3. 🚀 Check out the [complete setup example](examples/complete-setup/)
4. 🔧 Build your network configuration
5. 📝 Share your feedback!
