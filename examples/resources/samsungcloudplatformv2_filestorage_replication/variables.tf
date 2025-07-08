variable "volume_id" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "region" {
  type    = string
  default = ""
}

variable "replication_frequency" {
  type    = string
  default = ""
}

variable "cifs_password" {
  type    = string
  default = ""
}

variable "backup_retention_count" {
  type    = number
  default = 0
}

variable "replication_type" {
  type    = string
  default = ""
}

variable "replication_update_type" {
  type    = string
  default = ""
}

variable "replication_policy" {
  type    = string
  default = ""
}

