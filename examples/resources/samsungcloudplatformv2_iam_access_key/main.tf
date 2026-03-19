provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_iam_access_key" "access_key" {
  access_key_type = var.access_key_access_key_type
  description = var.access_key_description
  is_enabled = var.access_key_is_enabled
}
