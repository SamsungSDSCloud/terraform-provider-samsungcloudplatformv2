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
    address_type               = ""
    auto_extension             = false
    description                = ""
    domestic_first_address_en  = ""
    domestic_first_address_ko  = ""
    domestic_second_address_en = ""
    domestic_second_address_ko = ""
    name                       = ""
    overseas_first_address     = ""
    overseas_second_address    = ""
    overseas_third_address     = ""
    postal_code                = ""
    register_email             = ""
    register_name_en           = ""
    register_name_ko           = ""
    register_telno             = ""
  }
}

variable "tag" {
  type = object({
    test_terraform_tag_key = string
  })
  default = {
    test_terraform_tag_key = ""
  }
}

