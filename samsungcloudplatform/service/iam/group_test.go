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

func TestAccGroupResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupCreate("test-acc-group", "group for test-acc",
					map[string]string{
						"test-acc-key": "test-acc-value",
					},
					[]string{"d258561c5cf04106bb0f2a8d02a7479e"}),
			},
			{
				Config: testAccGroupUpdate("change-test-acc-group", "change group for test-acc"),
			},
		},
	})
}

func testAccGroupCreate(name string, description string, tags map[string]string, policyIds []string) string {
	tagsJson, _ := json.Marshal(tags)
	policyIdsJson, _ := json.Marshal(policyIds)

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_group" "group" {
				name = "%s"
  				description = "%s"
				tags = %s
				policy_ids = %s
			}`, name, description, tagsJson, policyIdsJson)
}

func testAccGroupUpdate(name string, description string) string {

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_group" "group" {
				name = "%s"
  				description = "%s"
			}`, name, description)
}
func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_group", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_group",
		F:    sweepGroup,
	})
}

func sweepGroup(region string) error {
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	r := iam.GroupDataSource{Size: types.Int32Value(1000)}
	groups, err := scpClient.Client.Iam.GetGroups(nil, r)

	var deleteResourceList []string

	for _, group := range groups.GetGroups() {
		if group.Description.Get() != nil {
			description := *group.Description.Get()

			if strings.HasPrefix(description, "test-acc") {
				deleteResourceList = append(deleteResourceList, group.Id)
			}
		}

	}

	for _, id := range deleteResourceList {
		err = scpClient.Client.Iam.DeleteGroup(nil, id)
		if err != nil {
			return fmt.Errorf("error deleting group %s: %v", id, err)
		}
	}

	return nil
}
