provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_role" "role" {
  id = var.id
}
