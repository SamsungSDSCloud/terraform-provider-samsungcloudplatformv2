package dns

import (
	"context"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/library/dns/1.1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	Config    *scpsdk.Configuration
	sdkClient *dns.APIClient // 서비스의 client 를 구조체에 추가한다.
}

func NewClient(config *scpsdk.Configuration) *Client { // client 생성 함수를 추가한다.
	return &Client{
		Config:    config,
		sdkClient: dns.NewAPIClient(config),
	}
}

func (client *Client) GetPrivateDnsList(ctx context.Context, request PrivateDnsDataSource) (*dns.PrivateDnsListResponse, error) {
	req := client.sdkClient.DnsV1PrivateDnsApiAPI.ListPrivateDns(ctx)

	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Id.IsNull() {
		req = req.Id(request.Id.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.VpcId.IsNull() {
		req = req.VpcId(request.VpcId.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPrivateDns(ctx context.Context, privateDnsId string) (*dns.PrivateDnsShowResponse, error) {
	req := client.sdkClient.DnsV1PrivateDnsApiAPI.ShowPrivateDns(ctx, privateDnsId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreatePrivateDns(ctx context.Context, request PrivateDnsResource) (*dns.PrivateDnsShowResponse, error) {
	req := client.sdkClient.DnsV1PrivateDnsApiAPI.CreatePrivateDns(ctx) // 호출을 위한 구조체를 반환 받는다.

	var PrivateDnsTags []dns.Tag

	for k, v := range request.Tags.Elements() {
		tagObject := dns.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}

		PrivateDnsTags = append(PrivateDnsTags, tagObject)
	}

	connectedVpcIds := make([]string, len(request.PrivateDnsCreate.ConnectedVpcIds))

	for idx, connectedVpcId := range request.PrivateDnsCreate.ConnectedVpcIds {
		connectedVpcIds[idx] = connectedVpcId.ValueString()
	}

	privateDnsElement := dns.PrivateDnsCreateRequest{
		ConnectedVpcIds: connectedVpcIds,
		Description:     *dns.NewNullableString(request.PrivateDnsCreate.Description.ValueStringPointer()),
		Name:            request.PrivateDnsCreate.Name.ValueString(),
		Tags:            PrivateDnsTags,
	}

	req = req.PrivateDnsCreateRequest(privateDnsElement)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) ActivatePrivateDns(ctx context.Context, request PrivateDnsResource) (*dns.PrivateDnsShowResponse, error) {
	req := client.sdkClient.DnsV1PrivateDnsApiAPI.ActivatePrivateDns(ctx) // 호출을 위한 구조체를 반환 받는다.

	privateDnsElement := dns.PrivateDnsActivateRequest{
		Name: request.PrivateDnsCreate.Name.ValueString(),
	}

	req = req.PrivateDnsActivateRequest(privateDnsElement)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdatePrivateDns(ctx context.Context, privateDnsId string, request PrivateDnsResource) (*dns.PrivateDnsShowResponse, error) {
	req := client.sdkClient.DnsV1PrivateDnsApiAPI.SetPrivateDns(ctx, privateDnsId)

	privateDnsSet := request.PrivateDnsCreate

	connectedVpcIds := make([]string, len(privateDnsSet.ConnectedVpcIds))

	for idx, connectedVpcId := range privateDnsSet.ConnectedVpcIds {
		connectedVpcIds[idx] = connectedVpcId.ValueString()
	}

	req = req.PrivateDnsSetRequest(dns.PrivateDnsSetRequest{
		ConnectedVpcIds: connectedVpcIds,
		Description:     *dns.NewNullableString(privateDnsSet.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeletePrivateDns(ctx context.Context, privateDnsId string) error {
	req := client.sdkClient.DnsV1PrivateDnsApiAPI.DeletePrivateDns(ctx, privateDnsId)

	_, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return err
}

func (client *Client) GetPublicDomainNameList(ctx context.Context, request PublicDomainNameDataSource) (*dns.PublicDomainListResponse, error) {
	req := client.sdkClient.DnsV1PublicDomainNameApiAPI.ListPublicDomains(ctx)

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
	if !request.CreatedBy.IsNull() {
		req = req.CreatedBy(request.CreatedBy.ValueString())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetPublicDomainName(ctx context.Context, publicDomainId string) (*dns.PublicDomainDetailResponse, error) {
	req := client.sdkClient.DnsV1PublicDomainNameApiAPI.GetPublicDomainDetail(ctx, publicDomainId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreatePublicDomainName(ctx context.Context, request PublicDomainNameResource) (*dns.CreatePublicDomainResponse, error) {
	req := client.sdkClient.DnsV1PublicDomainNameApiAPI.CreatePublicDomain(ctx) // 호출을 위한 구조체를 반환 받는다.

	var publicDomainNameTags []dns.Tag

	for k, v := range request.Tags.Elements() {
		tagObject := dns.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}

		publicDomainNameTags = append(publicDomainNameTags, tagObject)
	}

	publicDomainNameElement := dns.CreatePublicDomainRequest{
		AddressType:             request.PublicDomainNameCreate.AddressType.ValueString(),
		AutoExtension:           *dns.NewNullableBool(request.PublicDomainNameCreate.AutoExtension.ValueBoolPointer()),
		Description:             *dns.NewNullableString(request.PublicDomainNameCreate.Description.ValueStringPointer()),
		DomesticFirstAddressEn:  *dns.NewNullableString(request.PublicDomainNameCreate.DomesticFirstAddressEn.ValueStringPointer()),
		DomesticFirstAddressKo:  *dns.NewNullableString(request.PublicDomainNameCreate.DomesticFirstAddressKo.ValueStringPointer()),
		DomesticSecondAddressEn: *dns.NewNullableString(request.PublicDomainNameCreate.DomesticSecondAddressEn.ValueStringPointer()),
		DomesticSecondAddressKo: request.PublicDomainNameCreate.DomesticSecondAddressKo.ValueString(),
		Name:                    request.PublicDomainNameCreate.Name.ValueString(),
		OverseasFirstAddress:    *dns.NewNullableString(request.PublicDomainNameCreate.OverseasFirstAddress.ValueStringPointer()),
		OverseasSecondAddress:   *dns.NewNullableString(request.PublicDomainNameCreate.OverseasSecondAddress.ValueStringPointer()),
		OverseasThirdAddress:    *dns.NewNullableString(request.PublicDomainNameCreate.OverseasThirdAddress.ValueStringPointer()),
		PostalCode:              request.PublicDomainNameCreate.PostalCode.ValueString(),
		RegisterEmail:           request.PublicDomainNameCreate.RegisterEmail.ValueString(),
		RegisterNameEn:          request.PublicDomainNameCreate.RegisterNameEn.ValueString(),
		RegisterNameKo:          request.PublicDomainNameCreate.RegisterNameKo.ValueString(),
		RegisterTelno:           request.PublicDomainNameCreate.RegisterTelno.ValueString(),
		Tags:                    publicDomainNameTags,
	}

	req = req.CreatePublicDomainRequest(publicDomainNameElement)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdatePublicDomainName(ctx context.Context, publicDomainNameId string, request PublicDomainNameResource) (*dns.PublicDomainPartialUpdateResponse, error) {
	req := client.sdkClient.DnsV1PublicDomainNameApiAPI.PutPublicDomain(ctx, publicDomainNameId)

	publicDomainNameSet := request.PublicDomainNameCreate

	req = req.PublicDomainPartialUpdateRequest(dns.PublicDomainPartialUpdateRequest{
		AutoExtension: *dns.NewNullableBool(publicDomainNameSet.AutoExtension.ValueBoolPointer()),
		Description:   *dns.NewNullableString(publicDomainNameSet.Description.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) UpdatePublicDomainNameInfomation(ctx context.Context, publicDomainNameId string, request PublicDomainNameResource) (*dns.PubblicDomainWhoisInfoUpdateResponse, error) {
	req := client.sdkClient.DnsV1PublicDomainNameApiAPI.UpdateWhoisInfoPublicDomain(ctx, publicDomainNameId)

	publicDomainNameSet := request.PublicDomainNameCreate

	req = req.PubblicDomainWhoisInfoUpdateRequest(dns.PubblicDomainWhoisInfoUpdateRequest{
		AddressType:             publicDomainNameSet.AddressType.ValueString(),
		DomesticFirstAddressEn:  publicDomainNameSet.DomesticFirstAddressEn.ValueString(),
		DomesticFirstAddressKo:  publicDomainNameSet.DomesticFirstAddressKo.ValueString(),
		DomesticSecondAddressEn: *dns.NewNullableString(publicDomainNameSet.DomesticSecondAddressEn.ValueStringPointer()),
		DomesticSecondAddressKo: publicDomainNameSet.DomesticSecondAddressKo.ValueString(),
		OverseasFirstAddress:    *dns.NewNullableString(publicDomainNameSet.OverseasFirstAddress.ValueStringPointer()),
		OverseasSecondAddress:   *dns.NewNullableString(publicDomainNameSet.OverseasSecondAddress.ValueStringPointer()),
		OverseasThirdAddress:    *dns.NewNullableString(publicDomainNameSet.OverseasThirdAddress.ValueStringPointer()),
		PostalCode:              publicDomainNameSet.PostalCode.ValueString(),
		RegisterEmail:           publicDomainNameSet.RegisterEmail.ValueString(),
		RegisterTelno:           publicDomainNameSet.RegisterTelno.ValueString(),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetHostedZoneList(ctx context.Context, request HostedZoneDataSource) (*dns.HostedZoneListResponseV1Dot1, error) {
	req := client.sdkClient.DnsV1HostedZonesApiAPI.ListHostedZone(ctx)

	if !request.Page.IsNull() {
		req = req.Page(request.Page.ValueInt32())
	}
	if !request.Size.IsNull() {
		req = req.Size(request.Size.ValueInt32())
	}
	if !request.Sort.IsNull() {
		req = req.Sort(request.Sort.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.ExactName.IsNull() {
		req = req.ExactName(request.ExactName.ValueString())
	}
	if !request.Type.IsNull() {
		req = req.Type_(request.Type.ValueString())
	}
	if !request.Email.IsNull() {
		req = req.Email(request.Email.ValueString())
	}
	if !request.Status.IsNull() {
		req = req.Status(request.Status.ValueString())
	}
	if !request.Description.IsNull() {
		req = req.Description(request.Description.ValueString())
	}
	if !request.Ttl.IsNull() {
		req = req.Ttl(request.Ttl.ValueInt32())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetHostedZone(ctx context.Context, hostedZoneId string) (*dns.HostedZoneShowResponse, error) {
	req := client.sdkClient.DnsV1HostedZonesApiAPI.ShowHostedZone(ctx, hostedZoneId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreateHostedZone(ctx context.Context, request HostedZoneResource) (*dns.HostedZoneCreateResponse, error) {
	req := client.sdkClient.DnsV1HostedZonesApiAPI.CreateHostedZone(ctx) // 호출을 위한 구조체를 반환 받는다.

	var hostedZoneTags []dns.Tag

	for k, v := range request.Tags.Elements() {
		tagObject := dns.Tag{
			Key:   k,
			Value: v.(types.String).ValueString(),
		}

		hostedZoneTags = append(hostedZoneTags, tagObject)
	}

	hostedZoneElement := dns.HostedZoneCreateRequest{
		Description: *dns.NewNullableString(request.HostedZoneCreate.Description.ValueStringPointer()),
		Email:       request.HostedZoneCreate.Email.ValueString(),
		Name:        request.HostedZoneCreate.Name.ValueString(),
		Type:        *dns.NewNullableString(request.HostedZoneCreate.Type.ValueStringPointer()),
		Tags:        hostedZoneTags,
	}

	req = req.HostedZoneCreateRequest(hostedZoneElement)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateHostedZone(ctx context.Context, hostedZoneId string, request HostedZoneResource) (*dns.HostedZoneSetResponse, error) {
	req := client.sdkClient.DnsV1HostedZonesApiAPI.SetHostedZone(ctx, hostedZoneId)

	hostedZoneSet := request.HostedZoneCreate

	req = req.HostedZoneSetRequest(dns.HostedZoneSetRequest{
		Description: *dns.NewNullableString(hostedZoneSet.Description.ValueStringPointer()),
		Email:       *dns.NewNullableString(hostedZoneSet.Email.ValueStringPointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteHostedZone(ctx context.Context, hostedZoneId string) (*dns.HostedZoneDeleteResponse, error) {
	req := client.sdkClient.DnsV1HostedZonesApiAPI.DeleteHostedZone(ctx, hostedZoneId)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) GetRecordList(ctx context.Context, request RecordDataSource) (*dns.RecordListResponse, error) {
	req := client.sdkClient.DnsV1RecordsApiAPI.ListRecords(ctx, request.HostedZoneId.ValueString())

	if !request.Limit.IsNull() {
		req = req.Limit(request.Limit.ValueInt32())
	}
	if !request.Marker.IsNull() {
		req = req.Marker(request.Marker.ValueString())
	}
	if !request.SortDir.IsNull() {
		req = req.SortDir(request.SortDir.ValueString())
	}
	if !request.SortKey.IsNull() {
		req = req.SortKey(request.SortKey.ValueString())
	}
	if !request.Name.IsNull() {
		req = req.Name(request.Name.ValueString())
	}
	if !request.ExactName.IsNull() {
		req = req.ExactName(request.ExactName.ValueString())
	}
	if !request.Type.IsNull() {
		req = req.Type_(request.Type.ValueString())
	}
	if !request.Data.IsNull() {
		req = req.Data(request.Data.ValueString())
	}
	if !request.Status.IsNull() {
		req = req.Status(request.Status.ValueString())
	}
	if !request.Description.IsNull() {
		req = req.Description(request.Description.ValueString())
	}
	if !request.Ttl.IsNull() {
		req = req.Ttl(request.Ttl.ValueInt32())
	}

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) GetRecord(ctx context.Context, hostedZoneId string, recordId string) (*dns.RecordShowResponse, error) {
	req := client.sdkClient.DnsV1RecordsApiAPI.ShowRecord(ctx, hostedZoneId, recordId)
	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) CreateRecord(ctx context.Context, hostedZoneId string, request RecordResource) (*dns.RecordCreateResponse, error) {
	req := client.sdkClient.DnsV1RecordsApiAPI.CreateRecord(ctx, hostedZoneId) // 호출을 위한 구조체를 반환 받는다.

	records := make([]interface{}, len(request.RecordCreate.Records))
	for i, r := range request.RecordCreate.Records {
		records[i] = r.ValueString() // Assuming non-null values; handle null if needed
	}

	recordElement := dns.RecordCreateRequest{
		Description: *dns.NewNullableString(request.RecordCreate.Description.ValueStringPointer()),
		Name:        request.RecordCreate.Name.ValueString(),
		Records:     records,
		Ttl:         *dns.NewNullableInt32(request.RecordCreate.Ttl.ValueInt32Pointer()),
		Type:        request.RecordCreate.Type.ValueString(),
	}

	req = req.RecordCreateRequest(recordElement)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}

func (client *Client) UpdateRecord(ctx context.Context, hostedZoneId string, recordId string, request RecordResource) (*dns.RecordSetResponse, error) {
	req := client.sdkClient.DnsV1RecordsApiAPI.SetRecord(ctx, hostedZoneId, recordId)

	recordSet := request.RecordCreate

	records := make([]interface{}, len(recordSet.Records))
	for i, r := range recordSet.Records {
		records[i] = r.ValueString() // Assuming non-null values; handle null if needed
	}

	req = req.RecordSetRequest(dns.RecordSetRequest{
		Records: records,
		Ttl:     *dns.NewNullableInt32(recordSet.Ttl.ValueInt32Pointer()),
	})

	resp, _, err := req.Execute()
	return resp, err
}

func (client *Client) DeleteRecord(ctx context.Context, hostedZoneId string, recordId string) (*dns.RecordCreateResponse, error) {
	req := client.sdkClient.DnsV1RecordsApiAPI.DeleteRecord(ctx, hostedZoneId, recordId)

	resp, _, err := req.Execute() // Execute 메서드를 호출하여 실행한다.
	return resp, err
}
