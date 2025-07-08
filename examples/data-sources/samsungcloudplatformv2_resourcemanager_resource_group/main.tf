provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_resourcemanager_resource_group" "resource_group" {
  tags = var.resource_group_tags
  filter {
    name = var.resource_group_filter_name
    values = var.resource_group_filter_values
    use_regex = var.resource_group_filter_use_regex
  }
}
