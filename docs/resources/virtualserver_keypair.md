---
page_title: "samsungcloudplatformv2_virtualserver_keypair Resource - samsungcloudplatformv2"
subcategory: Keypair
description: |-
  keypair
---

# samsungcloudplatformv2_virtualserver_keypair (Resource)

keypair

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_virtualserver_keypair" "keypair" {
  name = var.name
  tags = {
    "test_terraform_tag_key": "test_terraform_tag_value"
  }
}



output "keypair_output" {
  value = samsungcloudplatformv2_virtualserver_keypair.keypair
}

variable "name" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name

### Optional

- `tags` (Map of String) A map of key-value pairs representing tags for the resource.
  - Keys must be a maximum of 128 characters.
  - Values must be a maximum of 256 characters.

### Read-Only

- `created_at` (String) Created at
- `fingerprint` (String) Fingerprint
- `id` (Number) Identifier of the resource.
- `private_key` (String) Private key
- `public_key` (String) Public key
- `type` (String) Type
- `user_id` (String) User ID