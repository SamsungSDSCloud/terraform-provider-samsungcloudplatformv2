---
page_title: "samsungcloudplatformv2_security_group_security_group_rules Data Source - samsungcloudplatformv2"
subcategory: Security Group Rule
description: |-
  List of security group rule
---

# samsungcloudplatformv2_security_group_security_group_rules (Data Source)

List of security group rule

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_security_group_rules" "ids" {
  security_group_id = var.security_group_id
  size = var.size
  page = var.page
  id = var.id
  remote_ip_prefix = var.remote_ip_prefix
  remote_group_id = var.remote_group_id
  description = var.description
  direction = var.direction
  service = var.service
}


output "ids" {
  value = data.samsungcloudplatformv2_security_group_security_group_rules.ids.ids
}

variable "security_group_id" {
  type    = string
  default = ""
}

variable "size" {
  type    = number
  default = 0
}

variable "page" {
  type    = number
  default = 0
}

variable "id" {
  type    = string
  default = ""
}

variable "remote_ip_prefix" {
  type    = string
  default = ""
}

variable "remote_group_id" {
  type    = string
  default = ""
}

variable "description" {
  type    = string
  default = ""
}

variable "direction" {
  type    = string
  default = ""
}

variable "service" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `security_group_id` (String) SecurityGroupId

### Optional

- `description` (String) Description
- `direction` (String) Direction
- `id` (String) Id
- `page` (Number) Page
- `remote_group_id` (String) RemoteGroupId
- `remote_ip_prefix` (String) RemoteIpPrefix
- `service` (String) Service
- `size` (Number) Size
- `sort` (String) Sort

### Read-Only

- `ids` (List of String) Security group Id List