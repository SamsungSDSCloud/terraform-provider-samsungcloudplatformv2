provider "samsungcloudplatformv2" {
}

# Assign baseline to an Organization Unit
resource "samsungcloudplatformv2_cloudcontrol_baseline_assignment" "ou" {
  assignment_id   = "ou-3b8908b1c0c14f9f9c269f7d1784258c"
  landing_zone_id = "aebc7c53544748fcbfbf5f2904c307e9"
  resource_type   = "OU"
  agree_yn        = "Y"

  # Change this value to trigger re-registration
  reregister_trigger = "1"
}

# Assign baseline to an Account
resource "samsungcloudplatformv2_cloudcontrol_baseline_assignment" "account" {
  assignment_id      = "e15ed1e264d447e9b6c61c03c1499091"
  landing_zone_id    = "aebc7c53544748fcbfbf5f2904c307e9"
  resource_type      = "ACCOUNT"
  agree_yn           = "Y"
  parent_unit_id     = "ou-fc8c29a138d78e24bf1fa86812fc8b"
  sso_user_name      = "testuser"
  sso_user_real_name = "test user"
}