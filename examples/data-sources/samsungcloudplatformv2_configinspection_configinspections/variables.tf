variable "with_count" {
  type    = string
  default = null // true, false
}

variable "limit" {
  type    = number
  default = null
}

variable "marker" {
  type    = string
  default = null
}

variable "sort" {
  type    = string
  default = null
}

variable "is_mine" {
  type    = bool
  default = null
}

variable "diagnosis_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DIAGNOSIS_ID"
}

variable "diagnosis_name" {
  type    = string
  default = null
}

variable "csp_type" {
  type    = string
  default = null
}

variable "diagnosis_account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S DIAGNOSIS_ACCOUNT_ID"
}

variable "recent_diagnosis_state" {
  type    = list(string)
  default = null
}

variable "start_date" {
  type    = string
  default = null
}

variable "end_date" {
  type    = string
  default = null
}



