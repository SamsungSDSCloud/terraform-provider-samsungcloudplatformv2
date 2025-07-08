provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_virtualserver_server_group" "server_group" {
  name = var.name
  policy = var.policy
  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}
