variable "dbaas_engine_version_id" {
  type = string
  default = "03cf36ee3162415e854c411e3690c85e"
}

variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32"]
}

variable "ha_enabled" {
  type = bool
  default = false
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "init_config_option" {
  type = object({
    database_name = string
    database_user_name = string
    database_user_password = string
    database_port = number
    database_character_set = string
    database_case_sensitive = bool
    backup_option = object({
      retention_period_day = string
      starting_time_hour = string
      archive_frequency_minute = string
    })
  })
  default = {
    database_name = "dbname"
    database_user_name = "username"
    database_user_password = "Paswrd0!!"
    database_port = 2866
    database_character_set = "utf8mb3"
    database_case_sensitive = true
    backup_option = {
      retention_period_day = "7"
      starting_time_hour = "12"
      archive_frequency_minute = "60"
    }
  }
}

variable "instance_groups" {
  type = list(object({
    role_type = string
    server_type_name = string
    block_storage_groups = list(object({
      role_type = string
      volume_type = string
      size_gb = number
    }))
    instances = list(object({
      role_type = string
      public_ip_id = string
    }))
  }))
  default = [{
    role_type = "ACTIVE"
    server_type_name = "db1v2m4"
    block_storage_groups = [
      {
        "role_type": "OS",
        "volume_type": "SSD",
        "size_gb": 104
      },
      {
        "role_type": "DATA",
        "volume_type": "SSD",
        "size_gb": 16
      }
    ]
    instances = [
      {
        "role_type": "ACTIVE",
        "public_ip_id" : "8a463aa4b1dc4f279c3f53b94dc45e74"
      }
    ]
  }]
}

variable "instance_name_prefix" {
  type = string
  default = "instname"
}

variable "name" {
  type = string
  default = "name"
}

variable "subnet_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "timezone" {
  type = string
  default = "Asia/Seoul"
}

variable "maintenance_option" {
  type = object({
    period_hour = string
    starting_day_of_week = string
    starting_time = string
    use_maintenance_option = bool
  })
  default = {
    period_hour = "1"
    starting_day_of_week = "MON"
    starting_time = "0100"
    use_maintenance_option = false
  }
}

variable "vip_public_ip_id" {
  type    = string
  default = ""
}

variable "virtual_ip_address" {
  type = string
  default = ""
}

variable "service_state" {
  type = string
  default = "RUNNING"
}

variable "tags" {
  type = map(string)
  default = {
    "key" : "value"
  }
}
