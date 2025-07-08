variable "account_id" {
  type    = string
  default = ""
}
variable "bucket_name" {
  type    = string
  default = ""
}

variable "bucket_region" {
  type    = string
  default = ""
}
variable "log_type_total_yn" {
  type    = string
  default = ""
}
variable "log_verification_yn" {
  type    = string
  default = ""
}
variable "region_total_yn" {
  type    = string
  default = ""
}

variable "resource_type_total_yn" {
  type    = string
  default = ""
}

variable "trail_description" {
  type    = string
  default = ""
}

variable "trail_name" {
  type    = string
  default = ""
}

variable "trail_save_type" {
  type    = string
  default = ""
}

variable "user_total_yn" {
  type    = string
  default = ""
}


variable "tag_create_requests" {
  type    = list(map(string))
  default = [null]
}




