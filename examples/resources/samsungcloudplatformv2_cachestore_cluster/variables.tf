variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.1/32", "192.168.10.2/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID" // Redis
}

variable "ha_enabled" {
  type    = bool
  default = false
}

variable "init_config_option" {
  type = object({
    database_port          = number
    database_user_password = string
    sentinel_port          = number
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
    database_port          = 6378
    database_user_password = "ENTER YOUR RESOURCE'S DATABASE_USER_PASSWORD"
    sentinel_port          = 26378
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
      size_gb     = 56
      volume_type = "SSD"
    }]
    instances = [{
      role_type = "MASTER"
    }]
    role_type        = "MASTER"
    server_type_name = "redis1v1m2"
  }]
}

variable "instance_name_prefix" {
  type    = string
  default = "terra"
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

variable "name" {
  type    = string
  default = "terrad"
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "replica_count" {
  type    = number
  default = 0
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "timezone" {
  type    = string
  default = "Asia/Seoul"
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


