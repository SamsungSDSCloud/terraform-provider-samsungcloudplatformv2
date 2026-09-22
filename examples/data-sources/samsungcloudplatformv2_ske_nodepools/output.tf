output "nodepools" {
  value = data.samsungcloudplatformv2_ske_nodepools.nodepools
}

output "nodepool_subnet_ids" {
  description = "Subnet IDs of all nodepools in the cluster"
  value       = data.samsungcloudplatformv2_ske_nodepools.nodepools.nodepools != null ? [for np in data.samsungcloudplatformv2_ske_nodepools.nodepools.nodepools : np.subnet_id] : []
}
