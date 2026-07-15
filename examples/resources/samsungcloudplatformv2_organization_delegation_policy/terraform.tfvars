organization_id = "o-2b63982e88b74dbcb71ee972b13e2ce1"

document = {
  version = "2024-07-01"
  statement = [
    {
      sid        = "statement2"
      effect     = "Allow"
      action    = ["organization:CreateServiceControlPolicy",
                   "organization:ListAccounts"]
      principal = {
        "scp": [
          "srn:dev2::f045159f40c64125a1fe61bd71d1c14c:::iam:user/e3a6e3f99c1040639e9a3c8f8b7427de"
        ]
      }
      resource  = ["*"]
    }
  ]
}