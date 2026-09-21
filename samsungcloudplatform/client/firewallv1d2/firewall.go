package firewall

import (
	"context"
	"fmt"
	"net/http"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/firewall/1.2"
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

//------------------- Firewall Rule -------------------//

func (client *Client) CreateFirewallRule(ctx context.Context, request FirewallRuleResource) (*scpfirewall.FirewallRuleShowResponseV1Dot2, error) {
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.CreateFirewallRule(ctx)

	frc := request.FirewallRuleCreate
	if frc == nil {
		return nil, fmt.Errorf("firewall_rule_create block is required")
	}

	firewallRuleReq := scpfirewall.FirewallRuleCreateRequest{
		SourceAddress:      frc.SourceAddress,
		DestinationAddress: frc.DestinationAddress,
		Service:            convertFirewallPorts(frc.Service),
		Action:             scpfirewall.FirewallRuleAction(frc.Action.ValueString()),
		Direction:          scpfirewall.FirewallRuleDirection(frc.Direction.ValueString()),
		OrderRuleId:        frc.OrderRuleId.ValueStringPointer(),
		OrderDirection:     convertOrderDirection(frc.OrderDirection.ValueStringPointer()),
		Status:             scpfirewall.FirewallStatusType(frc.Status.ValueString()),
		Description:        *scpfirewall.NewNullableString(frc.Description.ValueStringPointer()),
	}

	req = req.FirewallRuleCreateSingleRequest(scpfirewall.FirewallRuleCreateSingleRequest{
		FirewallId:   request.FirewallId.ValueString(),
		FirewallRule: firewallRuleReq,
	})

	dataResp, httpResp, err := req.Execute()
	if httpResp != nil {
		defer httpResp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("error creating firewall rule v1d2: %w", err)
	}

	return dataResp, nil
}

func (client *Client) GetFirewallRule(id string) (*scpfirewall.FirewallRuleShowResponseV1Dot1, error) {
	ctx := context.Background()
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.ShowFirewallRule(ctx, id)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteFirewallRule(ctx context.Context, id string) (*http.Response, error) {
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.DeleteFirewallRule(ctx, id)

	httpResp, err := req.Execute()
	return httpResp, err
}

func (client *Client) GetFirewallRuleList(page types.Int32, size types.Int32, sort types.String, firewallId types.String,
	srcIp types.String, dstIp types.String, description types.String, state types.List, status types.String) (*scpfirewall.FirewallRuleListResponseV1Dot1, error) {

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

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdateFirewallRule(ctx context.Context, firewallRuleId string, request FirewallRuleResource) (*scpfirewall.FirewallRuleShowResponseV1Dot2, error) {
	req := client.sdkClient.FirewallV1FirewallRulesApiAPI.SetFirewallRule(ctx, firewallRuleId)

	frc := request.FirewallRuleCreate
	if frc == nil {
		return nil, fmt.Errorf("firewall_rule_create block is required")
	}

	req = req.FirewallRuleUpdateSingleRequest(scpfirewall.FirewallRuleUpdateSingleRequest{
		FirewallRule: scpfirewall.FirewallRuleCreateRequest{
			Action:             scpfirewall.FirewallRuleAction(frc.Action.ValueString()),
			Description:        *scpfirewall.NewNullableString(frc.Description.ValueStringPointer()),
			DestinationAddress: frc.DestinationAddress,
			Direction:          scpfirewall.FirewallRuleDirection(frc.Direction.ValueString()),
			OrderDirection:     convertOrderDirection(frc.OrderDirection.ValueStringPointer()),
			OrderRuleId:        frc.OrderRuleId.ValueStringPointer(),
			Service:            convertFirewallPorts(frc.Service),
			SourceAddress:      frc.SourceAddress,
			Status:             scpfirewall.FirewallStatusType(frc.Status.ValueString()),
		},
	})

	resp, _, err := req.Execute()
	return resp, err
}
