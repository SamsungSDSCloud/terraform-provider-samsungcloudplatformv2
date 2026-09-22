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

func TestAccIamIdentityCenterUserResourceTest(t *testing.T) {
	existingInstanceId := getExistingInstanceId()
	testUserId := fmt.Sprintf("test-acc-u-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserCreate(
					existingInstanceId,
					testUserId,
					"Test User Name",
					"test1234@example.com",
					"P@ssw0rd!2024Secure",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("samsungcloudplatformv2_iam_identity_center_user.user", "instance_id"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_user.user", "user_id", testUserId),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_user.user", "name", "Test User Name"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_user.user", "email", "test1234@example.com"),
				),
			},
			{
				Config: testAccUserUpdate(
					existingInstanceId,
					testUserId,
					"Updated User Name",
					"updated@example.com",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_user.user", "name", "Updated User Name"),
					resource.TestCheckResourceAttr("samsungcloudplatformv2_iam_identity_center_user.user", "email", "updated@example.com"),
				),
			},
			{
				Config: testAccUserDelete(
					existingInstanceId,
					testUserId,
					"Updated User Name",
					"updated@example.com",
				),
			},
		},
	})
}

func testAccUserCreate(instanceId string, userId string, name string, email string, password string) string {
	if instanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = "%s"
			user_id     = "%s"
			name        = "%s"
			email       = "%s"
			password    = "%s"
		}`, instanceId, userId, name, email, password)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "test-acc-u-instance"
			description         = "test instance for user test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			user_id     = "%s"
			name        = "%s"
			email       = "%s"
			password    = "%s"
		}`, userId, name, email, password)
}

func testAccUserUpdate(instanceId string, userId string, name string, email string) string {
	if instanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = "%s"
			user_id     = "%s"
			name        = "%s"
			email       = "%s"
		}`, instanceId, userId, name, email)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "test-acc-u-instance"
			description         = "test instance for user test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			user_id     = "%s"
			name        = "%s"
			email       = "%s"
		}`, userId, name, email)
}

func testAccUserDelete(instanceId string, userId string, name string, email string) string {
	if instanceId != "" {
		return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = "%s"
			user_id     = "%s"
			name        = "%s"
			email       = "%s"
		}`, instanceId, userId, name, email)
	}
	return fmt.Sprintf(`
		resource "samsungcloudplatformv2_iam_identity_center_instance" "instance" {
			name                = "test-acc-u-instance"
			description         = "test instance for user test"
			identity_store_type = "IDENTITY_CENTER_DIRECTORY"
		}

		resource "samsungcloudplatformv2_iam_identity_center_user" "user" {
			instance_id = samsungcloudplatformv2_iam_identity_center_instance.instance.id
			user_id     = "%s"
			name        = "%s"
			email       = "%s"
		}`, userId, name, email)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_identity_center_user", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_identity_center_user",
		F:    sweepUser,
	})
}

func sweepUser(region string) error {
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

		users, err := scpClient.Client.IamIdentityCenter.ListUsers(ctx, instanceId, "", 100, 0, "", "", "")
		if err != nil {
			fmt.Printf("warning: error listing users for instance %s: %v\n", instanceId, err)
			continue
		}

		for _, u := range users.GetUsers() {
			if strings.HasPrefix(u.UserId, "test-acc-u-") {
				_, _, err = scpClient.Client.IamIdentityCenter.DeleteUser(ctx, u.Id, instanceId)
				if err != nil {
					fmt.Printf("warning: failed to delete user %s: %v\n", u.Id, err)
				}
			}
		}
	}

	return nil
}
