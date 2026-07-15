provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_cloudmonitoring_producttypes" "producttypes" {
  product_category_code = var.ProductCategoryCode
}
