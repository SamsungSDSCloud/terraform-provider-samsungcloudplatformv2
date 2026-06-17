package servicewatch_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/servicewatch"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLogStreamResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// step1. Log Group with Log Stream Create
				Config: testAccLogGroupWithLogStreamUpdate("test-acc-log-group", 1,
					map[string]string{
						"test-acc-key": "test-acc-value"}, "testacclogstream"),
			},
			{
				// step2. Log Group Update
				Config: testAccLogGroupWithLogStreamUpdate("test-acc-log-group", 3,
					map[string]string{
						"test-acc-key": "test-acc-value"}, "testacclogstream"),
			},
			{
				// step3. Log Stream Delete
				Config: testAccLogGroupUpdate("test-acc-log-group", 3,
					map[string]string{
						"test-acc-key": "test-acc-value"}),
			},
		}, //step3. Log Group Delete
	})
}

func testAccLogGroupUpdate(name string, retentionPeriod int32, tags map[string]string) string {
	tagsJson, _ := json.Marshal(tags)
	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_servicewatch_log_group" "log_group"{
				  name = "%s"
				  retention_period = "%d"
				  tags = %s
			}
	`, name, retentionPeriod, tagsJson)
}

func testAccLogGroupWithLogStreamUpdate(log_group_name string, retentionPeriod int32, tags map[string]string, log_stream_name string) string {
	return testAccLogGroupUpdate(log_group_name, retentionPeriod, tags) +
		fmt.Sprintf(`
			resource "samsungcloudplatformv2_servicewatch_log_stream" "log_stream"{
				  log_group_id = samsungcloudplatformv2_servicewatch_log_group.log_group.id
				  name = "%s"
			}
	`, log_stream_name)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_servicewatch_log_group", &resource.Sweeper{
		Name: "samsungcloudplatformv2_servicewatch_log_group",
		F:    sweepLogGroup,
	})
}

func sweepLogGroup(region string) error {
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	r := servicewatch.LogGroupDataSources{}
	logGroups, err := scpClient.Client.ServiceWatch.GetLogGroupList(nil, r)
	if err != nil {
		return err
	}

	// 삭제 대상 ids 조회 (name: test_acc 로 시작)
	var deleteLogGroupsIds []string
	for _, logGroup := range logGroups.GetLogGroups() {
		name := logGroup.Name

		if strings.HasPrefix(name, "test-acc") {
			deleteLogGroupsIds = append(deleteLogGroupsIds, logGroup.Id)
			fmt.Println(">> [Delete] Log Group Name: " + name + ", ID: " + logGroup.Id)
		}
	}

	_, err = scpClient.Client.ServiceWatch.DeleteLogGroup(nil, deleteLogGroupsIds)
	if err != nil {
		return nil
	}

	return nil
}
