provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_instances" "instances" {
  size = var.size
  page = var.page
  sort = var.sort
}
