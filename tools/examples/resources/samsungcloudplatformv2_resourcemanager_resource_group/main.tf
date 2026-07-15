provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_resourcemanager_resource_group" "resource_group" {
  name        = var.resource_group_name
  description = var.resource_group_description
  resource_types = ["iam:policy"]
  group_definition_tags = {
    "testkey1" : "testvalue1"
  }
  tags = {
    "tf_key1": "tf_val1"
  }
}
