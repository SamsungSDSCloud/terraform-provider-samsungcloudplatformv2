---
page_title: "samsungcloudplatformv2_sqlserver_cluster Resource - samsungcloudplatformv2"
subcategory: Sqlserver
description: |-
  sqlserver
---

# samsungcloudplatformv2_sqlserver_cluster (Resource)

sqlserver

## Example Usage

```terraform
provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_sqlserver_cluster" "cluster" {
  allowable_ip_addresses  = var.allowable_ip_addresses
  dbaas_engine_version_id = var.dbaas_engine_version_id
  nat_enabled             = var.nat_enabled
  ha_enabled              = var.ha_enabled
  init_config_option      = var.init_config_option
  instance_groups         = var.instance_groups
  instance_name_prefix    = var.instance_name_prefix
  name                    = var.name
  subnet_id               = var.subnet_id
  tags                    = var.tags
  service_state           = var.service_state
  timezone                = var.timezone
  maintenance_option      = var.maintenance_option
  vip_public_ip_id        = var.vip_public_ip_id
  virtual_ip_address      = var.virtual_ip_address
  service_watch_log_collection = var.service_watch_log_collection
}


output "cluster_output" {
  value = samsungcloudplatformv2_sqlserver_cluster.cluster
}

variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.1/32", "192.168.10.2/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID"
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "ha_enabled" {
  type    = bool
  default = false
}

variable "init_config_option" {
  type = object({
    ad_config = object({
      ad_dns_servers        = set(string)
      ad_domain_name        = string
      ad_netbios_name       = string
      ad_user_id            = string
      ad_user_password      = string
      failover_cluster_name = string
    })
    ad_enabled             = bool
    audit_enabled          = bool
    database_collation     = string
    database_port          = number
    database_service_name  = string
    database_user_name     = string
    database_user_password = string
    license                = string
    backup_option = object({
      retention_period_day     = string
      starting_time_hour       = string
      archive_frequency_minute = string
      full_backup_day_of_week  = string
    })
    databases = list(object({
      database_name = string
      drive_letter  = string
    }))
  })
  default = {
    ad_config = {
      ad_dns_servers        = ["192.168.35.218"]
      ad_domain_name        = "scp.dev2"
      ad_netbios_name       = "SCP"
      ad_user_id            = "ENTER YOUR RESOURCE'S AD_USER_ID"
      ad_user_password      = "ENTER YOUR RESOURCE'S AD_USER_PASSWORD"
      failover_cluster_name = "Chlwjrghk001"
    }
    ad_enabled    = true
    audit_enabled = false
    backup_option = {
      archive_frequency_minute = null
      full_backup_day_of_week  = null
      retention_period_day     = null
      starting_time_hour       = null
    }
    database_collation     = "SQL_Latin1_General_CP1_CI_AS"
    database_port          = 2866
    database_service_name  = "Sqlserver"
    database_user_name     = "sdsv"
    database_user_password = "ENTER YOUR RESOURCE'S DATABASE_USER_PASSWORD"
    databases = [{
      database_name = "sqlserver"
      drive_letter  = "E"
    }]
    license = "HMWJ3-KY3J2-NMVD7-KG4JR-X2G8G"
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
      role_type = "ACTIVE"
    }]
    role_type        = "ACTIVE"
    server_type_name = "db1v2m8"
  }]
}

variable "instance_name_prefix" {
  type    = string
  default = "sqlserverb"
}

variable "name" {
  type    = string
  default = "sqlserverTd"
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
    period_hour            = null
    starting_day_of_week   = null
    starting_time          = null
    use_maintenance_option = null
  }
}

variable "vip_public_ip_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VIP_PUBLIC_IP_ID"
}

variable "virtual_ip_address" {
  type    = string
  default = null
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

variable "service_watch_log_collection" {
  type    = bool
  default = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allowable_ip_addresses` (Set of String) Allowed IP addresses list  
  - example: ['192.168.10.1/32']
- `dbaas_engine_version_id` (String) DBaaS engine version ID 
  - example: YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID
- `ha_enabled` (Boolean) HA availability 
  - example: False
