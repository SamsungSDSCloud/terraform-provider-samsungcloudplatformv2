variable "volume_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VOLUME_ID"
}

variable "snapshot_retention_count" {
  type    = number
  default = 10
}

variable "snapshot_schedule" {
  type = object({
    frequency   = string
    day_of_week = string
    hour        = string
  })
  default = {
    day_of_week = null
    frequency   = "DAILY"
    hour        = "10"
  }
}


