variable "dbaas_engine_version_id" {
  type = string
  default = "1cd2c28ba72447daaaf7e4d7e8dd720b"
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
    audit_enabled = bool
    database_name = string
    database_user_name = string
    database_user_password = string
    database_port = number
    database_character_set = string
    backup_option = object({
      retention_period_day = string
      starting_time_hour = string
      archive_frequency_minute = string
    })
  })
  default = {
    audit_enabled = true
    database_name = "dbname"
    database_user_name = "username"
    database_user_password = "dbpwd00!"
    database_port = 2866
    database_character_set = "utf8"
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
  default = "instacename"
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
    period_hour = null
    starting_day_of_week = null
    starting_time = null
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
