provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_baremetal_baremetal" "baremetal" {
  id = var.id
}