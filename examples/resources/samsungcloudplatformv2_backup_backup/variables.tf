variable "region" {
  type    = string
  default = "kr-west1"
}

variable "name" {
  type = string
  default = "terraformtestbackup01"
}

variable "policy_category" {
  type = string
  default = "AGENTLESS"
}

variable "policy_type" {
  type = string
  default = "VM_IMAGE"
}

variable "server_uuid" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "server_category" {
  type = string
  default = "VIRTUAL_SERVER"
}

variable "encrypt_enabled" {
  type = bool
  default = false
}

variable "retention_period" {
  type = string
  default = "MONTH_1"
}

variable "schedules" {
  type = list(object({
    type = string
    frequency = string
    start_time = string
    start_day = string
    start_week = string
  }))
  default = [
    {
      type = "FULL"
      frequency = "DAILY"
      start_time = "11:00:00"
      start_day = null
      start_week = null
    },
    {
      type = "INCREMENTAL"
      frequency = "WEEKLY"
      start_time = "12:30:00"
      start_day = "THU"
      start_week = null
    },
    {
      type = "INCREMENTAL"
      frequency = "MONTHLY"
      start_time = "13:00:00"
      start_day = "FRI"
      start_week = "WEEK_3"
    }]
}

variable "tags" {
  type    = map(string)
  default = {
    "test_tag_key": "tag_key",
    "test_tag_key2": "tag_key_2"
  }
}
