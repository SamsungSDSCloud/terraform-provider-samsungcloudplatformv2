resource "samsungcloudplatformv2_ske_cluster_linked_resources" "example" {
  cluster_id = var.cluster_id

  linked_resources = [
    {
      id   = var.linked_resource_id
      name = var.linked_resource_name
      type = var.linked_resource_type
    }
  ]
}
