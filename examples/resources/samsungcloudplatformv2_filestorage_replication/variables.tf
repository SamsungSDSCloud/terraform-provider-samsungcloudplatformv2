variable "volume_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S VOLUME_ID"
}

variable "name" {
  type    = string
  default = "backuptest"
}

variable "region" {
  type    = string
  default = "kr-west1"
}

variable "replication_frequency" {
  type    = string
  default = "hourly"
}

variable "cifs_password" {
  type    = string
  default = "ENTER YOUR RESOURCE'S CIFS_PASSWORD"
}

variable "backup_retention_count" {
  type    = number
  default = 10
}

variable "replication_type" {
  type    = string
  default = "backup"
}

variable "replication_update_type" {
  type    = string
  default = "replication update type"
}

variable "replication_policy" {
  type    = string
  default = "replication policy"
}


