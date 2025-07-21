variable "volume_id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "snapshot_retention_count" {
  type = number
  default = 10
}

variable "snapshot_schedule" {
  type = object({
    frequency = string
    day_of_week = string
    hour = string
  })
  default = {
    frequency = "DAILY"
    day_of_week = null
    hour = "10"
  }
}