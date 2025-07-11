---
page_title: "samsungcloudplatformv2_quota_account_quota Data Source - samsungcloudplatformv2"
subcategory: Account Quota
description: |-
  The account quota
---

# samsungcloudplatformv2_quota_account_quota (Data Source)

The account quota

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_quota_account_quota" "account_quota"{
  id = var.account_quota_account_quota_id
}


output "account_quota" {
  value = data.samsungcloudplatformv2_quota_account_quota.account_quota
}


variable "account_quota_account_quota_id" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Account Quota ID

### Read-Only

- `account_quota` (Attributes) (see [below for nested schema](#nestedatt--account_quota))

<a id="nestedatt--account_quota"></a>
### Nested Schema for `account_quota`

Read-Only:

- `account_id` (String) Unique identifier for the account
- `account_name` (String) Name of the account
- `adjustable` (Boolean) Flag indicating if additional quota is being requested
- `applied_value` (Number)
- `approval` (Boolean) Approval
- `class_value` (String) Value associated with the request class
- `created_at` (String) Created At
- `description` (String) Detailed description of the quota item
- `free_rate` (Number) Free Rate
- `id` (String) Account Quota ID
- `initial_value` (Number) Initial quota value allocated
- `modified_at` (String) Modified At
- `quota_item` (String) Specific quota item within the resource
- `reduction` (Boolean) Reduction
- `request` (Boolean) Reqeust
- `request_class` (String) Classification of the quota request (e.g., Account, Region)
- `resource_type` (String) Type of the resource (e.g., Virtual Server, Storage)
- `service` (String) Name of the service to which quota applies
- `srn` (String) Service Resource Name for the quota item
- `unit` (String) Unit in which the quota value is measured (e.g., EA, GB)