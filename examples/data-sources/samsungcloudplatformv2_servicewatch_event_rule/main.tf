provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_servicewatch_event_rule" "event_rule" {
  event_rule_id = var.id
}
