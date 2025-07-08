provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_resourcemanager_resource_groups" "ids" {
  region = var.resource_groups_region
  filter {
    name = var.resource_groups_filter_name
    values = var.resource_groups_filter_values
    use_regex = var.resource_groups_filter_use_regex
  }
}
