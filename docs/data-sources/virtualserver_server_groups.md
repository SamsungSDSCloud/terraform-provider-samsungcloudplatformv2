---
page_title: "samsungcloudplatformv2_virtualserver_server_groups Data Source - samsungcloudplatformv2"
subcategory: Server Group
description: |-
  list of server groups.
---

# samsungcloudplatformv2_virtualserver_server_groups (Data Source)

list of server groups.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_server_groups" "ids" {
  filter {
    name = var.server_groups_filter_name
    values = var.server_groups_filter_values
    use_regex = var.server_groups_filter_use_regex
  }
}

output "servers" {
  value = data.samsungcloudplatformv2_virtualserver_server_groups.ids
}

variable "server_groups_filter_name" {
  type    = string
  default = ""
}

variable "server_groups_filter_values" {
  type    = list(string)
  default = [""]
}

variable "server_groups_filter_use_regex" {
  type    = bool
  default = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `ids` (List of String) Server Group ID List

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)