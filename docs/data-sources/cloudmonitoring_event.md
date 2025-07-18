---
page_title: "samsungcloudplatformv2_cloudmonitoring_event Data Source - samsungcloudplatformv2"
subcategory: Cloud Monitoring
description: |-
  Event Detail.
---

# samsungcloudplatformv2_cloudmonitoring_event (Data Source)

Event Detail.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_event" "eventDetail" {
  event_id = var.EventId
  x_resource_type = var.XResourceType

}


output "eventDetail" {
  value = data.samsungcloudplatformv2_cloudmonitoring_event.eventDetail
}

variable "EventId" {
  type    = string
  default = ""
}

variable "XResourceType" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `event_id` (String) EventId
- `event_level` (String) EventLevel
- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `x_resource_type` (String) xResourceType

### Read-Only

- `duration_second` (Number) DurationSecond
- `end_dt` (String) EndDt
- `event_message` (String) EventMessage
- `event_policy_id` (Number) EventPolicyId
- `event_state` (String) EventState
- `metric_key` (String) MetricKey
- `metric_name` (String) MetricName
- `object_display_name` (String) ObjectDisplayName
- `object_name` (String) ObjectName
- `object_type` (String) ObjectType
- `object_type_name` (String) ObjectTypeName
- `product_ip_address` (String) ProductIpAddress
- `product_name` (String) ProductName
- `product_resource_id` (String) ProductResourceId
- `product_type_code` (String) ProductTypeCode
- `start_dt` (String) StartDt

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)