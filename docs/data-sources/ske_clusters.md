---
page_title: "samsungcloudplatformv2_ske_clusters Data Source - samsungcloudplatformv2"
subcategory: Kubernetes Engine Cluster
description: |-
  list of cluster.
---

# samsungcloudplatformv2_ske_clusters (Data Source)

list of cluster.

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_ske_clusters" "ids" {
    page               = var.page
    size               = var.size
    sort               = var.sort
    name               = var.name
    subnet_id          = var.subnet_id
    status             = var.status
    kubernetes_version = var.kubernetes_version
    region = var.clusters_region
}


output "clusters" {
  value = data.samsungcloudplatformv2_ske_clusters.ids
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

variable "name" {
  type    = string
  default = ""
}

variable "subnet_id" {
  type    = string
  default = ""
}

variable "status" {
  type    = list(string)
  default = [""]
}

variable "kubernetes_version" {
  type    = list(string)
  default = [""]
}




variable "clusters_region" {
  type    = string
  default = ""
}

variable "clusters_filter_name" {
  type    = string
  default = ""
}

variable "clusters_filter_values" {
  type    = list(string)
  default = [""]
}

variable "clusters_filter_use_regex" {
  type    = bool
  default = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block List) Filter (see [below for nested schema](#nestedblock--filter))
- `kubernetes_version` (List of String) KubernetesVersion List
- `name` (String) Name
- `page` (Number) Page
- `region` (String) Region
- `size` (Number) Size (between 1 and 10000)
- `sort` (String) Sort
- `status` (List of String) Status List
- `subnet_id` (String) SubnetId (validation)
- `tags` (Map of String) A map of key-value pairs representing tags for the resource.
 - Keys must be a maximum of 128 characters.
 - Values must be a maximum of 256 characters.

### Read-Only

- `ids` (List of String) ID List

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `use_regex` (Boolean) Enable regex match for values
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)