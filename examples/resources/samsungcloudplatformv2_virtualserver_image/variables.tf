variable "name" {
  type = string
  default = "terraform_image"
}

variable "os_distro" {
  type = string
  default = "ubuntu"
}

variable "disk_format" {
  type = string
  default = "qcow2"
}

variable "container_format" {
  type = string
  default = "bare"
}

variable "min_disk" {
  type = number
  default = 10
}

variable "min_ram" {
  type = number
  default = 1
}

variable "visibility" {
  type = string
  default = "private"
}

variable "protected" {
  type = bool
  default = false
}

variable "url" {
  type = string
  default = "https://object-store.private.kr-west1.s.samsungsdscloud.com/8a463aa4b1dc4f279c3f53b94dc45e74:terraformvm/ubuntu.qcow2"
}
variable "instance_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}
