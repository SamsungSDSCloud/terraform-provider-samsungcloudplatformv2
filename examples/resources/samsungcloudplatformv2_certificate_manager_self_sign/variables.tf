variable "cn" {
  type    = string
  default = "test.go.kr.sds"
}

variable "not_after_dt" {
  type    = string
  default = "20251231"
}

variable "not_before_dt" {
  type    = string
  default = "20251212"
}

variable "name" {
  type    = string
  default = "test_go_gdcv_251027"
}

variable "organization" {
  type    = string
  default = "samsungSDS"
}

variable "region" {
  type    = string
  default = "kr-west1"
}

variable "timezone" {
  type    = string
  default = "Asia/Seoul"
}

variable "recipients" {
  type = list(map(string))
  default = [{
    region    = "Asia/Seoul"
    user_id   = "ENTER YOUR RESOURCE'S USER_ID"
    user_name = "userA@samsung.com"
  }]
}

variable "tags" {
  type = map(string)
  default = {
    test_tag_key  = "test_tag_value"
    test_tag_key2 = "test_tag_value2"
  }
}


