provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_policy" "policy" {
  id = var.id
}
