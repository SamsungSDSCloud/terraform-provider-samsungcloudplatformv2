provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_resourcemanager_tags" "tags" {
  key = "tag1"
  value = "11"
}
