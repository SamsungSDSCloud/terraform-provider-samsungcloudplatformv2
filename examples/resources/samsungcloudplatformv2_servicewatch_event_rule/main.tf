provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_servicewatch_event_rule" "event_rule" {
  description = var.description
  event_ids = var.event_ids
  name    = var.event_rule_name
  recipient_ids = var.recipient_ids
  service_id = var.service_id
  resource_type_id = var.resource_type_id
  srn_list = var.srn_list
  tags = var.tags
}
