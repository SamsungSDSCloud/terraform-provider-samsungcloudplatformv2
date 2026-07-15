variable "size" {
  type    = number
  default = 10000
}

variable "page" {
  type    = number
  default = null
}

variable "sort" {
  type    = string
  default = null
}

variable "name" {
  type    = string
  default = null
}

variable "subnet_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S SUBNET_ID"
}

variable "status" {
  type    = list(string)
  default = []
}

variable "kubernetes_version" {
  type    = list(string)
  default = []
}




variable "clusters_region" {
  type    = string
  default = "kr-west1"
}

variable "clusters_filter_name" {
  type    = string
  default = "name"
}

variable "clusters_filter_values" {
  type    = list(string)
  default = ["test"]
}

variable "clusters_filter_use_regex" {
  type    = bool
  default = true
}



