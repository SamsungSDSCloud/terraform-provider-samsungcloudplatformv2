provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_billing_planned_computes" "planned_computes" {
  limit = 10
  page = 1
}
