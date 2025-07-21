variable "account_id" {
  type = string
  default = "8a463aa4b1dc4f279c3f53b94dc45e74"
}
variable "bucket_name" {
  type = string
  default = "bucketname"
}

variable "bucket_region" {
  type = string
  default = "kr-west1"
}
variable "log_type_total_yn" {
  type = string
  default = "N"
}
variable "log_verification_yn" {
  type = string
  default = "N"
}
variable "region_total_yn" {
  type = string
  default = "Y"
}

variable "resource_type_total_yn" {
  type = string
  default = "Y"
}

variable "trail_description" {
  type = string
  default = "description info"
}

variable "trail_name" {
  type = string
  default = "TrailName"
}

variable "trail_save_type" {
  type = string
  default = "JSON"
}

variable "user_total_yn" {
  type = string
  default = "Y"
}


variable "tag_create_requests" {
  type = list(map(string))
  default = [
    {
      key   = "vpn_tag_key"
      value = "vpn_tag_value"
    },
    {
      key   = "another_tag_key"
      value = "another_tag_value"
    }
  ]
}


