variable "size" {
  type    = number
  default = 0
}

variable "page" {
  type    = number
  default = 0
}

variable "sort" {
  type    = string
  default = ""
}

variable "name" {
  type    = string
  default = ""
}

variable "subnet_id" {
  type    = string
  default = ""
}

variable "status" {
  type    = list(string)
  default = [""]
}

variable "kubernetes_version" {
  type    = list(string)
  default = [""]
}




variable "clusters_region" {
  type    = string
  default = ""
}

variable "clusters_filter_name" {
  type    = string
  default = ""
}

variable "clusters_filter_values" {
  type    = list(string)
  default = [""]
}

variable "clusters_filter_use_regex" {
  type    = bool
  default = false
}


