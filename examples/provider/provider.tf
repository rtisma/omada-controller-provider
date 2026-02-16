terraform {
  required_providers {
    omada = {
      source  = "your-org/omada"
      version = "~> 0.1"
    }
  }
}

provider "omada" {
  host     = "https://192.168.1.1:8043"
  username = "admin"
  password = "password"
  site_id  = "Default"
  insecure = true # Set to true to skip TLS verification for self-signed certificates
}
