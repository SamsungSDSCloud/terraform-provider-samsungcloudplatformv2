output "id" {
  value = data.samsungcloudplatformv2_scr_tag.tag.id
}

output "image_id" {
  value = data.samsungcloudplatformv2_scr_tag.tag.image_id
}

output "registry_id" {
  value = data.samsungcloudplatformv2_scr_tag.tag.registry_id
}

output "repository_id" {
  value = data.samsungcloudplatformv2_scr_tag.tag.repository_id
}

output "hash_digest" {
  value = data.samsungcloudplatformv2_scr_tag.tag.hash_digest
}

output "manifest" {
  value       = data.samsungcloudplatformv2_scr_tag.tag.manifest
  sensitive   = true
}

output "manifest_media_type" {
  value = data.samsungcloudplatformv2_scr_tag.tag.manifest_media_type
}

output "state" {
  value = data.samsungcloudplatformv2_scr_tag.tag.state
}

output "reference_tags" {
  value = data.samsungcloudplatformv2_scr_tag.tag.reference_tags
}

output "lock_policy" {
  value = data.samsungcloudplatformv2_scr_tag.tag.lock_policy
}

output "created_at" {
  value = data.samsungcloudplatformv2_scr_tag.tag.created_at
}

output "modified_at" {
  value = data.samsungcloudplatformv2_scr_tag.tag.modified_at
}
