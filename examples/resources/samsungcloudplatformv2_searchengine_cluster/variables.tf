variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}

variable "is_combined" {
  type    = bool
  default = true
}

variable "nat_enabled" {
  type    = bool
  default = false
}

variable "init_config_option" {
  type = object({
    database_port          = number
    database_user_name     = string
    database_user_password = string
    backup_option = object({
      retention_period_day = string
      starting_time_hour   = string
    })
  })
  default = {
    database_port          = 9201
    database_user_name     = "username"
    database_user_password = "password001!"
    backup_option = {
      retention_period_day = "7"
      starting_time_hour   = "12"
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
      role_type = string
    }))
  }))
  default = [
    {
      role_type        = "MASTER_DATA"
      server_type_name = "ses1v2m4"
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
          "role_type" : "MASTER_DATA",
        }
      ]
    },
    {
      role_type        = "KIBANA"
      server_type_name = "ses1v2m4"
      block_storage_groups = [
        {
          "role_type" : "OS",
          "volume_type" : "SSD",
          "size_gb" : 104
        }
      ]
      instances = [
        {
          "role_type" : "KIBANA",
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
  type    = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
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
    period_hour            = "0.5"
    starting_day_of_week   = "MON"
    starting_time          = "0000"
    use_maintenance_option = true
  }
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

variable "license" {
  type    = string
  default = "{\n   \"license\":{\n      \"uid\":\"8a463aa4-b1dc-4f27-9c3f-53b94dc45e74\",\n      \"type\":\"trial\",\n      \"issue_date_in_millis\":1745193600000,\n      \"expiry_date_in_millis\":1752969599999,\n      \"max_nodes\":3,\n      \"issued_to\":\"Samsung SDS Co., Ltd.\",\n      \"issuer\":\"API\",\n      \"signature\":\"signature data\",\n      \"start_date_in_millis\":1745193600000\n   }\n}"
}