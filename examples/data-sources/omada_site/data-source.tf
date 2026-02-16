data "omada_site" "default" {
  name = "Default"
}

output "site_id" {
  value = data.omada_site.default.id
}

output "site_info" {
  value = {
    name     = data.omada_site.default.name
    type     = data.omada_site.default.type
    location = data.omada_site.default.location
    timezone = data.omada_site.default.timezone
  }
}
