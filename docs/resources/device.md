---
page_title: "omada_device Resource - terraform-provider-omada"
subcategory: ""
description: |-
  Manages configuration of an Omada device (Access Point, Switch, or Gateway).
---

# omada_device (Resource)

Manages configuration of an Omada device. This resource configures existing devices that are already connected to the Omada Controller.

**Important:** This resource does not create new devices. Devices must be physically connected to your network and visible in the Omada Controller before they can be managed by Terraform.

## Example Usage

### Configure a Single Device

```terraform
resource "omada_device" "lobby_ap" {
  mac         = "AA:BB:CC:DD:EE:FF"
  name        = "Lobby Access Point"
  led_enabled = false
  location    = "Main Lobby - 1st Floor"
}
```

### Configure Multiple Devices

```terraform
resource "omada_device" "conference_aps" {
  for_each = {
    "conf_a" = { mac = "11:22:33:44:55:66", location = "Conference Room A" }
    "conf_b" = { mac = "11:22:33:44:55:77", location = "Conference Room B" }
  }

  mac         = each.value.mac
  name        = "AP - ${upper(each.key)}"
  led_enabled = true
  location    = each.value.location
}
```

### Manage All Access Points

```terraform
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
  location    = "Managed by Terraform"
}
```

## Schema

### Required

- `mac` (String) MAC address of the device in format `AA:BB:CC:DD:EE:FF`. Changes to this force resource replacement
- `name` (String) Name/description of the device

### Optional

- `site_id` (String) Site where the device belongs. Defaults to the provider's `site_id`
- `led_enabled` (Boolean) Enable or disable LED indicator. Defaults to `true`
- `location` (String) Physical location description

### Read-Only

- `type` (String) Device type: `ap` (Access Point), `switch`, or `gateway`
- `model` (String) Device model number
- `status` (String) Connection status

## Import

Devices can be imported using their MAC address:

```shell
terraform import omada_device.lobby_ap AA:BB:CC:DD:EE:FF
```

## Notes

### Device Adoption

If a device is pending adoption when you create this resource, Terraform will automatically adopt it. This means the device will be added to the controller's management.

### Device Deletion

When you destroy this resource (`terraform destroy`), the device will be "forgotten" by the controller, removing it from management. The physical device will remain on the network but will need to be re-adopted if you want to manage it again.

### Finding Device MAC Addresses

You can find device MAC addresses in several ways:

1. Using the `omada_devices` data source
2. In the Omada Controller web interface
3. On the physical device label
4. Using network scanning tools

### Best Practices

- Use the `omada_devices` data source to dynamically discover and manage devices
- Keep LEDs disabled on APs in visible areas to reduce distraction
- Use descriptive location strings that match your physical site documentation
- Consider using `for_each` to manage multiple similar devices
