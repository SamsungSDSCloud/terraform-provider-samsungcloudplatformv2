package iamidentitycenter_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform"
	iamidentitycentercommon "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common/iamidentitycenter"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamIdentityCenterGroupResourceTest(t *testing.T) {
	testGroupName := fmt.Sprintf("test-acc-g-%d", time.Now().UnixNano())
	testInstanceName := fmt.Sprintf("test-acc-i-%d", time.Now().UnixNano())

	existingInstanceId := getExistingInstanceId()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGroupCreate(
					testInstanceName,
					testGroupName,
					"test acc IAM Identity Center Group description",
					existingInstanceId,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_group.group", "name", testGroupName),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_group.group", "description", "test acc IAM Identity Center Group description"),
				),
			},
			{
				Config: testAccGroupUpdate(
					testInstanceName,
					testGroupName,
					"test acc IAM Identity Center Group updated description",
					existingInstanceId,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_group.group", "name", testGroupName),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_group.group", "description", "test acc IAM Identity Center Group updated description"),
				),
			},
		},
	})
}

func testAccGroupCreate(instanceName string, name string, description string, existingInstanceId string) string {
	if existingInstanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = "%s"
			name = "%s"
			description = "%s"
		}`, existingInstanceId, name, description)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name = "%s"
			description = "test instance for group test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name = "%s"
			description = "%s"
		}`, instanceName, name, description)
}

func testAccGroupUpdate(instanceName string, name string, description string, existingInstanceId string) string {
	if existingInstanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = "%s"
			name = "%s"
			description = "%s"
		}`, existingInstanceId, name, description)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name = "%s"
			description = "test instance for group test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_group" "group" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name = "%s"
			description = "%s"
		}`, instanceName, name, description)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_identity_center_group", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_identity_center_group",
		F:    sweepGroup,
	})
}

func sweepGroup(region string) error {
	ctx := context.Background()
	scpClient, err := iamidentitycentercommon.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	instances, err := scpClient.Client.IamIdentityCenter.ListInstances(ctx, 0, 0, "")
	if err != nil {
		return fmt.Errorf("error listing instances: %v", err)
	}

	for _, inst := range instances.GetInstances() {
		instanceId := inst.Id

		groups, err := scpClient.Client.IamIdentityCenter.ListGroups(ctx, instanceId, "", 100, 0, "", "", "")
		if err != nil {
			fmt.Printf("warning: error listing groups for instance %s: %v\n", instanceId, err)
			continue
		}

		for _, g := range groups.GetGroups() {
			if strings.HasPrefix(g.Name, "test-acc-g-") {
				_, _, err = scpClient.Client.IamIdentityCenter.DeleteGroup(ctx, g.Id, instanceId)
				if err != nil {
					fmt.Printf("warning: failed to delete group %s: %v\n", g.Id, err)
				}
			}
		}

		instanceDetail, _, err := scpClient.Client.IamIdentityCenter.GetInstance(ctx, instanceId)
		if err != nil {
			continue
		}
		if strings.HasPrefix(instanceDetail.Instance.Name, "test-acc-i-") {
			_, err = scpClient.Client.IamIdentityCenter.DeleteInstance(ctx, instanceId)
			if err != nil {
				fmt.Printf("warning: failed to delete instance %s: %v\n", instanceId, err)
			}
		}
	}

	return nil
}
