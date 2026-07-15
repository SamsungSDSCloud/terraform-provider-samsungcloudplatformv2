provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_budget_budgets" "ids" {
  filter {
    name = var.budget_filter_name
    values = var.budget_filter_values
    use_regex = var.budget_filter_use_regex
  }
}
