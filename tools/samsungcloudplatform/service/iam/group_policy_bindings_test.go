package iam_test

import (
	"fmt"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupPolicyResourceTest(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupPolicyWithGroupAndPolicy("samsungcloudplatformv2_iam_policy.policy_a.id"),
			},
			{
				Config: testAccGroupPolicyWithGroupAndPolicy("samsungcloudplatformv2_iam_policy.policy_b.id"),
			},
			{
				Config: testAccGroupPolicyRemove(),
			},
			{
				Config: testAccGroupRemove(),
			},
		},
	})
}

func getGroupForGroupPolicy() string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_group" "group" {
				name = "test-acc-group"
  				description = "test-acc group desc"
				tags = {
					"test-acc-key" = "test-acc-value"
				}
				policy_ids = []
			}
	`)
}

func getPolicyAForGroupPolicy() string {
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

func getPolicyBForGroupPolicy() string {
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

func testAccGroupPolicyWithGroupAndPolicy(policyId string) string {
	return fmt.Sprintf(`
			%s
			%s
			%s
	
			resource "samsungcloudplatformv2_iam_group_policy_bindings" "group_policy_bindings" {
			  group_id = samsungcloudplatformv2_iam_group.group.id
			  policy_ids = [%s]
			}

		`, getGroupForGroupPolicy(), getPolicyAForGroupPolicy(), getPolicyBForGroupPolicy(), policyId)
}

func testAccGroupPolicyRemove() string {
	return fmt.Sprintf(`
			%s
			%s
			%s

		`, getGroupForGroupPolicy(), getPolicyAForGroupPolicy(), getPolicyBForGroupPolicy())
}

func testAccGroupRemove() string {
	return fmt.Sprintf(`
			%s
			%s

		`, getPolicyAForGroupPolicy(), getPolicyBForGroupPolicy())
}
