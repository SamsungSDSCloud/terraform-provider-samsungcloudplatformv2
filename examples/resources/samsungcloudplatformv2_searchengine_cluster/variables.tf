variable "allowable_ip_addresses" {
  type    = set(string)
  default = ["192.168.10.2/32"]
}

variable "dbaas_engine_version_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DBAAS_ENGINE_VERSION_ID"
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
    backup_option = {
      retention_period_day = "7"
      starting_time_hour   = "12"
    }
    database_port          = 9201
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
      role_type = "MASTER_DATA"
    }]
    role_type        = "MASTER_DATA"
    server_type_name = "ses1v2m4"
    }, {
    block_storage_groups = [{
      role_type   = "OS"
      size_gb     = 104
      volume_type = "SSD"
    }]
    instances = [{
      role_type = "KIBANA"
    }]
    role_type        = "KIBANA"
    server_type_name = "ses1v2m4"
  }]
}


variable "instance_name_prefix" {
  type    = string
  default = "searchb"
}

variable "name" {
  type    = string
  default = "searchb"
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
    key = "value"
  }
}

variable "license" {
  type    = string
  default = "{\n   \"license\":{\n      \"uid\":\"f5c002e6-c29d-4dde-bce4-cdb35bf85ab8\",\n      \"type\":\"trial\",\n      \"issue_date_in_millis\":1745193600000,\n      \"expiry_date_in_millis\":1752969599999,\n      \"max_nodes\":3,\n      \"issued_to\":\"S-Core Co., Ltd. (non-production environments)\",\n      \"issuer\":\"API\",\n      \"signature\":\"AAAAAwAAAA2ZzeOwjqKpsuwZQJCbAAABmC9ZN0hjZDBGYnVyRXpCOW5Bb3FjZDAxOWpSbTVoMVZwUzRxVk1PSmkxaktJRVl5MUYvUWh3bHZVUTllbXNPbzBUemtnbWpBbmlWRmRZb25KNFlBR2x0TXc2K2p1Y1VtMG1UQU9TRGZVSGRwaEJGUjE3bXd3LzRqZ05iLzRteWFNekdxRGpIYlFwYkJiNUs0U1hTVlJKNVlXekMrSlVUdFIvV0FNeWdOYnlESDc3MWhlY3hSQmdKSjJ2ZTcvYlBFOHhPQlV3ZHdDQ0tHcG5uOElCaDJ4K1hob29xSG85N0kvTWV3THhlQk9NL01VMFRjNDZpZEVXeUtUMXIyMlIveFpJUkk2WUdveEZaME9XWitGUi9WNTZVQW1FMG1DenhZU0ZmeXlZakVEMjZFT2NvOWxpZGlqVmlHNC8rWVVUYzMwRGVySHpIdURzKzFiRDl4TmM1TUp2VTBOUlJZUlAyV0ZVL2kvVk10L0NsbXNFYVZwT3NSU082dFNNa2prQ0ZsclZ4NTltbU1CVE5lR09Bck93V2J1Y3c9PQAAAQCdy6pBHZq1yaRofO7pDJJrEGmcnsUpa4BmsWywjbOU3zc3zLp6hvVuoWQ8bys6wl+lflToI51p6WinnyGmiFJoSvkNUMquuJMvEza5MlBzu4yb7mZKEUa4hxvz7IjQOltXP7KUXnMa98SrHx9Fkf/N80ZFhGD9t25UBkBgYvKEWTCztmOmNXTX/Sdq+JixfuijiA75EVGWhge7tGc6OpfHBgpZOJIxGOTAIUQiBxWK1ZdBF76CwhEVTbkMiKutvNsbwwo+yWiNGjq0mCUDYZUMXp6T4xk0VLzrmgGhBCqSR2BJzESq1Yk2VmIBv2Sn7JzokbGB3SB7FYQHKArVsinU\",\n      \"start_date_in_millis\":1745193600000\n   }\n}"
}


