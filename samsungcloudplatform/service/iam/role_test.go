package iam_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRoleCreate("test-acc-role", "test-acc role desc",
					map[string]string{
						"test-acc-key": "test-acc-value",
					}),
			},
			{
				Config: testAccRoleUpdate("test-acc role desc change"),
			},
		},
	})
}

func testAccRoleCreate(name string, description string, tags map[string]string) string {
	tagsJson, _ := json.Marshal(tags)

	return fmt.Sprintf(`
			%s

			resource "samsungcloudplatformv2_iam_role" "role" {
				  name = "%s"
				  description = "%s"
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
				  tags = %s
			}
	`, getAccountIdForRole(), name, description, tagsJson)
}

func testAccRoleUpdate(description string) string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_role" "role" {
				  description = "%s"
				  max_session_duration = 7200
			}
	`, description)
}

func getAccountIdForRole() string {
	return fmt.Sprintf(`
			data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
				limit = 1
			}
	`)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_role", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_role",
		F:    sweepRole,
	})
}

func sweepRole(region string) error {
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	r := iam.RoleDataSource{Size: types.Int32Value(1000)}
	roles, err := scpClient.Client.Iam.GetRoles(nil, r)

	var deleteResourceList []string

	for _, role := range roles.GetRoles() {
		if role.Description.Get() != nil {
			description := *role.Description.Get()

			if strings.HasPrefix(description, "test-acc") {
				deleteResourceList = append(deleteResourceList, role.Id)
			}
		}
	}

	for _, id := range deleteResourceList {
		err = scpClient.Client.Iam.DeleteRole(nil, id)
		if err != nil {
			return fmt.Errorf("error deleting role %s: %v", id, err)
		}
	}

	return nil
}
