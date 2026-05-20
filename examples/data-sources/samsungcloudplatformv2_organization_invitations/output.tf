output "invitations" {
  value = {
    organization_invitations = data.samsungcloudplatformv2_organization_invitations.invitations.organization_invitations
    total_count              = data.samsungcloudplatformv2_organization_invitations.invitations.total_count
    page                     = data.samsungcloudplatformv2_organization_invitations.invitations.page
    size                     = data.samsungcloudplatformv2_organization_invitations.invitations.size
    sort                     = data.samsungcloudplatformv2_organization_invitations.invitations.sort_result
  }
}