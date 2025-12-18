package certificatemanager

import (
	"context"
	"fmt"

	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	certificatemanager "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/certificatemanager/1.1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *certificatemanager.APIClient
}

func NewClient(config *scpsdk.Configuration) *Client {
	return &Client{
		Config:    config,
		sdkClient: certificatemanager.NewAPIClient(config),
	}
}

//------------ certificatemanager -------------------//

func (client *Client) GetCertificateManagerList(ctx context.Context, request CertificateManagerDataSource) (*certificatemanager.CertificateListResponse, error) {
	req := client.sdkClient.CertificatemanagerV1CertificateManagerApiAPI.ListCertificates(ctx)
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.IsMine.IsNull() {
		req = req.IsMine(request.IsMine.ValueBool())
	}
	if !request.Cn.IsNull() {
		req = req.Cn(request.Cn.ValueString())
	}
	if len(request.State) > 0 {
		var stra []*string
		for _, v := range request.State {
			stra = append(stra, v.ValueStringPointer())
		}
		req = req.State(certificatemanager.State{ArrayOfPtrString: &stra})
	}
	resp, _, err := req.Execute()
	fmt.Printf("--------------------Start CALL GetCertificateManagerList error: %v\n", err)
	fmt.Printf("--------------------Start CALL GetCertificateManagerList resp: %v\n", resp)

	if err != nil {
		return nil, err
	}
	return resp, err
}

func (client *Client) SelfSignCreateCertificateManager(ctx context.Context, request CertificateManagerSelfSignResource) (*certificatemanager.CertificateDetailResponse, error) {
	req := client.sdkClient.CertificatemanagerV1CertificateManagerApiAPI.SelfSignCert(ctx)
	tags := convertToTags(request.Tags.Elements())
	recipients, _ := convertToRecipients(request.Recipients)

	req = req.SelfSignCreateRequest(certificatemanager.SelfSignCreateRequest{
		Cn:           request.Cn.ValueString(),
		Organization: request.Organization.ValueString(),
		NotAfterDt:   request.NotAfterDt.ValueString(),
		NotBeforeDt:  request.NotBeforeDt.ValueString(),
		Region:       request.Region.ValueString(),
		Timezone:     request.Timezone.ValueString(),
		Recipients:   recipients,
		Name:         request.Name.ValueString(),
		Tags:         tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) CreateCertificateManager(ctx context.Context, request CertificateManagerResource) (*certificatemanager.CertificateCreateResponse, error) {
	req := client.sdkClient.CertificatemanagerV1CertificateManagerApiAPI.CreateCertificate(ctx)
	tags := convertToTags(request.Tags.Elements())
	recipients, _ := convertToRecipients(request.Recipients)

	req = req.CertificateCreateRequest(certificatemanager.CertificateCreateRequest{
		CertBody:   request.CertBody.ValueString(),
		CertChain:  *certificatemanager.NewNullableString(request.CertChain.ValueStringPointer()),
		PrivateKey: request.PrivateKey.ValueString(),
		Region:     request.Region.ValueString(),
		Timezone:   request.Timezone.ValueString(),
		Recipients: recipients,
		Name:       request.Name.ValueString(),
		Tags:       tags,
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetCertificateManager(ctx context.Context, certificateManagerId string) (*certificatemanager.CertificateDetailResponse, error) {
	req := client.sdkClient.CertificatemanagerV1CertificateManagerApiAPI.DetailCertificate(ctx, certificateManagerId)

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteCertificateManager(ctx context.Context, certificateId string) error {
	req := client.sdkClient.CertificatemanagerV1CertificateManagerApiAPI.DeleteCertificate(ctx, certificateId)

	_, err := req.Execute()
	return err
}

func convertToTags(elements map[string]attr.Value) []certificatemanager.Tag {
	var tags []certificatemanager.Tag
	for k, v := range elements {
		tagObject := certificatemanager.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}
		tags = append(tags, tagObject)
	}
	return tags
}

func convertToRecipients(elements []types.Map) ([]certificatemanager.Recipient, error) {
	var recipients []certificatemanager.Recipient
	for _, r := range elements {
		recipientObject := certificatemanager.Recipient{}

		for k, v := range r.Elements() {
			if k == "region" {
				recipientObject.Region = v.(types.String).ValueString()
			}
			if k == "user_id" {
				recipientObject.UserId = v.(types.String).ValueString()
			}
			if k == "user_name" {
				recipientObject.UserName = v.(types.String).ValueString()
			}
		}
		recipients = append(recipients, recipientObject)
	}
	return recipients, nil
}
