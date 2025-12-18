provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_certificate_manager" "certificate" {
        id = var.id
}
