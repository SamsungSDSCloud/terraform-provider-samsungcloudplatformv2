package iam_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	util "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/iam"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserCreate("test-acc-user", "test-acc user desc",
					map[string]string{
						"test-acc-key": "test-acc-value"}, 2),
			},
			{
				Config: testAccUserUpdate("test-acc user desc change", 3),
			},
		},
	})
}

func testAccUserCreate(name string, description string, tags map[string]string, passwordReuseCount int) string {
	tagsJson, _ := json.Marshal(tags)

	return fmt.Sprintf(`
			%s

			resource "samsungcloudplatformv2_iam_user" "user"{
				  account_id = data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys[0].account_id
				  description = "%s"
				  password = "U2Ftc3VuZ3Nkc0Nsb3VkITIz"
				  user_name = "%s"
				  tags = %s
				  policy_ids = ["9bfa3a0668a146b6a90320743be400eb"]
                  password_reuse_count = %d
			}
	`, getAccountIdForUser(), description, name, tagsJson, passwordReuseCount)
}
func testAccUserUpdate(description string, passwordReuseCount int) string {
	return fmt.Sprintf(`
			%s

			resource "samsungcloudplatformv2_iam_user" "user"{
                  account_id = data.samsungcloudplatformv2_iam_access_keys.access_keys.access_keys[0].account_id
				  description = "%s"
                  password_reuse_count = %d
			}
	`, getAccountIdForUser(), description, passwordReuseCount)
}

func getAccountIdForUser() string {
	return fmt.Sprintf(`
			data "samsungcloudplatformv2_iam_access_keys" "access_keys" {
				limit = 1
			}
	`)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_user", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_user",
		F:    sweepUser,
	})

}

func sweepUser(region string) error {
	scpClient, err := util.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	accountId := getAccountIdUsingClient(err, scpClient)

	r := iam.UserDataSource{Size: types.Int32Value(1000)}
	users, err := scpClient.Client.Iam.GetUsers(nil, accountId, r)

	var deleteResourceList []string

	for _, user := range users.GetUsers() {
		if user.Description.Get() != nil {
			description := *user.Description.Get()

			if strings.HasPrefix(description, "test-acc") {
				deleteResourceList = append(deleteResourceList, user.Id)
			}
		}

	}

	for _, id := range deleteResourceList {
		err = scpClient.Client.Iam.DeleteUser(nil, accountId, id)
		if err != nil {
			return fmt.Errorf("error deleting user %s: %v", id, err)
		}
	}

	return nil
}

func getAccountIdUsingClient(err error, scpClient client.Instance) string {
	r := iam.AccessKeyDataSource{}
	accessKeys, err := scpClient.Client.Iam.GetAccessKeyList(nil, r)
	var accountId string
	for _, accessKey := range accessKeys.GetAccessKeys() {
		accountId = accessKey.AccountId
		break
	}

	return accountId
}