- `init_config_option` (Attributes) Init config option (see [below for nested schema](#nestedatt--init_config_option))
- `instance_groups` (Attributes List) Instance groups (see [below for nested schema](#nestedatt--instance_groups))
- `instance_name_prefix` (String) Instance name prefix 
  - example: 'test'  
  - minLength: 3  
  - maxLength: 13  
  - pattern: ^[a-z][a-zA-Z0-9\-]*$
- `maintenance_option` (Attributes) Maintenance option (see [below for nested schema](#nestedatt--maintenance_option))
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
  - example: YOUR RESOURCE'S SUBNET_ID
- `timezone` (String) Timezone 
  - example: 'Asia/Seoul'

### Optional

- `service_watch_log_collection` (Boolean) ServiceWatchLogCollection
 - example: false
- `tags` (Map of String) A map of key-value pairs representing tags for the resource.
  - Keys must be a maximum of 128 characters.
  - Values must be a maximum of 256 characters.
- `vip_public_ip_id` (String) VIP Public IP ID (Required when NatEnabled=True & HaEnabled=True)
  - example: YOUR RESOURCE'S VIP_PUBLIC_IP_ID
- `virtual_ip_address` (String) Virtual IP address
  - example: 192.168.4.30

### Read-Only

- `id` (String) Identifier of the resource.
  - example: YOUR RESOURCE'S ID
- `origin_cluster_id` (String) Origin Cluster Id.
  - example: YOUR RESOURCE'S ORIGIN_CLUSTER_ID

<a id="nestedatt--init_config_option"></a>
### Nested Schema for `init_config_option`

Required:

- `ad_config` (Attributes) AdConfig (see [below for nested schema](#nestedatt--init_config_option--ad_config))
- `ad_enabled` (Boolean) AdEnabled
  - example: false
- `audit_enabled` (Boolean) Audit Log Setting
  - example: true
- `backup_option` (Attributes) Backup option (see [below for nested schema](#nestedatt--init_config_option--backup_option))
- `database_collation` (String) Database collation 
  - allowed values: 'SQL_Latin1_General_CP1_CI_AS','Korean_Wansung_CS_AS','Chinese_PRC_CI_AS'
- `database_port` (Number) Database service port 
  - example: 2866
- `database_service_name` (String) Database service name 
  - example: 'Test' 
  - minLength: 1  
  - maxLength: 15  
  - pattern: ^[A-Z][a-zA-Z]*$
- `database_user_name` (String) Database user name 
  - example: 'test' 
  - minLength: 2  
  - maxLength: 20  
  - pattern: ^[a-z]*$
- `database_user_password` (String) Database user password 
  - minLength: 8  
  - maxLength: 30  
  - pattern: ^(?=.*[a-zA-Z])(?=.*[`\-[\]~!@#$%^&*()_+={};:,<.>/?])(?=.*[0-9])(?=\S*[^\w\s]).{8,30} ("'exclude)
- `databases` (Attributes List) Databases (see [below for nested schema](#nestedatt--init_config_option--databases))
- `license` (String) License
  - example: license

Read-Only:

- `origin_region` (String) Origin Region
 -example: kr-west1

<a id="nestedatt--init_config_option--ad_config"></a>
### Nested Schema for `init_config_option.ad_config`

Optional:

- `ad_dns_servers` (Set of String) AD DNS Servers
  - example: 192.168.10.10
- `ad_domain_name` (String) AD Domain Name
  - example: test
- `ad_netbios_name` (String) AD NetBIOS Name
  - example: test
- `ad_user_id` (String) AD User ID
  - example: YOUR RESOURCE'S AD_USER_ID
- `ad_user_password` (String) AD User Password 
  - example: YOUR RESOURCE'S AD_USER_PASSWORD
- `failover_cluster_name` (String) Failover Cluster Name
  - example: testcluster


<a id="nestedatt--init_config_option--backup_option"></a>
### Nested Schema for `init_config_option.backup_option`

Optional:

- `archive_frequency_minute` (String) Backup starting time (minute) 
  - example: 60 
  - pattern: 5 / 10 / 30 / 60
- `full_backup_day_of_week` (String) Days of the week when full backups are performed 
  - allowed values: 'MON','TUE','WED','THU','FRI','SAT','SUN' 
  - default : 'SUN'
- `retention_period_day` (String) Backup retention period (day) 
  - example: 7 
  - min: 7 
  - max: 35
- `starting_time_hour` (String) Backup starting time (hour) 
  - example: 12 
  - min: 00 
  - max: 23


<a id="nestedatt--init_config_option--databases"></a>
### Nested Schema for `init_config_option.databases`

Required:

- `database_name` (String) Database name 
  - example: 'test' 
  - minLength: 3  
  - maxLength: 20  
  - pattern: ^[a-zA-Z][a-zA-Z0-9]*$
- `drive_letter` (String) Drive Letter
  - example: C



<a id="nestedatt--instance_groups"></a>
### Nested Schema for `instance_groups`

Required:

- `block_storage_groups` (Attributes List) BlockStorage groups (see [below for nested schema](#nestedatt--instance_groups--block_storage_groups))
- `instances` (Attributes List) Instances (see [below for nested schema](#nestedatt--instance_groups--instances))
- `role_type` (String) Role type 
  - example: 'ACTIVE' 
  - pattern: ACTIVE (HaEnabled=False) / PRIMARY_SECONDARY (HaEnabled=True)
- `server_type_name` (String) Server type name 
  - example: 'db1v2m4'

Read-Only:

- `id` (String) Instance group ID.
  - example: YOUR RESOURCE'S ID

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

- `id` (String) Block storage group ID
  - example: YOUR RESOURCE'S ID
- `name` (String) Block storage group name
  - example: cluster-Disk-00


<a id="nestedatt--instance_groups--instances"></a>
### Nested Schema for `instance_groups.instances`

Required:

- `role_type` (String) Role type 
  - example: 'PRIMARY' 
  - pattern: ACTIVE / PRIMARY / SECONDARY

Optional:

- `public_ip_id` (String) Public IP ID (Required when NatEnabled=True & HaEnabled=False)
  - example: YOUR RESOURCE'S PUBLIC_IP_ID
- `service_ip_address` (String) User subnet IP address
  - example: 192.168.4.22

Read-Only:

- `name` (String) Instance name
  - example: test001



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