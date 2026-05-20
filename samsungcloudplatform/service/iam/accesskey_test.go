package iam_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	util "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/iam"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessKeyResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyUpdate("PERMANENT", true), // step1. access_key 생성
			},
			{
				Config: testAccAccessKeyUpdate("PERMANENT", false), // step2. access_key disable 하게 업데이트 처리
			},
			// step3. 이후 step이 끝나면 access_key에 대한 delete 메소드가 수행되며 삭제 됨
		},
	})
}

func testAccAccessKeyUpdate(accesskeyType string, is_enabled bool) string {
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_iam_access_key" "access_key" {
				access_key_type = "%s"
  				description = "test-acc access key desc"
  				is_enabled = %t
			}`, accesskeyType, is_enabled)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_iam_access_key", &resource.Sweeper{
		Name: "samsungcloudplatformv2_iam_access_key",
		F:    sweepAccessKey,
	})
}

func sweepAccessKey(region string) error {
	scpClient, err := util.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client for region %s: %v", region, err)
	}

	r := iam.AccessKeyDataSource{}
	accessKeys, err := scpClient.Client.Iam.GetAccessKeyList(nil, r)

	var deleteResourceList []string

	for _, accessKey := range accessKeys.GetAccessKeys() {
		if accessKey.Description.Get() != nil {
			description := *accessKey.Description.Get()

			if strings.HasPrefix(description, "test-acc") {
				deleteResourceList = append(deleteResourceList, accessKey.Id)
			}
		}
	}

	for _, id := range deleteResourceList {
		r := iam.AccessKeyResource{IsEnabled: types.BoolValue(false)}

		_, err := scpClient.Client.Iam.UpdateAccessKey(nil, id, r)
		if err != nil {
			return fmt.Errorf("error updating access key %s: %v", id, err)
		}

		err = scpClient.Client.Iam.DeleteAccessKey(nil, id)
		if err != nil {
			return fmt.Errorf("error deleting access key %s: %v", id, err)
		}
	}

	return nil
}
