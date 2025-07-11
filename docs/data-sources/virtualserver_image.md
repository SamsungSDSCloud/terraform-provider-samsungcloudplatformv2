---
page_title: "samsungcloudplatformv2_virtualserver_image Data Source - samsungcloudplatformv2"
subcategory: Image
description: |-
  Image.
---

# samsungcloudplatformv2_virtualserver_image (Data Source)

Image.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_virtualserver_image" "image" {
  id = var.id

  scp_image_type = var.scp_image_type
  scp_original_image_type = var.scp_original_image_type
  name = var.name
  os_distro = var.os_distro
  status = var.status
  visibility = var.visibility

  filter {
    name = var.image_filter_name
    values = var.image_filter_values
    use_regex = var.image_filter_use_regex
  }
}


output "image" {
  value = data.samsungcloudplatformv2_virtualserver_image.image
}

variable "id" {
  type    = string
  default = ""
}

variable "scp_image_type" {
  type    = string
  default = ""
}

variable "scp_original_image_type" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "os_distro" {
  type    = string
  default = ""
}

variable "status" {
  type    = string
  default = ""
}

variable "visibility" {
  type    = string
  default = ""
}

variable "image_filter_name" {
  type    = string
  default = ""
}

variable "image_filter_values" {
  type    = list(string)
  default = [""]
}

variable "image_filter_use_regex" {
  type    = bool
  default = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `id` (String) ID
- `name` (String) Name
- `os_distro` (String) OS Distro
- `scp_image_type` (String) SCP Image type
- `scp_original_image_type` (String) SCP Original Image type
- `status` (String) Status
- `visibility` (String) Visibility

### Read-Only

- `image` (Attributes) Image. (see [below for nested schema](#nestedatt--image))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)


<a id="nestedatt--image"></a>
### Nested Schema for `image`

Read-Only:

- `checksum` (String) Checksum
- `container_format` (String) Container format
- `created_at` (String) Created at
- `disk_format` (String) Disk format
- `file` (String) File
- `id` (String) ID
- `min_disk` (Number) Min disk
- `min_ram` (Number) Min ram
- `name` (String) Name
- `os_distro` (String) OS Distro
- `os_hash_algo` (String) OS Hash algo
- `os_hash_value` (String) OS Hash value
- `os_hidden` (Boolean) OS hidden
- `owner` (String) Owner
- `owner_account_name` (String) Owner account name
- `owner_user_name` (String) Owner user name
- `protected` (Boolean) Protected
- `root_device_name` (String) Root device name
- `scp_image_type` (String) SCP Image type
- `scp_k8s_version` (String) SCP K8s version
- `scp_original_image_type` (String) SCP original Image type
- `scp_os_version` (String) SCP OS version
- `size` (Number) Size
- `status` (String) Status
- `updated_at` (String) Updated at
- `url` (String) Url
- `virtual_size` (Number) Virtual size
- `visibility` (String) Visibility
- `volumes` (String) Volumes