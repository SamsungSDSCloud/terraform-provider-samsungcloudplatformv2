variable "account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S ACCOUNT_ID"
}
variable "bucket_name" {
  type    = string
  default = "scoretest1"
}

variable "bucket_region" {
  type    = string
  default = "kr-west1"
}
variable "log_type_total_yn" {
  type    = string
  default = "N"
}
variable "log_verification_yn" {
  type    = string
  default = "N"
}
variable "region_total_yn" {
  type    = string
  default = "Y"
}

variable "resource_type_total_yn" {
  type    = string
  default = "Y"
}

variable "trail_description" {
  type    = string
  default = "This is a test trail"
}

variable "trail_name" {
  type    = string
  default = "TestTrail200321"
}

variable "trail_save_type" {
  type    = string
  default = "JSON"
}

variable "user_total_yn" {
  type    = string
  default = "Y"
}

variable "organization_trail_yn" {
  type    = string
  default = "N"
}

variable "log_archive_account_id" {
  type    = string
  default = "ENTER YOUR RESOURCE'S LOG_ARCHIVE_ACCOUNT_ID"
}

variable "tag_create_requests" {
  type = list(map(string))
  default = [{
    key   = "vpn_tag_key"
    value = "vpn_tag_value"
    }, {
    key   = "another_tag_key"
    value = "another_tag_value"
  }]
}





