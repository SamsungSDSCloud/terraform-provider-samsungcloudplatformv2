package iam_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	util "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/iam"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyCreate("test-acc-policy", "test-acc policy desc",
					map[string]string{
						"test-acc-key": "test-acc-value",
					}),
			},
			{
				Config: testAccPolicyUpdate("change-test-acc-policy", "test-acc policy change desc"),
			},
		},
	})
}

func testAccPolicyCreate(name string, description string, tags map[string]string) string {
	tagsJson, _ := json.Marshal(tags)

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_policy" "policy"{
				policy_name = "%s"
				description = "%s"
				tags = %s
				policy_version = {
									"policy_document" = {
									  "version" = "2024-07-01",
									  "statement" = [
										{
										  "action" = [
											"*"
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
	`, name, description, tagsJson)
}

func testAccPolicyUpdate(name string, description string) string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_policy" "policy"{
				policy_name = "%s"
				description = "%s"
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
	`, name, description)

}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_policy", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_policy",
		F:    sweepPolicy,
	})
}

func sweepPolicy(region string) error {
	scpClient, err := util.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	r := iam.PolicyDatasource{Size: types.Int32Value(1000)}
	policies, err := scpClient.Client.Iam.GetPolicies(nil, r)

	var deleteResourceList []string

	for _, policy := range policies.GetPolicies() {
		if policy.Description.Get() != nil {
			description := *policy.Description.Get()

			if strings.HasPrefix(description, "test-acc") {
				deleteResourceList = append(deleteResourceList, *policy.Id)
			}

		}
	}

	for _, id := range deleteResourceList {
		err = scpClient.Client.Iam.DeletePolicy(nil, id)
		if err != nil {
			return fmt.Errorf("error deleting policy %s: %v", id, err)
		}
	}

	return nil
}
