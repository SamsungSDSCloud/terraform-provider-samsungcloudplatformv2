package firewall

import (
	"context"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/firewall/1.0"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	config    *scpsdk.Configuration
	sdkClient *scpfirewall.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: scpfirewall.NewAPIClient(config),
	}
}

//------------------- Firewall -------------------//

func (client *Client) GetFirewallList(page types.Int32, size types.Int32, sort types.String, name types.String,
	vpcName types.String, productType types.List, state types.List) (*scpfirewall.FirewallListResponse, error) {

	ctx := context.Background()
	req := client.sdkClient.FirewallV1FirewallApiAPI.ListFirewalls(ctx)

	if !page.IsNull() {
		req = req.Page(page.ValueInt32())
	}
	if !size.IsNull() {
		req = req.Size(size.ValueInt32())
	}
	if !sort.IsNull() {
		req = req.Sort(sort.ValueString())
	}
	if !name.IsNull() {
		req = req.Name(name.ValueString())
	}
	if !vpcName.IsNull() {
		req = req.VpcName(vpcName.ValueString())
	}
	if !productType.IsNull() {
		var productTypes []string
		productType.ElementsAs(ctx, &productTypes, false)
		reqProductTypeList := make([]scpfirewall.FirewallProductType, 0, len(productTypes))
		for _, productType := range productTypes {
			reqProductTypeList = append(reqProductTypeList, scpfirewall.FirewallProductType(productType))
		}
		req = req.ProductType(reqProductTypeList)
	}
	if !state.IsNull() {
		var states []string
		state.ElementsAs(ctx, &states, false)
		reqStateList := make([]scpfirewall.FirewallState, 0, len(states))
		for _, state := range states {
			reqStateList = append(reqStateList, scpfirewall.FirewallState(state))
		}
		req = req.State(reqStateList)
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetFirewall(firewallId string) (*scpfirewall.FirewallShowResponse, error) {

	ctx := context.Background()
	req := client.sdkClient.FirewallV1FirewallApiAPI.ShowFirewall(ctx, firewallId)

	resp, _, err := req.Execute()
	return resp, err
}

//------------------- Firewall Rule -------------------//

func (client *Client) GetFirewallRule(firewallRuleId string) (*scpfirewall.FirewallRuleShowResponse, error) {

	ctx := context.Background()
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.ShowFirewallRule(ctx, firewallRuleId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetFirewallRuleList(page types.Int32, size types.Int32, sort types.String, firewallId types.String,
	srcIp types.String, dstIp types.String, description types.String, state types.List, status types.String, fetchAll types.Bool) (*scpfirewall.FirewallRuleListResponse, error) {

	ctx := context.Background()
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.ListFirewallRules(ctx)

	if !page.IsNull() {
		req = req.Page(page.ValueInt32())
	}
	if !size.IsNull() {
		req = req.Size(size.ValueInt32())
	}
	if !sort.IsNull() {
		req = req.Sort(sort.ValueString())
	}
	if !firewallId.IsNull() {
		req = req.FirewallId(firewallId.ValueString())
	}
	if !srcIp.IsNull() {
		req = req.SrcIp(srcIp.ValueString())
	}
	if !dstIp.IsNull() {
		req = req.DstIp(dstIp.ValueString())
	}
	if !description.IsNull() {
		req = req.Description(description.ValueString())
	}
	if !state.IsNull() {
		var states []string
		state.ElementsAs(ctx, &states, false)
		reqStateList := make([]scpfirewall.FirewallRuleState, 0, len(states))
		for _, s := range states {
			reqStateList = append(reqStateList, scpfirewall.FirewallRuleState(s))
		}
		req = req.State(reqStateList)
	}
	if !status.IsNull() {
		req = req.Status(scpfirewall.FirewallStatusType(status.ValueString()))
	}
	if !fetchAll.IsNull() {
		req = req.FetchAll(fetchAll.ValueBool())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateFirewallRule(ctx context.Context, request FirewallRuleResource) (*scpfirewall.FirewallRuleShowResponse, error) {
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.CreateFirewallRule(ctx)

	firewallRuleElement := scpfirewall.FirewallRuleCreateRequest{
		SourceAddress:      request.FirewallRuleCreate.SourceAddress,
		DestinationAddress: request.FirewallRuleCreate.DestinationAddress,
		Service:            convertFirewallPorts(request.FirewallRuleCreate.Service),
		Action:             scpfirewall.FirewallRuleAction(request.FirewallRuleCreate.Action.ValueString()),
		Direction:          scpfirewall.FirewallRuleDirection(request.FirewallRuleCreate.Direction.ValueString()),
		OrderRuleId:        request.FirewallRuleCreate.OrderRuleId.ValueStringPointer(),
		OrderDirection:     convertOrderDirection(request.FirewallRuleCreate.OrderDirection.ValueStringPointer()),
		Status:             scpfirewall.FirewallStatusType(request.FirewallRuleCreate.Status.ValueString()),
		Description:        *scpfirewall.NewNullableString(request.FirewallRuleCreate.Description.ValueStringPointer()),
	}

	req = req.FirewallRuleCreateSingleRequest(scpfirewall.FirewallRuleCreateSingleRequest{
		FirewallId:   request.FirewallId.ValueString(),
		FirewallRule: firewallRuleElement,
	})

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, fmt.Errorf("error creating firewall rule: %w", err)
	}
	defer httpResp.Body.Close()

	return resp, nil
}

func (client *Client) UpdateFirewallRule(ctx context.Context, firewallRuleId string, request FirewallRuleResource) (*scpfirewall.FirewallRuleShowResponse, error) {
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.SetFirewallRule(ctx, firewallRuleId)

	req = req.FirewallRuleSetRequest(scpfirewall.FirewallRuleSetRequest{
		SourceAddress:      request.FirewallRuleCreate.SourceAddress,
		DestinationAddress: request.FirewallRuleCreate.DestinationAddress,
		Service:            convertFirewallPorts(request.FirewallRuleCreate.Service),
		Action:             scpfirewall.FirewallRuleAction(request.FirewallRuleCreate.Action.ValueString()),
		Direction:          scpfirewall.FirewallRuleDirection(request.FirewallRuleCreate.Direction.ValueString()),
		Description:        *scpfirewall.NewNullableString(request.FirewallRuleCreate.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteFirewallRule(ctx context.Context, firewallRuleId string) error {
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.DeleteFirewallRule(ctx, firewallRuleId)

	_, err := req.Execute()
	return err
}
