provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_addressbookmembers" "addressbookmembers" {
  addrbook_id = var.AddrbookId
}
