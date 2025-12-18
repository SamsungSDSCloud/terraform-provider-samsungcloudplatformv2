variable "with_count" {
  type    = string
  default = ""
}

variable "limit" {
  type    = number
  default = 0
}

variable "marker" {
  type    = string
  default = ""
}

variable "sort" {
  type    = string
  default = ""
}

variable "is_mine" {
  type    = bool
  default = false
}

variable "diagnosis_id" {
  type    = string
  default = ""
}

variable "diagnosis_name" {
  type    = string
  default = ""
}

variable "csp_type" {
  type    = string
  default = ""
}

variable "diagnosis_account_id" {
  type    = string
  default = ""
}

variable "recent_diagnosis_state" {
  type    = list(string)
  default = [""]
}

variable "start_date" {
  type    = string
  default = ""
}

variable "end_date" {
  type    = string
  default = ""
}


