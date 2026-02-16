---
page_title: "omada_site Data Source - terraform-provider-omada"
subcategory: ""
description: |-
  Retrieves information about an Omada site.
---

# omada_site (Data Source)

Retrieves information about an Omada site. Sites are the top-level organizational unit in Omada Controller, containing devices, networks, and configurations.

## Example Usage

### Get Default Site

```terraform
data "omada_site" "default" {
  name = "Default"
}

output "site_id" {
  value = data.omada_site.default.id
}
```

### Use Provider's Configured Site

```terraform
# Uses the site_id from provider configuration
data "omada_site" "current" {}

output "timezone" {
  value = data.omada_site.current.timezone
}
```

### Get Specific Site

```terraform
data "omada_site" "branch_office" {
  name = "Branch-NYC"
}
```

## Schema

### Optional

- `name` (String) Name of the site to retrieve. If not specified, uses the provider's configured `site_id`

### Read-Only

- `id` (String) Unique identifier of the site
- `name` (String) Name of the site
- `type` (String) Type of the site
- `location` (String) Physical location of the site
- `timezone` (String) Timezone configuration for the site
- `scenario` (String) Site scenario/template

## Notes

- If no `name` is specified, the data source uses the `site_id` from the provider configuration
- Site names are case-sensitive
- The `Default` site is created automatically in new Omada Controller installations
