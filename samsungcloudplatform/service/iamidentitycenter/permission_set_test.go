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

const (
	testPermissionSetResourceRef = "samsungcloudplatformv2_iam_identity_center_permission_set.permission_set"
)

func TestAccIamIdentityCenterPermissionSetResourceTest(t *testing.T) {
	existingInstanceId := getExistingInstanceId()
	testPermissionSetName := fmt.Sprintf("test-acc-ps-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionSetCreate(
					existingInstanceId,
					testPermissionSetName,
					"test acc IAM Identity Center Permission Set description",
					3600,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testPermissionSetResourceRef, "name", testPermissionSetName),
					resource.TestCheckResourceAttr(testPermissionSetResourceRef, "description", "test acc IAM Identity Center Permission Set description"),
					resource.TestCheckResourceAttr(testPermissionSetResourceRef, "session_duration", "3600"),
					resource.TestCheckResourceAttrSet(testPermissionSetResourceRef, "id"),
					resource.TestCheckResourceAttrSet(testPermissionSetResourceRef, "srn"),
					resource.TestCheckResourceAttrSet(testPermissionSetResourceRef, "state"),
				),
			},
			{
				Config: testAccPermissionSetUpdate(
					existingInstanceId,
					testPermissionSetName,
					"test acc IAM Identity Center Permission Set updated description",
					7200,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testPermissionSetResourceRef, "name", testPermissionSetName),
					resource.TestCheckResourceAttr(testPermissionSetResourceRef, "description", "test acc IAM Identity Center Permission Set updated description"),
					resource.TestCheckResourceAttr(testPermissionSetResourceRef, "session_duration", "7200"),
				),
			},
			{
				Config: testAccPermissionSetDelete(
					existingInstanceId,
					testPermissionSetName,
					"test acc IAM Identity Center Permission Set updated description",
					7200,
				),
			},
		},
	})
}

func testAccPermissionSetCreate(instanceId string, name string, description string, sessionDuration int) string {
	if instanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = "%s"
			name             = "%s"
			description      = "%s"
			session_duration = %d
		}`, instanceId, name, description, sessionDuration)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "%s"
			description         = "test instance for permission set test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name             = "%s"
			description      = "%s"
			session_duration = %d
		}`, "test-acc-ps-instance", name, description, sessionDuration)
}

func testAccPermissionSetUpdate(instanceId string, name string, description string, sessionDuration int) string {
	if instanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = "%s"
			name             = "%s"
			description      = "%s"
			session_duration = %d
		}`, instanceId, name, description, sessionDuration)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "test-acc-ps-instance"
			description         = "test instance for permission set test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name             = "%s"
			description      = "%s"
			session_duration = %d
		}`, name, description, sessionDuration)
}

func testAccPermissionSetDelete(instanceId string, name string, description string, sessionDuration int) string {
	if instanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = "%s"
			name             = "%s"
			description      = "%s"
			session_duration = %d
		}`, instanceId, name, description, sessionDuration)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "test-acc-ps-instance"
			description         = "test instance for permission set test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_permission_set" "permission_set" {
			instance_id      = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			name             = "%s"
			description      = "%s"
			session_duration = %d
		}`, name, description, sessionDuration)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_identity_center_permission_set", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_identity_center_permission_set",
		F:    sweepPermissionSet,
	})
}

func sweepPermissionSet(region string) error {
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

		permissionSets, err := scpClient.Client.IamIdentityCenter.ListPermissionSets(ctx, instanceId, "", 100, 0, "")
		if err != nil {
			fmt.Printf("warning: error listing permission sets for instance %s: %v\n", instanceId, err)
			continue
		}

		for _, ps := range permissionSets.GetPermissionSets() {
			if !strings.HasPrefix(ps.Name, "test-acc-ps-") {
				continue
			}

			err = scpClient.Client.IamIdentityCenter.DeletePermissionSet(ctx, ps.Id, instanceId)
			if err != nil {
				fmt.Printf("warning: failed to delete permission set %s: %v (continuing...)\n", ps.Id, err)
				continue
			}
			fmt.Printf("deleted IAM Identity Center permission set: %s\n", ps.Id)
		}
	}

	return nil
}
