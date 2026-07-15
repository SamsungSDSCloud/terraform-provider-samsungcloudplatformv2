provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_gslb_gslb" "gslb" {
    id = var.id
}
