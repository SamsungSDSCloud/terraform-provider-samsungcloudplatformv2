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


