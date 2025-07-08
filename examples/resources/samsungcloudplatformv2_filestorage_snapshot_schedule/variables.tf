variable "volume_id" {
  type    = string
  default = ""
}

variable "snapshot_retention_count" {
  type    = number
  default = 0
}

variable "snapshot_schedule" {
  type = object({
    frequency   = string
    day_of_week = string
    hour        = string
  })
  default = {
    day_of_week = ""
    frequency   = ""
    hour        = ""
  }
}

