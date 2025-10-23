provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_loggingaudit_trail" "trail" {

  account_id = var.account_id
  bucket_name = var.bucket_name
  bucket_region = var.bucket_region
  log_type_total_yn = var.log_type_total_yn
  log_verification_yn = var.log_type_total_yn
  region_names=[]
  region_total_yn= var.region_total_yn
  resource_type_total_yn= var.resource_type_total_yn
  tag_create_requests =  [
                            {
                              "key": "111",
                              "value": "111"
                            }
                          ]
  target_log_types= []
  target_resource_types=  ["iam:access-key",  "iam:group",]
  target_users= []
  trail_description= var.trail_description
  trail_name= var.trail_name
  trail_save_type= var.trail_save_type
  user_total_yn= var.user_total_yn
  organization_trail_yn = var.organization_trail_yn
  log_archive_account_id = var.log_archive_account_id
}


