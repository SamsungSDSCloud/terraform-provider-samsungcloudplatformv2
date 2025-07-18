---
page_title: "samsungcloudplatformv2_virtualserver_server Data Source - samsungcloudplatformv2"
subcategory: Virtual Server
description: |-
  list of servers.
---

# samsungcloudplatformv2_virtualserver_server (Data Source)

list of servers.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_server" "server" {
  id = var.id

  name = var.name
  ip = var.ip
  state = var.state
  product_category = var.product_category
  vpc_id = var.vpc_id
  server_type_id = var.server_type_id
  auto_scaling_group_id = var.auto_scaling_group_id

  filter {
    name = var.server_filter_name
    values = var.server_filter_values
    use_regex = var.server_filter_use_regex
  }
}


output "image" {
  value = data.samsungcloudplatformv2_virtualserver_server.server
}

variable "id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "ip" {
  type    = string
  default = ""
}

variable "state" {
  type    = string
  default = ""
}

variable "product_category" {
  type    = string
  default = ""
}

variable "vpc_id" {
  type    = string
  default = ""
}

variable "server_type_id" {
  type    = string
  default = ""
}

variable "auto_scaling_group_id" {
  type    = string
  default = ""
}

variable "server_filter_name" {
  type    = string
  default = ""
}

variable "server_filter_values" {
  type    = list(string)
  default = [""]
}

variable "server_filter_use_regex" {
  type    = bool
  default = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `auto_scaling_group_id` (String) Auto scaling group ID
- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `id` (String) ID
- `ip` (String) Ip
- `name` (String) Name
- `product_category` (String) Product category
- `product_offering` (String) Product offering
- `server_type_id` (String) Server type ID
- `state` (String) State
- `vpc_id` (String) VPC ID

### Read-Only

- `server` (Attributes) Server. (see [below for nested schema](#nestedatt--server))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)


<a id="nestedatt--server"></a>
### Nested Schema for `server`

Read-Only:

- `account_id` (String) Account ID
- `addresses` (Attributes List) Addresses (see [below for nested schema](#nestedatt--server--addresses))
- `auto_scaling_group_id` (String) Auto scaling group ID
- `created_at` (String) Created at
- `created_by` (String) Created by
- `disk_config` (String) Disk config
- `id` (String) ID
- `image_id` (String) Image ID
- `keypair_name` (String) Keypair name
- `launch_configuration_id` (String) Launch Configuration ID
- `locked` (Boolean) Locked
- `metadata` (Map of String) Metadata
- `modified_at` (String) Modified at
- `name` (String) Name
- `planned_compute_os_type` (String) Planned compute os type
- `product_category` (String) Product category
- `product_offering` (String) Product offering
- `security_groups` (Attributes List) Security groups (see [below for nested schema](#nestedatt--server--security_groups))
- `server_group_id` (String) Server group ID
- `server_type` (Attributes) Server type (see [below for nested schema](#nestedatt--server--server_type))
- `state` (String) State
- `volumes` (Attributes List) Volumes (see [below for nested schema](#nestedatt--server--volumes))
- `vpc_id` (String) Vpc ID

<a id="nestedatt--server--addresses"></a>
### Nested Schema for `server.addresses`

Read-Only:

- `ip_addresses` (Attributes List) IP addresses (see [below for nested schema](#nestedatt--server--addresses--ip_addresses))
- `subnet_name` (String) Subnet name

<a id="nestedatt--server--addresses--ip_addresses"></a>
### Nested Schema for `server.addresses.ip_addresses`

Read-Only:

- `ip_address` (String) IP address
- `version` (Number) Version



<a id="nestedatt--server--security_groups"></a>
### Nested Schema for `server.security_groups`

Read-Only:

- `name` (String) Name


<a id="nestedatt--server--server_type"></a>
### Nested Schema for `server.server_type`

Read-Only:

- `disk` (Number) Disk
- `ephemeral` (Number) Ephemeral
- `extra_specs` (Map of String) Extra specs
- `id` (String) ID
- `name` (String) Name
- `ram` (Number) Ram
- `swap` (Number) Swap
- `vcpus` (Number) Vcpus


<a id="nestedatt--server--volumes"></a>
### Nested Schema for `server.volumes`

Read-Only:

- `delete_on_termination` (Boolean) Delete on termination
- `id` (String) ID