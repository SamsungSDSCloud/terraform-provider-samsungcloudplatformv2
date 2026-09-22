provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_identity_center_user" "example" {
  instance_id = "efi74xj2n61q"
  id          = var.id
}