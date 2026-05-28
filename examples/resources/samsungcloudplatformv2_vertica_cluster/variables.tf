variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID"
}

variable "init_config_option" {
  type = object({
    database_locale        = string
    database_name          = string
    database_user_name     = string
    database_user_password = string
    backup_option = object({
      retention_period_day = string
      starting_time_hour   = string
    })
  })
  default = {
    backup_option = {
      retention_period_day = null
      starting_time_hour   = null
    }
    database_locale        = "ko_KR.utf8"
    database_name          = "terra"
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
      role_type = "DATA"
    }]
    role_type        = "DATA"
    server_type_name = "db1v1m2"
  }]
}

variable "instance_name_prefix" {
  type    = string
  default = "terra"
}

variable "name" {
  type    = string
  default = "terraa"
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
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
    use_maintenance_option = false
  }
}

variable "service_state" {
  type    = string
  default = "RUNNING"
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "timezone" {
  type    = string
  default = "Asia/Seoul"
}

variable "license" {
  type    = string
  default = "vertica_community"
}

variable "tags" {
  type = map(string)
  default = {
    key = "value"
  }
}


