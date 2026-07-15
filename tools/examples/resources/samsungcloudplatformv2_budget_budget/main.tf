provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_budget_budget" "budget" {
  name             = var.budget_name
  amount           = var.budget_amount
  start_month      = var.budget_start_month
  unit             = var.budget_unit
  notifications    = var.budget_notifications
  prevention       = var.budget_prevention
}