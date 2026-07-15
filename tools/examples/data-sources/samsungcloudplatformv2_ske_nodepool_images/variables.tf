variable "scp_original_image_type" {
  type    = string
  default = "k8s"
}

variable "kubernetes_version" {
  type    = string
  default = "v1.34.3"
}

variable "os" {
  type    = string
  default = null
}

variable "size" {
  type    = string
  default = null
}

variable "page" {
  type    = string
  default = null
}

variable "sort" {
  type    = string
  default = null
}


