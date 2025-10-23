variable "name" {
  type    = string
  default = ""
}

variable "protocol" {
  type    = string
  default = ""
}

variable "type_name" {
  type    = string
  default = ""
}

variable "cifs_password" {
  type    = string
  default = ""
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
    object_id   = ""
    object_type = ""
  }]
}

variable "tags" {
  type    = map(string)
  default = null
}


