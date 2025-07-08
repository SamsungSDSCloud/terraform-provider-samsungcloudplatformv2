provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_quota_account_quotas" "ids" {
  filter {
    name = var.account_quotas_filter_name
    values = var.account_quotas_filter_values
    use_regex = var.account_quotas_filter_use_regex
  }
}
