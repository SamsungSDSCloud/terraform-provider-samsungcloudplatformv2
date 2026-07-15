variable "dashboard_name" {
  type    = string
  default = "tf_test_01"
}

variable "height" {
  type    = number
  default = 1
}

variable "width" {
  type    = number
  default = 1
}

variable "order" {
  type    = number
  default = 1
}


variable "period" {
  type    = number
  default = 300
}

variable "statistic_type" {
  type    = string
  default = "AVG"
}

variable "title" {
  type    = string
  default = "Virtual Server | CPU Usage"
}

variable "metric_name" {
  type    = string
  default = "CPU Usage"
}

variable "namespace_name" {
  type    = string
  default = "Virtual Server"
}

variable "display_name" {
  type    = string
  default = "Virtual Server | CPU Usage"
}

variable "dimension_key" {
  type    = string
  default = "resource_id"
}

variable "dimension_value" {
  type    = string
  default = "d5b49100-e3e3-4d10-b2e9-9da68aed7747"
}

variable "color" {
  type    = string
  default = "#ff7f0e"
}


