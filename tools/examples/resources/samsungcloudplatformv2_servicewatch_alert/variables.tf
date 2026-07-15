variable "alert_name" {
  type    = string
  default = "tf_alert_test_01"
}

variable "alert_description" {
  type    = string
  default = "Kubernetes Alert for Terraform."
}

variable "alert_type" {
  type    = string
  default = "METRIC_ALERT"
}

variable "alert_activated_yn" {
  type    = string
  default = "Y"
}

variable "namespace_name" {
  type    = string
  default = "Kubernetes Engine"
}

variable "metric_name" {
  type    = string
  default = "namespace_number_of_running_pods"
}

variable "dimension_key_1" {
  type    = string
  default = "resource_id"
}

variable "dimension_value_1" {
  type    = string
  default = "7dcc310f6d3841cd82480cf0f1654bba"
}

variable "dimension_key_2" {
  type    = string
  default = "namespace"
}

variable "dimension_value_2" {
  type    = string
  default = "kube-system"
}

variable "period" {
  type    = number
  default = 300
}

variable "statistic_type" {
  type    = string
  default = "SUM"
}

variable "level" {
  type    = string
  default = "MIDDLE"
}

variable "operator" {
  type    = string
  default = "GTE"
}

variable "threshold" {
  type    = number
  default = 3
}

variable "missing_data_option" {
  type    = string
  default = "MISSING"
}



