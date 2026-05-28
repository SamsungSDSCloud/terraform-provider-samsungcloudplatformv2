variable "name" {
  type    = string
  default = "terraform_test_volume"
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

variable "access_rules" {
  type = list(object({
    object_type = string,
    object_id   = string
  }))
  default = [{
    object_id   = "ENTER YOUR RESOURCE'S OBJECT_ID"
    object_type = "VM"
  }]
}

variable "tags" {
  type = map(string)
  default = {
    test_terraform = "test_terraform_value"
  }
}



