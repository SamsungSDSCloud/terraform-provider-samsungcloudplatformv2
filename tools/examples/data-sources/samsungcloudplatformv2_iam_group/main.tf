provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_group" "group" {
  id = var.id
}
