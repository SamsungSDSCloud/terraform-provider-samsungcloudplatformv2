provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_virtualserver_keypair" "keypair" {
  name = var.name
  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}
