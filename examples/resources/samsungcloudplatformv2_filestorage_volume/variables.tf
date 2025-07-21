variable "name" {
  type = string
  default = "my_volume"
}

variable "protocol" {
  type = string
  default = "NFS"
}

variable "type_name" {
  type = string
  default = "HDD"
}

variable "cifs_password" {
  type = string
  default = "cifspwd0!!"
}

variable "file_unit_recovery_enabled" {
  type = bool
  default = true
}

variable "access_rules" {
  type = list(object({
    object_type = string,
    object_id = string
  }))
  default = [{
    object_type="VM",
    object_id="8a463aa4-b1dc-4f27-9c3f-53b94dc45e74"
  }]
}

variable "tags" {
  type    = map(string)
  default = {
    "terraform_key" = "terraform_value"
  }
}