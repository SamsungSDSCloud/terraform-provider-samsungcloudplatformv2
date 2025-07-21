variable "public_domain_name" {
  type = object({
    name                       = string,
    description                = string,
    register_name_ko           = string,
    register_name_en           = string,
    register_email             = string,
    address_type               = string,
    postal_code                = string,
    domestic_first_address_ko  = string,
    domestic_second_address_ko = string,
    domestic_first_address_en  = string,
    domestic_second_address_en = string,
    overseas_first_address     = string,
    overseas_second_address    = string,
    overseas_third_address     = string,
    register_telno             = string,
    auto_extension             = bool
  })
  default = {
    name                       = "terraform.com",
    description                = "description info",
    register_name_ko           = "생성이름",
    register_name_en           = "regsisterName",
    register_email             = "sung.sam@samsung.com",
    address_type               = "DOMESTIC",
    postal_code                = "05503",
    domestic_first_address_ko  = "서울특별시 송파구 올림픽로 35길 125 (신천동)",
    domestic_second_address_ko = "삼성 SDS 타워",
    domestic_first_address_en  = "125, Olympic-ro 35-gil, Songpa-gu, Seoul",
    domestic_second_address_en = "Samsung SDS Tower",
    overseas_first_address     = "",
    overseas_second_address    = "",
    overseas_third_address     = "",
    register_telno             = "010-1234-5678",
    auto_extension             = true
  }
}

variable "tag" {
  type = object({
    terraform_tag_key = string
  })
  default = {
    terraform_tag_key = "terraform_tag_value"
  }
}