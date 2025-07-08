variable "dbaas_engine_version_id" {
  type    = string
  default = ""
}

variable "allowable_ip_addresses" {
  type    = list(string)
  default = [""]
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
    audit_enabled = false
    backup_option = {
      archive_frequency_minute = ""
      retention_period_day     = ""
      starting_time_hour       = ""
    }
    database_character_set = ""
    database_name          = ""
    database_port          = 0
    database_user_name     = ""
    database_user_password = ""
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
      #      public_ip_id = string
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
  type = string
  #  default = "ce93a65b18164072a856c1c31adc1108"
  default = ""
}

variable "timezone" {
  type    = string
  default = ""
}

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
  type = string
  #  default = "STOPPED"
  default = ""
}

variable "tags" {
  type    = map(string)
  default = null
}


