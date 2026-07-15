provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_quota_account_quota" "account_quota"{
  id = var.account_quota_account_quota_id
}
