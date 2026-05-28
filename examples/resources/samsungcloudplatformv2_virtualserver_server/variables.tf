variable "name" {
  type    = string
  default = "terraform_test_server"
}

variable "state" {
  type    = string
  default = "ACTIVE"
}

variable "image_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S IMAGE_ID"
}

variable "server_type_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SERVER_TYPE_ID"
}

variable "keypair_name" {
  type    = string
  default = "terraform_test_keypair"
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
    delete_on_termination = true
    size                  = 16
    type                  = "SSD"
  }
}

variable "server_group_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SERVER_GROUP_ID"
}

variable "partition_number" {
  type    = number
  default = "1"
}

variable "security_groups" {
  type    = list(string)
  default = []
}

variable "networks_interface_1_subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S NETWORKS_INTERFACE_1_SUBNET_ID"
}

variable "extra_volumes_volume_1_size" {
  type    = number
  default = 16
}

variable "extra_volumes_volume_1_type" {
  type    = string
  default = "SSD"
}

variable "extra_volumes_volume_1_delete_on_termination" {
  type    = bool
  default = true
}



