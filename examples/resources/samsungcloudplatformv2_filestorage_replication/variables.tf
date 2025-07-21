variable "volume_id" {
  type = string
  default = "8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
}

variable "name" {
  type = string
  default = "my_volume"
}

variable "region" {
  type = string
  default = "kr-west1"
}

variable "replication_frequency" {
  type = string
  default = "5min"
}

variable "cifs_password" {
  type = string
  default = "cifspwd0!!"
}

variable "backup_retention_count" {
  type = number
  default = 10
}

variable "replication_type"{
  type = string
  default = "backup"
}

variable "replication_update_type" {
  type = string
  default = "policy"
}

variable "replication_policy" {
  type = string
  default = "use"
}