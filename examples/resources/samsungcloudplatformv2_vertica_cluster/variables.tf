variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type = string
  default = "09c2fe88089040ffa035604e38f7e4e9"
}

variable "init_config_option" {
  type = object({
    database_locale        = string
    database_name          = string
    database_user_name     = string
    database_user_password = string
    backup_option = object({
      retention_period_day     = string
      starting_time_hour       = string
    })
  })
  default = {
    database_locale        = "ko_KR.utf8"
    database_name          = "dbname"
    database_user_name     = "username"
    database_user_password = "Paswrd0!!"
    backup_option = {
      retention_period_day     = "7"
      starting_time_hour       = "12"
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
      role_type          = string
      service_ip_address = string
      public_ip_id       = string
    }))
  }))
  default = [
    {
      role_type        = "DATA"
      server_type_name = "db1v1m2"
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
          "role_type" : "DATA",
          "service_ip_address" : "192.168.28.217",
          "public_ip_id" : "8a463aa4b1dc4f279c3f53b94dc45e74",
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
    "key" : "value"
  }
}