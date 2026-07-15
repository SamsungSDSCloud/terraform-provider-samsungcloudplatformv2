provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_groups" "groups" {
  size = var.size
  page = var.page
  sort = var.sort
  name = var.name
}
