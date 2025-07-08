provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_billing_planned_computes" "planned_computes" {
  account_id       = var.account_id # 계정 ID
  contract_type    = var.contract_type                              # 계약 유형
  os_type          = var.os_type                         # OS 유형
  server_type      = var.server_type
  service_id       = var.service_id    
  service_name     = var.service_name            # 리소스 이름 (server_name으로 매핑
  action           = var.action
  region           = var.region
  tags = {
    "tf_key1": "tf_val1"
  }
}