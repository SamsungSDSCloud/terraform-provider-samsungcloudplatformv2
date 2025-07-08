provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
  limit = var.access_key_limit
}
