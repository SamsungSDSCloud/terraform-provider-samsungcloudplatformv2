variable "volume_name" {
  type    = string
  default = "test_access_list"
}

variable "protocol" {
  type    = string
  default = "CIFS"
}

variable "type_name" {
  type    = string
  default = "HDD"
}

variable "cifs_password" {
  type    = string
  default = "ENTER YOUR RESOURCE'S CIFS_PASSWORD"
}

variable "file_unit_recovery_enabled" {
  type    = bool
  default = false
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform = "test_terraform_value"
  }
}

variable "zone" {
  type    = string
  default = "kr-west1-a"
}

variable "object_type" {
  type    = string
  default = "VM"
}

variable "object_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S OBJECT_ID"
}



