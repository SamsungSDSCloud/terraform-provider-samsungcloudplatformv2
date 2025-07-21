variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type = string
  default = "821110c054044f0188b0811e53ef9ed6"
}

variable "ha_enabled" {
  type    = bool
  default = false
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "init_config_option" {
  type = object({
    audit_enabled          = bool
    database_encoding      = string
    database_locale        = string
    database_name          = string
    database_port          = number
    database_user_name     = string
    database_user_password = string
    backup_option = object({
      retention_period_day     = string
      starting_time_hour       = string
      archive_frequency_minute = string
    })
  })
  default = {
    audit_enabled          = false
    database_encoding      = "UTF-8"
    database_locale        = "C"
    database_name          = "epaasdbname"
    database_port          = 2866
    database_user_name     = "terraformuser"
    database_user_password = "password001!"
    backup_option = {
      retention_period_day     = "7"
      starting_time_hour       = "12"
      archive_frequency_minute = "60"
    }
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
      role_type    = string
    }))
  }))
  default = [
    {
      role_type        = "ACTIVE"
      server_type_name = "db1v2m4"
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
  default = "terraprefix"
}

variable "name" {
  type    = string
  default = "terraname"
}

variable "subnet_id" {
  type    = string
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
    "key2" : "value2"
  }
}