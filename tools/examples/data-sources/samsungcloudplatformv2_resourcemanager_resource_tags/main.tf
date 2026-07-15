provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_resourcemanager_resource_tags" "resource_tags" {
  srn = "srn:dev2::afd580f490394896a6bceabf77683c6b:kr-west1::scp-network:vpc/b75d3a1aae2d4fd6aee4ac49ba47e681"
}
