variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID"
}

variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.1/32", "1.1.1.1/32"]
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
    database_name          = string
    database_user_name     = string
    database_user_password = string
    database_port          = number
    database_character_set = string
    backup_option = object({
      retention_period_day     = string
      starting_time_hour       = string
      archive_frequency_minute = string
    })
  })
  default = {
    audit_enabled = true
    backup_option = {
      archive_frequency_minute = null
      retention_period_day     = null
      starting_time_hour       = null
    }
    database_character_set = "utf8"
    database_name          = "terra"
    database_port          = 2866
    database_user_name     = "terra"
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
      role_type = "ACTIVE"
    }]
    role_type        = "ACTIVE"
    server_type_name = "db1v2m4"
  }]
}

variable "instance_name_prefix" {
  type    = string
  default = "terrac"
}

variable "name" {
  type    = string
  default = "terrac"
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
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
    period_hour            = null
    starting_day_of_week   = null
    starting_time          = null
    use_maintenance_option = false
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



