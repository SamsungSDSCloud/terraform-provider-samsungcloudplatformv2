variable "name" {
  type    = string
  default = ""
}

variable "state" {
  type    = string
  default = ""
}

variable "image_id" {
  type    = string
  default = ""
}

variable "server_type_id" {
  type    = string
  default = ""
}

variable "keypair_name" {
  type    = string
  default = ""
}

variable "lock" {
  type    = bool
  default = false
}

variable "user_data" {
  type    = string
  default = ""
}

variable "boot_volume" {
  type = object({
    size                  = number
    type                  = string
    delete_on_termination = bool
  })
  default = {
    delete_on_termination = false
    size                  = 0
    type                  = ""
  }
}

variable "server_group_id" {
  type    = string
  default = ""
}

variable "partition_number" {
  type    = number
  default = 0
}

variable "security_groups" {
  type    = list(string)
  default = [""]
}

variable "networks_interface_1_subnet_id" {
  type    = string
  default = ""
}

variable "extra_volumes_volume_1_size" {
  type    = number
  default = 0
}

variable "extra_volumes_volume_1_type" {
  type    = string
  default = ""
}

variable "extra_volumes_volume_1_delete_on_termination" {
  type    = bool
  default = false
}


