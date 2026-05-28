---
page_title: "samsungcloudplatformv2_searchengine_cluster Resource - samsungcloudplatformv2"
subcategory: Searchengine
description: |-
  searchengine
---

# samsungcloudplatformv2_searchengine_cluster (Resource)

searchengine

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_searchengine_cluster" "cluster" {
  allowable_ip_addresses  = var.allowable_ip_addresses
  dbaas_engine_version_id = var.dbaas_engine_version_id
  init_config_option      = var.init_config_option
  instance_groups         = var.instance_groups
  instance_name_prefix    = var.instance_name_prefix
  is_combined             = var.is_combined
  nat_enabled             = var.nat_enabled
  license                 = var.license
  name                    = var.name
  subnet_id               = var.subnet_id
  tags                    = var.tags
  service_state           = var.service_state
  timezone                = var.timezone
  maintenance_option      = var.maintenance_option
}


output "cluster_output" {
  value = samsungcloudplatformv2_searchengine_cluster.cluster
}

variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.2/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID"
}

variable "is_combined" {
  type    = bool
  default = true
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "init_config_option" {
  type = object({
    database_port          = number
    database_user_name     = string
    database_user_password = string
    backup_option = object({
      retention_period_day = string
      starting_time_hour   = string
    })
  })
  default = {
    backup_option = {
      retention_period_day = "7"
      starting_time_hour   = "12"
    }
    database_port          = 9201
    database_user_name     = "terraformtest"
    database_user_password = "ENTER YOUR RESOURCE'S DATABASE_USER_PASSWORD"
  }
}

variable "instance_groups" {
  type = list(object({
    role_type        = string
    server_type_name = string
    block_storage_groups = list(object({
      role_type   = string
      volume_type = string
      size_gb     = number
    }))
    instances = list(object({
      role_type = string
    }))
  }))
  default = [{
    block_storage_groups = [{
      role_type   = "OS"
      size_gb     = 104
      volume_type = "SSD"
      }, {
      role_type   = "DATA"
      size_gb     = 16
      volume_type = "SSD"
    }]
    instances = [{
      role_type = "MASTER_DATA"
    }]
    role_type        = "MASTER_DATA"
    server_type_name = "ses1v2m4"
    }, {
    block_storage_groups = [{
      role_type   = "OS"
      size_gb     = 104
      volume_type = "SSD"
    }]
    instances = [{
      role_type = "KIBANA"
    }]
    role_type        = "KIBANA"
    server_type_name = "ses1v2m4"
  }]
}


variable "instance_name_prefix" {
  type    = string
  default = "searchb"
}

variable "name" {
  type    = string
  default = "searchb"
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "timezone" {
  type    = string
  default = "Asia/Seoul"
}

// OPTION
variable "maintenance_option" {
  type = object({
    period_hour            = string
    starting_day_of_week   = string
    starting_time          = string
    use_maintenance_option = bool
  })
  default = {
    period_hour            = "0.5"
    starting_day_of_week   = "MON"
    starting_time          = "0000"
    use_maintenance_option = true
  }
}

variable "service_state" {
  type    = string
  default = "RUNNING"
}

variable "tags" {
  type = map(string)
  default = {
    key = "value"
  }
}

variable "license" {
  type    = string
  default = "{\n   \"license\":{\n      \"uid\":\"f5c002e6-c29d-4dde-bce4-cdb35bf85ab8\",\n      \"type\":\"trial\",\n      \"issue_date_in_millis\":1745193600000,\n      \"expiry_date_in_millis\":1752969599999,\n      \"max_nodes\":3,\n      \"issued_to\":\"S-Core Co., Ltd. (non-production environments)\",\n      \"issuer\":\"API\",\n      \"signature\":\"AAAAAwAAAA2ZzeOwjqKpsuwZQJCbAAABmC9ZN0hjZDBGYnVyRXpCOW5Bb3FjZDAxOWpSbTVoMVZwUzRxVk1PSmkxaktJRVl5MUYvUWh3bHZVUTllbXNPbzBUemtnbWpBbmlWRmRZb25KNFlBR2x0TXc2K2p1Y1VtMG1UQU9TRGZVSGRwaEJGUjE3bXd3LzRqZ05iLzRteWFNekdxRGpIYlFwYkJiNUs0U1hTVlJKNVlXekMrSlVUdFIvV0FNeWdOYnlESDc3MWhlY3hSQmdKSjJ2ZTcvYlBFOHhPQlV3ZHdDQ0tHcG5uOElCaDJ4K1hob29xSG85N0kvTWV3THhlQk9NL01VMFRjNDZpZEVXeUtUMXIyMlIveFpJUkk2WUdveEZaME9XWitGUi9WNTZVQW1FMG1DenhZU0ZmeXlZakVEMjZFT2NvOWxpZGlqVmlHNC8rWVVUYzMwRGVySHpIdURzKzFiRDl4TmM1TUp2VTBOUlJZUlAyV0ZVL2kvVk10L0NsbXNFYVZwT3NSU082dFNNa2prQ0ZsclZ4NTltbU1CVE5lR09Bck93V2J1Y3c9PQAAAQCdy6pBHZq1yaRofO7pDJJrEGmcnsUpa4BmsWywjbOU3zc3zLp6hvVuoWQ8bys6wl+lflToI51p6WinnyGmiFJoSvkNUMquuJMvEza5MlBzu4yb7mZKEUa4hxvz7IjQOltXP7KUXnMa98SrHx9Fkf/N80ZFhGD9t25UBkBgYvKEWTCztmOmNXTX/Sdq+JixfuijiA75EVGWhge7tGc6OpfHBgpZOJIxGOTAIUQiBxWK1ZdBF76CwhEVTbkMiKutvNsbwwo+yWiNGjq0mCUDYZUMXp6T4xk0VLzrmgGhBCqSR2BJzESq1Yk2VmIBv2Sn7JzokbGB3SB7FYQHKArVsinU\",\n      \"start_date_in_millis\":1745193600000\n   }\n}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allowable_ip_addresses` (Set of String) Allowed IP addresses list  
  - example: ['192.168.10.1/32']
- `dbaas_engine_version_id` (String) DBaaS engine version ID 
  - example: YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID
- `init_config_option` (Attributes) Init config option (see [below for nested schema](#nestedatt--init_config_option))
- `instance_groups` (Attributes List) Instance groups (see [below for nested schema](#nestedatt--instance_groups))
- `instance_name_prefix` (String) Instance name prefix 
  - example: 'test'  
  - minLength: 3  
  - maxLength: 13  
  - pattern: ^[a-z][a-zA-Z0-9\-]*$
- `is_combined` (Boolean) MASTER,DATA combined (IsCombined=true), MASTER,DATA seperated (IsCombined=False)
- `maintenance_option` (Attributes) MaintenanceOption (see [below for nested schema](#nestedatt--maintenance_option))
- `name` (String) Cluster name 
  - example: 'test'  
  - minLength: 3  
  - maxLength: 20  
  - pattern: ^[a-zA-Z]*$
- `nat_enabled` (Boolean) NAT availability 
  - example: False
- `service_state` (String) Service state 
  - example : 'RUNNING' (Create,Start) / 'STOPPED' (Stop)
- `subnet_id` (String) Subnet ID
- `timezone` (String) Timezone 
  - example: 'Asia/Seoul'

### Optional

- `license` (String) License
- `tags` (Map of String) A map of key-value pairs representing tags for the resource.
  - Keys must be a maximum of 128 characters.
  - Values must be a maximum of 256 characters.

### Read-Only

- `id` (String) Identifier of the resource.

<a id="nestedatt--init_config_option"></a>
### Nested Schema for `init_config_option`

Required:

- `backup_option` (Attributes) BackupOption (see [below for nested schema](#nestedatt--init_config_option--backup_option))
- `database_port` (Number) Database service port 
  - example: 9201
- `database_user_name` (String) Database user name 
  - example: 'test' 
  - minLength: 2  
  - maxLength: 20  
  - pattern: ^[a-z]*$
- `database_user_password` (String) Database user password 
  - minLength: 8  
  - maxLength: 30  
  - pattern: ^(?=.*[a-zA-Z])(?=.*[`\-[\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\S*[^\w\s]).{8,30} ("'제외)

<a id="nestedatt--init_config_option--backup_option"></a>
### Nested Schema for `init_config_option.backup_option`

Optional:

- `retention_period_day` (String) Backup retention period (day) 
  - example: 7 
  - min: 7 
  - max: 35
- `starting_time_hour` (String) Backup starting time (hour) 
  - example: 12 
  - min: 00 
  - max: 23



<a id="nestedatt--instance_groups"></a>
### Nested Schema for `instance_groups`

Required:

- `block_storage_groups` (Attributes List) BlockStorage groups (see [below for nested schema](#nestedatt--instance_groups--block_storage_groups))
- `instances` (Attributes List) Instances (see [below for nested schema](#nestedatt--instance_groups--instances))
- `role_type` (String) Role type 
  - example: 'MASTER_DATA' 
  - pattern: MASTER_DATA (IsCombined=True) / MASTER, DATA (IsCombined=False) / KIBANA, DASHBOARDS (required)
- `server_type_name` (String) Server type name 
  - example: 'se1v2m4'

Read-Only:

- `id` (String) Id

<a id="nestedatt--instance_groups--block_storage_groups"></a>
### Nested Schema for `instance_groups.block_storage_groups`

Required:

- `role_type` (String) Role type 
  - example: 'OS'
- `size_gb` (Number) Size in GB 
  - example: 104 
  - minLength: 16  
  - maxLength: 5120
- `volume_type` (String) Volume type 
  - example: 'SSD'

Read-Only:

- `id` (String) Id
- `name` (String) Name


<a id="nestedatt--instance_groups--instances"></a>
### Nested Schema for `instance_groups.instances`

Required:

- `role_type` (String) Role type 
  - example: 'MASTER_DATA' 
  - pattern: MASTER_DATA / MASTER / DATA / KIBANA / DASHBOARDS

Optional:

- `public_ip_id` (String) Public IP ID (Required when NatEnabled=True)
- `service_ip_address` (String) User subnet IP address

Read-Only:

- `name` (String) Name



<a id="nestedatt--maintenance_option"></a>
### Nested Schema for `maintenance_option`

Optional:

- `period_hour` (String) Period in hours 
  - example: 1
- `starting_day_of_week` (String) Starting day of week 
  - example: 'MON'
- `starting_time` (String) Starting time 
  - example: '0000'
- `use_maintenance_option` (Boolean) Use maintenance option 
  - example: False