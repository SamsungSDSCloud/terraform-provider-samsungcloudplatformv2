variable "allowable_ip_addresses" {
  type = list(string)
  #  default = []
  default = [""]
}

variable "dbaas_engine_version_id" {
  type = string
  #default = "e4dc55b4607f4a42be1baae17af4ccaa" // Valkey
  default = ""
}

variable "ha_enabled" {
  type = bool
  #  default = true
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
      retention_period_day = ""
      starting_time_hour   = ""
    }
    database_port          = 0
    database_user_password = ""
    sentinel_port          = 0
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
      #      service_ip_address    = string
      #      public_ip_id          = string
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

variable "name" {
  type    = string
  default = ""
}

variable "nat_enabled" {
  type = bool
  #  default = true
  default = false
}

variable "replica_count" {
  type    = number
  default = 0
}

variable "subnet_id" {
  type = string
  #default = "ce93a65b18164072a856c1c31adc1108"
  #default = "eb838b2fb1a5405eab29b1341ca912bb"
  default = ""
}

variable "timezone" {
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

