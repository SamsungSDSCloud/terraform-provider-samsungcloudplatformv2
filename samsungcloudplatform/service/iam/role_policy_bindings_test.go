package iam_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRolePolicyResourceTest(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRolePolicyWithRoleAndPolicy("samsungcloudplatformv2_iam_policy.policy_a.id"),
			},
			{
				Config: testAccRolePolicyWithRoleAndPolicy("samsungcloudplatformv2_iam_policy.policy_b.id"),
			},
			{
				Config: testAccRolePolicyRemove(),
			},
			{
				Config: testAccRoleRemove(),
			},
		},
	})
}

func getRoleForRolePolicy() string {
	return fmt.Sprintf(`
			%s

			resource "samsungcloudplatformv2_iam_role" "role" {
				  name = "test-acc-role"
				  description = "test-acc role desc"
				  max_session_duration = 3600
				  assume_role_policy_document = {
					"version" = "2024-07-01",
					"statement" = [
					  {
						"action" = [
						  "sts:AssumeRole"
						],
						"condition" = {},
						"effect"    = "Allow",
						"principal" = {
						  "principal_map" = {
							"Account" = [data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys[0].account_id]
						  }
						},
						"resource" = ["*"],
						"sid"      = "VisualEditor0"
					  }
					]
				  }
				  policy_ids = ["9bfa3a0668a146b6a90320743be400eb"]
			}
	`, getAccountIdForRolePolicy())
}

func getPolicyAForRolePolicy() string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_policy" "policy_a"{
				policy_name = "test-acc-policy-a"
				description = "test-acc policy desc"
				policy_version = {
									"policy_document" = {
									  "version" = "2024-07-01",
									  "statement" = [
										{
										  "action" = [
											"iam:*"
										  ],
										  "effect" = "Allow",
										  "resource" = [
											"*"
										  ],
										  "sid" = "VisualEditor0"
										}
									  ]
									}
								  }
			}
	`)
}

func getPolicyBForRolePolicy() string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_policy" "policy_b"{
				policy_name = "test-acc-policy-b"
				description = "test-acc policy desc"
				policy_version = {
									"policy_document" = {
									  "version" = "2024-07-01",
									  "statement" = [
										{
										  "action" = [
											"organization:*"
										  ],
										  "effect" = "Allow",
										  "resource" = [
											"*"
										  ],
										  "sid" = "VisualEditor0"
										}
									  ]
									}
								  }
			}
	`)
}

func getAccountIdForRolePolicy() string {
	return fmt.Sprintf(`
			data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
				limit = 1
			}
	`)
}

func testAccRolePolicyWithRoleAndPolicy(policyId string) string {
	return fmt.Sprintf(`
			%s
			%s
			%s
	
			resource "samsungcloudplatformv2_iam_role_policy_bindings" "role_policy_bindings" {
			  role_id = samsungcloudplatformv2_iam_role.role.id
			  policy_ids = [%s]
			}

		`, getRoleForRolePolicy(), getPolicyAForRolePolicy(), getPolicyBForRolePolicy(), policyId)
}

func testAccRolePolicyRemove() string {
	return fmt.Sprintf(`
			%s
			%s
			%s

		`, getRoleForRolePolicy(), getPolicyAForRolePolicy(), getPolicyBForRolePolicy())
}

func testAccRoleRemove() string {
	return fmt.Sprintf(`
			%s
			%s

		`, getPolicyAForRolePolicy(), getPolicyBForRolePolicy())
}
