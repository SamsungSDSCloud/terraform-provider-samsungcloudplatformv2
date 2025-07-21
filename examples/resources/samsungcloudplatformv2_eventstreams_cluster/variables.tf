variable "allowable_ip_addresses" {
  type = list(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "189299a34f464cac94a24f2d8d57afec"
}

variable "akhq_enabled" {
  type    = bool
  default = false
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
    zookeeper_sasl_id       = string
    zookeeper_sasl_password = string
    zookeeper_port          = number
    broker_sasl_id          = string
    broker_sasl_password    = string
    broker_port             = number
  })
  default = {
    zookeeper_sasl_id       = "eventstream",
    zookeeper_sasl_password = "Password001!",
    zookeeper_port          = 2180,
    broker_sasl_id          = "eventstream",
    broker_sasl_password    = "Password001!",
    broker_port             = 9091
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
      role_type        = "ZOOKEEPER_BROKER"
      server_type_name = "es1v2m4"
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
          "role_type" : "ZOOKEEPER_BROKER"
        }
      ]
    }
  ]
}


variable "instance_name_prefix" {
  type    = string
  default = "eventa"
}

variable "name" {
  type    = string
  default = "eventa"
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
    "key" : "value",
    "key1" : "value1"
  }
}