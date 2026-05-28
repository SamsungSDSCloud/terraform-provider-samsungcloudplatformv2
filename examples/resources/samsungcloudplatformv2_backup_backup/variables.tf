variable "region" {
  type    = string
  default = "kr-west1"
}

variable "name" {
  type    = string
  default = "terraformtestbackup01"
}

variable "policy_category" {
  type    = string
  default = "AGENTLESS"
}

variable "policy_type" {
  type    = string
  default = "VM_IMAGE"
}

variable "server_uuid" {
  type    = string
  default = "a16687f2-3abc-4f40-bb5d-ee79ea21249d"
}

variable "server_category" {
  type    = string
  default = "VIRTUAL_SERVER"
}

variable "encrypt_enabled" {
  type    = bool
  default = true
}

variable "retention_period" {
  type    = string
  default = "MONTH_1"
}

variable "schedules" {
  type = list(object({
    type       = string
    frequency  = string
    start_time = string
    start_day  = string
    start_week = string
  }))
  default = [{
    frequency  = "DAILY"
    start_day  = null
    start_time = "11:00:00"
    start_week = null
    type       = "FULL"
    }, {
    frequency  = "WEEKLY"
    start_day  = "THU"
    start_time = "12:30:00"
    start_week = null
    type       = "INCREMENTAL"
    }, {
    frequency  = "MONTHLY"
    start_day  = "FRI"
    start_time = "13:00:00"
    start_week = "WEEK_3"
    type       = "INCREMENTAL"
  }]
}

variable "tags" {
  type = map(string)
  default = {
    test_tag_key  = "test_tag_value"
    test_tag_key2 = "test_tag_value2"
  }
}



