package servicewatch_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDashboardResourceTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"samsungcloudplatformv2": providerserver.NewProtocol6WithError(samsungcloudplatform.NewProvider("test")),
		},
		Steps: []resource.TestStep{
			{
				// step1. Dashboard Create
				Config: testAccDashboardCreate("test-acc-dashboard")},
			{
				// step2. Dashboard Update
				Config: testAccDashboardUpdate("test-acc-dashboard",
					60,
					"AVG",
					"Test ACC",
					"CPU Usage",
					"Virtual Server"),
			},
			// step3. Dashboard Delete
		},
	})
}

func testAccDashboardCreate(name string) string {

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_servicewatch_dashboard" "dashboard"{
				  name = "%s"
			}
	`, name)
}

func testAccDashboardUpdate(name string,
	period int32,
	statistic string,
	title string,
	metricName string,
	namespaceName string) string {

	return fmt.Sprintf(`
			resource "samsungcloudplatformv2_servicewatch_dashboard" "dashboard"{
				  name = "%s"
				  widgets = [{
					width = "1"
					height = "1"
					order = "1"
					type = "metric"
					properties = {
						period = "%d"
						stacked = "false"
						statistic_type = "%s"
						title = "%s"
						view = "line"
						metrics = [{
							color = ""
							display_name = "Test ACC"
							name = "%s"
							namespace_name = "%s"
							period = "%d"
							statistic_type = "%s"
							dimensions = [{key="resource_id", value="d5b49100-e3e3-4d10-b2e9-9da68aed7747"}]
						}]
					}
				  }]
			}
	`, name, period, statistic, title, metricName, namespaceName, period, statistic)
}

func init() {
	resource.AddTestSweepers("samsungcloudplatformv2_servicewatch_dashboard", &resource.Sweeper{
		Name: "samsungcloudplatformv2_servicewatch_dashboard",
		F:    sweepDashboard,
	})
}

func sweepDashboard(region string) error {
	scpClient, err := SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	r := servicewatch.DashboardDataSources{}
	dashboards, err := scpClient.Client.ServiceWatch.GetDashboardList(r)
	if err != nil {
		return err
	}

	// 삭제 대상 ids 조회 (name: test_acc 로 시작)
	var deleteDashboardIds []string
	for _, dashboard := range dashboards.GetDashboards() {
		name := dashboard.Name

		if strings.HasPrefix(name, "test-acc") {
			deleteDashboardIds = append(deleteDashboardIds, dashboard.Id)
			fmt.Println(">> [Delete] Dashboard Name: " + name + ", ID: " + dashboard.Id)

		}
	}

	_, err = scpClient.Client.ServiceWatch.DeleteDashboard(nil, deleteDashboardIds)
	if err != nil {
		return nil
	}

	return nil
}
