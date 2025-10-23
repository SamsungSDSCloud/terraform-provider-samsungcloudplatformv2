variable "allowable_ip_addresses" {
  type    = list(string)
  default = [""]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = ""
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
    audit_enabled = false
    backup_option = {
      archive_frequency_minute = ""
      full_backup_day_of_week  = ""
      retention_period_day     = ""
      starting_time_hour       = ""
    }
    database_collation     = ""
    database_port          = 0
    database_service_name  = ""
    database_user_name     = ""
    database_user_password = ""
    databases = [{
      database_name = ""
      drive_letter  = ""
    }]
    license = ""
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
      role_type   = ""
      size_gb     = 0
      volume_type = ""
    }]
    instances = [{
      role_type = ""
    }]
    role_type        = ""
    server_type_name = ""
  }]
}

variable "instance_name_prefix" {
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

variable "timezone" {
  type    = string
  default = ""
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
    period_hour            = ""
    starting_day_of_week   = ""
    starting_time          = ""
    use_maintenance_option = false
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
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}

