---
page_title: "samsungcloudplatformv2_security_group_security_groups Data Source - samsungcloudplatformv2"
subcategory: Security Group
description: |-
  List of security group
---

# samsungcloudplatformv2_security_group_security_groups (Data Source)

List of security group

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_security_group_security_groups" "ids" {
  size = var.size
  page = var.page
  sort = var.sort
  id = var.id
  name = var.name
}


output "ids" {
  value = data.samsungcloudplatformv2_security_group_security_groups.ids.ids
}


variable "size" {
  type    = number
  default = 0
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = ""
}

variable "id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) Id
- `name` (String) Name
- `page` (Number) Page
- `size` (Number) Size
- `sort` (String) Sort

### Read-Only

- `ids` (List of String) Security group Id List