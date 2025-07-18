---
page_title: "samsungcloudplatformv2_virtualserver_image Resource - samsungcloudplatformv2"
subcategory: Image
description: |-
  image
---

# samsungcloudplatformv2_virtualserver_image (Resource)

image

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

// Image From URL
resource "samsungcloudplatformv2_virtualserver_image" "image" {
  name              = var.name
  os_distro         = var.os_distro
  disk_format       = var.disk_format
  container_format  = var.container_format
  min_disk          = var.min_disk
  min_ram           = var.min_ram
  visibility        = var.visibility
  protected         = var.protected
  url               = var.url
  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}

// Image From Server (Create)
resource "samsungcloudplatformv2_virtualserver_image" "image2" {
  name              = var.name
  instance_id       = var.instance_id
  tags = {}
}



output "image_output" {
  value = samsungcloudplatformv2_virtualserver_image.image
}

output "image_output2" {
  value = samsungcloudplatformv2_virtualserver_image.image2
}

variable "name" {
  type    = string
  default = ""
}

variable "os_distro" {
  type    = string
  default = ""
}

variable "disk_format" {
  type    = string
  default = ""
}

variable "container_format" {
  type    = string
  default = ""
}

variable "min_disk" {
  type    = number
  default = 0
}

variable "min_ram" {
  type    = number
  default = 0
}

variable "visibility" {
  type    = string
  default = ""
}

variable "protected" {
  type    = bool
  default = false
}

variable "url" {
  type    = string
  default = ""
}
variable "instance_id" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name

### Optional

- `container_format` (String) Container format
- `disk_format` (String) Disk format
- `instance_id` (String) Instance Id
- `min_disk` (Number) Min disk
- `min_ram` (Number) Min ram
- `os_distro` (String) OS Distro
- `protected` (Boolean) Protected
- `tags` (Map of String) A map of key-value pairs representing tags for the resource.
  - Keys must be a maximum of 128 characters.
  - Values must be a maximum of 256 characters.
- `url` (String) Url
- `visibility` (String) Visibility

### Read-Only

- `checksum` (String) Checksum
- `created_at` (String) Created at
- `file` (String) File
- `id` (String) Identifier of the resource.
- `os_hash_algo` (String) OS Hash algo
- `os_hash_value` (String) OS Hash value
- `os_hidden` (Boolean) OS hidden
- `owner` (String) Owner
- `owner_account_name` (String) Owner account name
- `owner_user_name` (String) Owner user name
- `root_device_name` (String) Root device name
- `scp_image_type` (String) SCP Image type
- `scp_k8s_version` (String) SCP K8s version
- `scp_original_image_type` (String) SCP original Image type
- `scp_os_version` (String) SCP OS version
- `size` (Number) Size
- `status` (String) Status
- `updated_at` (String) Updated at
- `virtual_size` (Number) Virtual size
- `volumes` (String) Volumes