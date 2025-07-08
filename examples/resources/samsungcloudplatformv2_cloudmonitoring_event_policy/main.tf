provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_cloudmonitoring_event_policy" "event_policy" {

    x_resource_type = var.XResourceType
    product_resource_id = var.ProductResourceId
    event_level = var.EventLevel
    disable_yn = var.DisableYn
    event_message_prefix = var.EventMessagePrefix
    ft_count = var.FtCount
    is_log_metric = var.IsLogMetric
    metric_key =  var.MetricKey
    object_name = var.ObjectName
#    object_display_name = var.ObjectDisplayName
    event_occur_time_zone = var.EventOccurTimeZone
    object_type = var.ObjectType
    object_type_name = var.ObjectTypeName
#    event_threshold = var.EventThreshold
#    event_policy_id = var.EventPolicyId
#    metric_key              = "libvirt.domain.cpu.scpm.usage"
    event_policy   = var.EventPolicy
#    notification_recipient = []
}