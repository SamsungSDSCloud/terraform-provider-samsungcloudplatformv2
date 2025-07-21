variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type = string
  default = "aef8e9ace6f54207bdf6266d4028cb74"
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
      retention_period_day     = string
      starting_time_hour       = string
    })
  })
  default = {
    database_port          = 2866
    database_user_password = "Pswrd^00^"
    sentinel_port          = 26378
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
      role_type             = string
      service_ip_address    = string
      public_ip_id          = string
    }))
  }))
  default = [
    {
      role_type        = "MASTER"
      server_type_name = "redis1v1m2"
      block_storage_groups = [
        {
          "role_type" : "OS",
          "volume_type" : "SSD",
          "size_gb" : 104
        },
        {
          "role_type" : "DATA",
          "volume_type" : "SSD",
          "size_gb" : 56
        }
      ]
      instances = [
        {
          "role_type" : "MASTER",
         "service_ip_address" : "192.168.10.1",
         "public_ip_id" : "83c3c73d457345e3829ee6d5557c0011",
        }
      ]
    }
  ]
}

variable "instance_name_prefix" {
  type    = string
  default = "terra"
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

variable "name" {
  type    = string
  default = "terrad"
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "replica_count" {
  type = number
  default = 0
}

variable "subnet_id" {
  type = string
  default = "83c3c73d457345e3829ee6d5557c0011"
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
    "key" : "value"
  }
}