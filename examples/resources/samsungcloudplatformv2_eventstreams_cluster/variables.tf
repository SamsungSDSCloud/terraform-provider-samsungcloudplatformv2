variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.1/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID"
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
    broker_port             = 9091
    broker_sasl_id          = "ENTER YOUR RESOURCE'S BROKER_SASL_ID"
    broker_sasl_password    = "ENTER YOUR RESOURCE'S BROKER_SASL_PASSWORD"
    zookeeper_port          = 2180
    zookeeper_sasl_id       = "ENTER YOUR RESOURCE'S ZOOKEEPER_SASL_ID"
    zookeeper_sasl_password = "ENTER YOUR RESOURCE'S ZOOKEEPER_SASL_PASSWORD"
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
      role_type = "ZOOKEEPER_BROKER"
    }]
    role_type        = "ZOOKEEPER_BROKER"
    server_type_name = "ess1v2m8"
  }]
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
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "timezone" {
  type    = string
  default = "Asia/Seoul"
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
    key  = "value"
    key1 = "value1"
  }
}

variable "service_watch_log_collection" {
  type    = bool
  default = false
}


