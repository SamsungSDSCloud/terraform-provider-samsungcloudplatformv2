variable "name" {
  type    = string
  default = ""
}

variable "os_distro" {
  type    = string
  default = ""
}

variable "disk_format" {
  type    = string
  default = ""
}

variable "container_format" {
  type    = string
  default = ""
}

variable "min_disk" {
  type    = number
  default = 0
}

variable "min_ram" {
  type    = number
  default = 0
}

variable "visibility" {
  type    = string
  default = ""
}

variable "protected" {
  type    = bool
  default = false
}

variable "url" {
  type    = string
  default = ""
}
variable "instance_id" {
  type    = string
  default = ""
}


