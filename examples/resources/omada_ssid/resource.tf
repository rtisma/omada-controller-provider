resource "omada_ssid" "corporate_wifi" {
  name            = "Corporate WiFi"
  ssid            = "CorpNet"
  security_mode   = "wpa3-personal"
  password        = "SecurePassword123!"
  vlan_id         = 10
  enabled         = true
  hide_ssid       = false
  guest_network   = false
  client_isolation = false
  band_2_4g_enabled = true
  band_5g_enabled   = true
  band_6g_enabled   = false
}

resource "omada_ssid" "guest_wifi" {
  name             = "Guest WiFi"
  ssid             = "Guest-Network"
  security_mode    = "wpa2-personal"
  password         = "GuestPass2024"
  vlan_id          = 100
  enabled          = true
  guest_network    = true
  client_isolation = true
  band_2_4g_enabled = true
  band_5g_enabled   = true
  max_clients      = 50
  rate_limit       = true
  downlink_limit   = 10000  # 10 Mbps
  uplink_limit     = 5000   # 5 Mbps
}

resource "omada_ssid" "iot_wifi" {
  name              = "IoT Devices"
  ssid              = "IoT-Network"
  security_mode     = "wpa2-personal"
  password          = "IoTPassword!"
  vlan_id           = 50
  enabled           = true
  hide_ssid         = true
  client_isolation  = true
  band_2_4g_enabled = true
  band_5g_enabled   = false
}

resource "omada_ssid" "enterprise_wifi" {
  name          = "Enterprise Network"
  ssid          = "EnterpriseCorp"
  security_mode = "wpa2-enterprise"
  radius_profile = "main-radius"
  vlan_id       = 20
  enabled       = true
}
