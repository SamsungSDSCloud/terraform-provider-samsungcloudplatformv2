provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_certificate_managers" "certificates" {
        size = var.size
        sort = var.sort
        page = var.page
        is_mine = var.is_mine
        name = var.name
        cn = var.cn
        state = var.state
}
