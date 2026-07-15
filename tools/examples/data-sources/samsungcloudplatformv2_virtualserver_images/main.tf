  provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_images" "ids" {
  scp_image_type = var.scp_image_type
  scp_original_image_type = var.scp_original_image_type
  name = var.name
  os_distro = var.os_distro
  status = var.status
  visibility = var.visibility
  filter {
    name = var.images_filter_name
    values = var.images_filter_values
    use_regex = var.images_filter_use_regex
  }
}
