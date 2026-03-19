package iam_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserPolicyResourceTest(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserPolicyWithUserAndPolicy("samsungcloudplatformv2_iam_policy.policy_a.id"),
			},
			{
				Config: testAccUserPolicyWithUserAndPolicy("samsungcloudplatformv2_iam_policy.policy_b.id"),
			},
			{
				Config: testAccUserPolicyRemove(),
			},
			{
				Config: testAccUserRemove(),
			},
		},
	})
}

func getUserForUserPolicy() string {
	return fmt.Sprintf(`
			%s

			resource "samsungcloudplatformv2_iam_user" "user"{
				  account_id = data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys[0].account_id
				  description = "test-acc user desc"
				  password = "U2Ftc3VuZ3Nkc0Nsb3VkITIz"
				  user_name = "test-acc-user"
				  policy_ids = ["9bfa3a0668a146b6a90320743be400eb"]
                  password_reuse_count = 2
			}
	`, getAccountIdForUserPolicy())
}

func getPolicyAForUserPolicy() string {
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

func getPolicyBForUserPolicy() string {
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

func getAccountIdForUserPolicy() string {
	return fmt.Sprintf(`
			data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
				limit = 1
			}
	`)
}

func testAccUserPolicyWithUserAndPolicy(policyId string) string {
	return fmt.Sprintf(`
			%s
			%s
			%s
	
			resource "samsungcloudplatformv2_iam_user_policy_bindings" "user_policy_bindings" {
			  user_id = samsungcloudplatformv2_iam_user.user.user_id
			  policy_ids = [%s]
			}

		`, getUserForUserPolicy(), getPolicyAForUserPolicy(), getPolicyBForUserPolicy(), policyId)
}

func testAccUserPolicyRemove() string {
	return fmt.Sprintf(`
			%s
			%s
			%s

		`, getUserForUserPolicy(), getPolicyAForUserPolicy(), getPolicyBForUserPolicy())
}

func testAccUserRemove() string {
	return fmt.Sprintf(`
			%s
			%s

		`, getPolicyAForUserPolicy(), getPolicyBForUserPolicy())
}
