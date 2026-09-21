output "tags_id" {
  value = data.samsungcloudplatformv2_scr_tag_secrets.secrets.tags_id
}

output "filtered_count" {
  value = data.samsungcloudplatformv2_scr_tag_secrets.secrets.filtered_count
}

output "last_scanned_at" {
  value = data.samsungcloudplatformv2_scr_tag_secrets.secrets.last_scanned_at
}

output "release_version" {
  value = data.samsungcloudplatformv2_scr_tag_secrets.secrets.release_version
}

output "secrets" {
  value       = data.samsungcloudplatformv2_scr_tag_secrets.secrets.secrets
  sensitive   = true
}

output "secret_summary" {
  value = data.samsungcloudplatformv2_scr_tag_secrets.secrets.secret_summary
}
