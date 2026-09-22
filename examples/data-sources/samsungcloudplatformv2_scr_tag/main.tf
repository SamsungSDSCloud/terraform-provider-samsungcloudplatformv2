provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_scr_tag" "tag" {
    id = var.id
}
