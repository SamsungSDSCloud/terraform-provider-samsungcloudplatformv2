---
page_title: "samsungcloudplatformv2_iam_access_keys Data Source - samsungcloudplatformv2"
subcategory: Access Key
description: |-
  list of access key.
---

# samsungcloudplatformv2_iam_access_keys (Data Source)

list of access key.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
  limit = var.access_key_limit
}


output "access_keys" {
  value = data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys
}

variable "access_key_limit" {
  type    = number
  default = 0
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `account_id` (String) Account Id
- `limit` (Number) Limit (between 1 and 10000)
- `marker` (String) Marker (between 1 and 64 characters)
- `sort` (String) Sort

### Read-Only

- `access_keys` (Attributes List) A list of access key. (see [below for nested schema](#nestedatt--access_keys))

<a id="nestedatt--access_keys"></a>
### Nested Schema for `access_keys`

Read-Only:

- `access_key` (String) AccessKey
- `access_key_type` (String) AccessKeyType
- `account_id` (String) AccountId
- `created_at` (String) CreatedAt
- `created_by` (String) CreatedBy
- `description` (String) Description
- `expiration_timestamp` (String) ExpirationTimestamp
- `id` (String) Id
- `modified_at` (String) ModifiedAt
- `modified_by` (String) ModifiedBy
- `parent_access_key_id` (String) ParentAccessKeyId
- `secret_key` (String) SecretKey