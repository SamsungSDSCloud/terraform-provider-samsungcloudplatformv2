output "ids" {
  value = data.samsungcloudplatformv2_budget_budgets.budgets.ids
}


output "budgets" {
  value = data.samsungcloudplatformv2_budget_budgets.budgets.budgets
}

output "count" {
  description = "Total number of budgets matching the filter"
  value       = data.samsungcloudplatformv2_budget_budgets.budgets.budget_count
}
output "sort" {
  description = "Sort criteria actually applied by the server"
  value       = data.samsungcloudplatformv2_budget_budgets.budgets.sort
}
output "size" {
  description = "Page size actually applied by the server"
  value       = data.samsungcloudplatformv2_budget_budgets.budgets.size
}
output "page" {
  description = "Page number actually applied by the server"
  value       = data.samsungcloudplatformv2_budget_budgets.budgets.page
}
