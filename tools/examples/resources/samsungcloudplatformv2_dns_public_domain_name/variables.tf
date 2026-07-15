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
    address_type               = "DOMESTIC"
    auto_extension             = true
    description                = "hahaha"
    domestic_first_address_en  = "51-29 Jamsil-ro, Songpa-gu, Seoul"
    domestic_first_address_ko  = "서울특별시 송파구 잠실로 51-29 (잠실동)"
    domestic_second_address_en = "asd"
    domestic_second_address_ko = "ㅁㄴㅇ"
    name                       = "terraform22.com"
    overseas_first_address     = "123 Main St"
    overseas_second_address    = "New York"
    overseas_third_address     = "NY 10001"
    postal_code                = "05503"
    register_email             = "asd@asd.com"
    register_name_en           = "asd"
    register_name_ko           = "ㅁㄴㅇ"
    register_telno             = "010-1111-1111"
  }
}

variable "tag" {
  type = object({
    test_terraform_tag_key = string
  })
  default = {
    test_terraform_tag_key = "test_terraform_tag_value"
  }
}


