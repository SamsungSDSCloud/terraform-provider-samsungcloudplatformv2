---
page_title: "samsungcloudplatformv2_cloudmonitoring_producttypes Data Source - samsungcloudplatformv2"
subcategory: samsungcloudplatformv2_cloudmonitoring_producttypes
description: |-
  The Schema of cloudMonitoringProductTypeDataSources.
---

# samsungcloudplatformv2_cloudmonitoring_producttypes (Data Source)

The Schema of cloudMonitoringProductTypeDataSources.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_producttypes" "producttypes" {
  product_category_code = var.ProductCategoryCode
}


output "producttypes" {
  value = data.samsungcloudplatformv2_cloudmonitoring_producttypes.producttypes
}

variable "ProductCategoryCode" {
  type    = string
  default = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `product_category_code` (String) ProductCategoryCode

### Read-Only

- `product_types` (Attributes List) ProductTypes (see [below for nested schema](#nestedatt--product_types))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)


<a id="nestedatt--product_types"></a>
### Nested Schema for `product_types`

Required:

- `parent_product_type_name` (String) ParentProductTypeName
- `product_type_code` (String) ProductTypeCode
- `product_type_name` (String) ProductTypeName
- `state_metric_key` (String) StateMetricKey
