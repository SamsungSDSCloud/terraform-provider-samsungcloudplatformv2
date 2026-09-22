output "baseline_assignments" {
  description = "List of baseline assignments matching the filter"
  value       = data.samsungcloudplatformv2_cloudcontrol_baseline_assignments.baseline_assignments.baseline_assignments
}