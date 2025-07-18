---
page_title: "samsungcloudplatformv2_backup_backup Data Source - samsungcloudplatformv2"
subcategory: Backup
description: |-
  Show Backup.
---

# samsungcloudplatformv2_backup_backup (Data Source)

Show Backup.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_backup_backup" "backup" {
  region = var.region

  id = var.id
  server_name = var.server_name
  name = var.name

  filter {
    name = var.backup_filter_name
    values = var.backup_filter_values
    use_regex = var.backup_filter_use_regex
  }
}


output "backup" {
  value = data.samsungcloudplatformv2_backup_backup.backup
}

variable "region" {
  type    = string
  default = ""
}

variable "id" {
  type    = string
  default = ""
}

variable "server_name" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "backup_filter_name" {
  type    = string
  default = ""
}

variable "backup_filter_values" {
  type    = list(string)
  default = [""]
}

variable "backup_filter_use_regex" {
  type    = bool
  default = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `id` (String) ID
- `name` (String) Backup name
- `region` (String) Region
- `server_name` (String) Backup server name

### Read-Only

- `backup` (Attributes) A detail of Backup. (see [below for nested schema](#nestedatt--backup))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)


<a id="nestedatt--backup"></a>
### Nested Schema for `backup`

Read-Only:

- `created_at` (String) Created At
- `created_by` (String) Created By
- `encrypt_enabled` (Boolean) Whether to use Encryption
- `id` (String) ID
- `modified_at` (String) Modified At
- `modified_by` (String) Modified By
- `name` (String) Backup server name
- `policy_category` (String) Backup policy category
- `policy_type` (String) Backup policy type
- `retention_period` (String) Backup retention period
- `role_type` (String) Backup role type
- `server_category` (String) Backup server category
- `server_name` (String) Backup server name
- `server_uuid` (String) Backup server UUID
- `state` (String) Backup state