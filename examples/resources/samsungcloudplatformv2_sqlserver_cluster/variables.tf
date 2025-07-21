variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32", "192.168.10.2/32"]
}

variable "dbaas_engine_version_id" {
  type = string
  default = "266fba6314a04d308f34a221c4f9c7a5"
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
    audit_enabled          = bool
    database_collation     = string
    database_port          = number
    database_service_name  = string
    database_user_name     = string
    database_user_password = string
    license = string
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
    audit_enabled          = false
    database_collation     = "SQL_Latin1_General_CP1_CI_AS"
    database_port          = 2866
    database_service_name  = "Sqlserver"
    database_user_name     = "sqlserver"
    database_user_password = "Password001!"
    license = "AAAAA-BBBBB-CCCCC-DDDDD-EEEEE"
    backup_option = {
      retention_period_day     = null
      starting_time_hour       = null
      archive_frequency_minute = null
      full_backup_day_of_week  = null
    }
    databases = [{
      database_name = "sqlserver"
      drive_letter = "E"
    }]
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
  default = [
    {
      role_type        = "ACTIVE"
      server_type_name = "db1v2m8"
      block_storage_groups = [
        {
          "role_type" : "OS",
          "volume_type" : "SSD",
          "size_gb" : 104
        },
        {
          "role_type" : "DATA",
          "volume_type" : "SSD",
          "size_gb" : 16
        }
      ]
      instances = [
        {
          "role_type" : "ACTIVE",
        }
      ]
    }
  ]
}

variable "instance_name_prefix" {
  type    = string
  default = "instname"
}

variable "name" {
  type    = string
  default = "name"
}

variable "subnet_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "timezone" {
  type    = string
  default = "Asia/Seoul"
}

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

variable "vip_public_ip_id" {
  type    = string
  default = ""
}

variable "virtual_ip_address" {
  type    = string
  default = ""
}

variable "service_state" {
  type    = string
  default = "RUNNING"
}

variable "tags" {
  type = map(string)
  default = {
    "key" : "value"
  }
}