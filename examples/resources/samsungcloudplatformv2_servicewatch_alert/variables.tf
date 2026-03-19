variable "alert_name" {
  type    = string
  default = ""
}

variable "alert_description" {
  type    = string
  default = ""
}

variable "alert_type" {
  type    = string
  default = ""
}

variable "alert_activated_yn" {
  type    = string
  default = ""
}

variable "namespace_name" {
  type    = string
  default = ""
}

variable "metric_name" {
  type    = string
  default = ""
}

variable "dimension_key_1" {
  type    = string
  default = ""
}

variable "dimension_value_1" {
  type    = string
  default = ""
}

variable "dimension_key_2" {
  type    = string
  default = ""
}

variable "dimension_value_2" {
  type    = string
  default = ""
}

variable "period" {
  type    = number
  default = 0
}

variable "statistic_type" {
  type    = string
  default = ""
}

variable "level" {
  type    = string
  default = ""
}

variable "operator" {
  type    = string
  default = ""
}

variable "threshold" {
  type    = number
  default = 0
}

variable "missing_data_option" {
  type    = string
  default = ""
}


