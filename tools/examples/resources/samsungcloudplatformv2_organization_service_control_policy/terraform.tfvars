# terraform.tfvars
policy_name        = "TerraformPolicyTestModified"
policy_description = "terraform test! modify!"
policy_type        = "USER_DEFINED"
organization_id    = "o-2b63982e88b74dbcb71ee972b13e2ce1"

policy_document = {
  statement = [
    {
      sid    = "statement1"
      effect = "Allow"
      action = ["iam:AddRolePolicyBinding",
        "iam:DeleteRole",
        "iam:DeleteBulkRole",
        "iam:RemoveBulkRolePolicyBinding",
        "iam:SetRoleTrustPolicyPrincipals",
        "iam:SetRoleTrustPolicy",
      "iam:SetRole"]
      not_action = []
      principal  = "*"
      resource   = ["*"]
      condition  = {}
    }
  ]
  version = "2024-07-01"
}