package dns

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-dns" // 해당 서비스의 서비스 타입(keystone 에 등록된 service type)을 추가한다.

type PrivateDnsDataSource struct { // Resource Group List request 모델을 참고하여 구조체를 구성한다.
	Size       types.Int32  `tfsdk:"size"`
	Page       types.Int32  `tfsdk:"page"`
	Sort       types.String `tfsdk:"sort"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	VpcId      types.String `tfsdk:"vpc_id"`
	PrivateDns []PrivateDns `tfsdk:"private_dns"`
}

type PrivateDns struct {
	AuthDnsName      types.String   `tfsdk:"auth_dns_name"`
	ConnectedVpcIds  []types.String `tfsdk:"connected_vpc_ids"`
	CreatedAt        types.String   `tfsdk:"created_at"`
	CreatedBy        types.String   `tfsdk:"created_by"`
	Description      types.String   `tfsdk:"description"`
	Id               types.String   `tfsdk:"id"`
	ModifiedAt       types.String   `tfsdk:"modified_at"`
	ModifiedBy       types.String   `tfsdk:"modified_by"`
	Name             types.String   `tfsdk:"name"`
	PoolId           types.String   `tfsdk:"pool_id"`
	PoolName         types.String   `tfsdk:"pool_name"`
	RegisteredRegion types.String   `tfsdk:"registered_region"`
	ResolverIp       types.String   `tfsdk:"resolver_ip"`
	ResolverName     types.String   `tfsdk:"resolver_name"`
	State            types.String   `tfsdk:"state"`
}

type PrivateDnsDataSourceDetail struct {
	Id         types.String `tfsdk:"id"`
	PrivateDns types.Object `tfsdk:"private_dns_detail"`
}

func (m PrivateDns) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"auth_dns_name": types.StringType,
		"connected_vpc_ids": types.ListType{
			ElemType: types.StringType,
		},
		"created_at":        types.StringType,
		"created_by":        types.StringType,
		"description":       types.StringType,
		"id":                types.StringType,
		"modified_at":       types.StringType,
		"modified_by":       types.StringType,
		"name":              types.StringType,
		"pool_id":           types.StringType,
		"pool_name":         types.StringType,
		"registered_region": types.StringType,
		"resolver_ip":       types.StringType,
		"resolver_name":     types.StringType,
		"state":             types.StringType,
	}
}

type PrivateDnsResource struct {
	Id               types.String     `tfsdk:"id"`
	PrivateDns       types.Object     `tfsdk:"private_dns"`
	PrivateDnsCreate PrivateDnsCreate `tfsdk:"private_dns_create"`
	Tags             types.Map        `tfsdk:"tags"`
}

type PrivateDnsCreate struct {
	ConnectedVpcIds []types.String `tfsdk:"connected_vpc_ids"`
	Description     types.String   `tfsdk:"description"`
	Name            types.String   `tfsdk:"name"`
}

type PublicDomainNameDataSource struct { // Resource Group List request 모델을 참고하여 구조체를 구성한다.
	Size              types.Int32        `tfsdk:"size"`
	Page              types.Int32        `tfsdk:"page"`
	Sort              types.String       `tfsdk:"sort"`
	Name              types.String       `tfsdk:"name"`
	CreatedBy         types.String       `tfsdk:"created_by"`
	PublicDomainNames []PublicDomainName `tfsdk:"public_domain_names"`
}

type PublicDomainName struct {
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	ExpiredDate types.String `tfsdk:"expired_date"`
	Id          types.String `tfsdk:"id"`
	ModifiedAt  types.String `tfsdk:"modified_at"`
	ModifiedBy  types.String `tfsdk:"modified_by"`
	Name        types.String `tfsdk:"name"`
	StartDate   types.String `tfsdk:"start_date"`
	Status      types.String `tfsdk:"status"`
}

type PublicDomainNameDataSourceDetail struct {
	Id                     types.String `tfsdk:"id"`
	PublicDomainNameDetail types.Object `tfsdk:"public_domain_name_detail"`
}

type PublicDomainNameDetail struct {
	AddressType             types.String `tfsdk:"address_type"`
	AutoExtension           types.Bool   `tfsdk:"auto_extension"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	Description             types.String `tfsdk:"description"`
	DomesticAddressEn       types.String `tfsdk:"domestic_address_en"`
	DomesticAddressKo       types.String `tfsdk:"domestic_address_ko"`
	DomesticFirstAddressEn  types.String `tfsdk:"domestic_first_address_en"`
	DomesticFirstAddressKo  types.String `tfsdk:"domestic_first_address_ko"`
	DomesticSecondAddressEn types.String `tfsdk:"domestic_second_address_en"`
	DomesticSecondAddressKo types.String `tfsdk:"domestic_second_address_ko"`
	ExpiredDate             types.String `tfsdk:"expired_date"`
	Id                      types.String `tfsdk:"id"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	Name                    types.String `tfsdk:"name"`
	OverseasAddress         types.String `tfsdk:"overseas_address"`
	OverseasFirstAddress    types.String `tfsdk:"overseas_first_address"`
	OverseasSecondAddress   types.String `tfsdk:"overseas_second_address"`
	OverseasThirdAddress    types.String `tfsdk:"overseas_third_address"`
	PostalCode              types.String `tfsdk:"postal_code"`
	RegisterEmail           types.String `tfsdk:"register_email"`
	RegisterNameEn          types.String `tfsdk:"register_name_en"`
	RegisterNameKo          types.String `tfsdk:"register_name_ko"`
	RegisterTelno           types.String `tfsdk:"register_telno"`
	StartDate               types.String `tfsdk:"start_date"`
	Status                  types.String `tfsdk:"status"`
}

func (m PublicDomainNameDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address_type":               types.StringType,
		"auto_extension":             types.BoolType,
		"created_at":                 types.StringType,
		"created_by":                 types.StringType,
		"description":                types.StringType,
		"domestic_address_en":        types.StringType,
		"domestic_address_ko":        types.StringType,
		"domestic_first_address_en":  types.StringType,
		"domestic_first_address_ko":  types.StringType,
		"domestic_second_address_en": types.StringType,
		"domestic_second_address_ko": types.StringType,
		"expired_date":               types.StringType,
		"id":                         types.StringType,
		"modified_at":                types.StringType,
		"modified_by":                types.StringType,
		"name":                       types.StringType,
		"overseas_address":           types.StringType,
		"overseas_first_address":     types.StringType,
		"overseas_second_address":    types.StringType,
		"overseas_third_address":     types.StringType,
		"postal_code":                types.StringType,
		"register_email":             types.StringType,
		"register_name_en":           types.StringType,
		"register_name_ko":           types.StringType,
		"register_telno":             types.StringType,
		"start_date":                 types.StringType,
		"status":                     types.StringType,
	}
}

type PublicDomainNameResource struct {
	Id                     types.String           `tfsdk:"id"`
	PublicDomainName       types.Object           `tfsdk:"public_domain_name"`
	PublicDomainNameCreate PublicDomainNameCreate `tfsdk:"public_domain_name_create"`
	Tags                   types.Map              `tfsdk:"tags"`
}

type PublicDomainNameCreate struct {
	AddressType             types.String `tfsdk:"address_type"`
	AutoExtension           types.Bool   `tfsdk:"auto_extension"`
	Description             types.String `tfsdk:"description"`
	DomesticFirstAddressEn  types.String `tfsdk:"domestic_first_address_en"`
	DomesticFirstAddressKo  types.String `tfsdk:"domestic_first_address_ko"`
	DomesticSecondAddressEn types.String `tfsdk:"domestic_second_address_en"`
	DomesticSecondAddressKo types.String `tfsdk:"domestic_second_address_ko"`
	Name                    types.String `tfsdk:"name"`
	OverseasFirstAddress    types.String `tfsdk:"overseas_first_address"`
	OverseasSecondAddress   types.String `tfsdk:"overseas_second_address"`
	OverseasThirdAddress    types.String `tfsdk:"overseas_third_address"`
	PostalCode              types.String `tfsdk:"postal_code"`
	RegisterEmail           types.String `tfsdk:"register_email"`
	RegisterNameEn          types.String `tfsdk:"register_name_en"`
	RegisterNameKo          types.String `tfsdk:"register_name_ko"`
	RegisterTelno           types.String `tfsdk:"register_telno"`
}

type HostedZoneDataSource struct { // Resource Group List request 모델을 참고하여 구조체를 구성한다.
	Page        types.Int32  `tfsdk:"page"`
	Size        types.Int32  `tfsdk:"size"`
	Sort        types.String `tfsdk:"sort"`
	Name        types.String `tfsdk:"name"`
	ExactName   types.String `tfsdk:"exact_name"`
	Type        types.String `tfsdk:"type"`
	Email       types.String `tfsdk:"email"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
	Ttl         types.Int32  `tfsdk:"ttl"`
	HostedZones []HostedZone `tfsdk:"hosted_zones"`
}

type HostedZone struct {
	Action         types.String   `tfsdk:"action"`
	Attributes     *Attributes    `tfsdk:"attributes"`
	CreatedAt      types.String   `tfsdk:"created_at"`
	Description    types.String   `tfsdk:"description"`
	Email          types.String   `tfsdk:"email"`
	HostedZoneType types.String   `tfsdk:"hosted_zone_type"`
	Id             types.String   `tfsdk:"id"`
	Links          *Links         `tfsdk:"links"`
	Masters        []types.String `tfsdk:"masters"`
	Name           types.String   `tfsdk:"name"`
	PoolId         types.String   `tfsdk:"pool_id"`
	PrivateDnsId   types.String   `tfsdk:"private_dns_id"`
	PrivateDnsName types.String   `tfsdk:"private_dns_name"`
	ProjectId      types.String   `tfsdk:"project_id"`
	Serial         types.Int32    `tfsdk:"serial"`
	Shared         types.Bool     `tfsdk:"shared"`
	Status         types.String   `tfsdk:"status"`
	TransferredAt  types.String   `tfsdk:"transferred_at"`
	Ttl            types.Int32    `tfsdk:"ttl"`
	Type           types.String   `tfsdk:"type"`
	UpdatedAt      types.String   `tfsdk:"updated_at"`
	Version        types.Int32    `tfsdk:"version"`
}

type Attributes struct {
	ServiceTier types.String `tfsdk:"service_tier"`
}

type Links struct {
	Self types.String `tfsdk:"self"`
}

type HostedZoneDataSourceDetail struct {
	Id               types.String `tfsdk:"id"`
	HostedZoneDetail types.Object `tfsdk:"hosted_zone_detail"`
}

func (m HostedZone) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action": types.StringType,
		"attributes": types.ObjectType{AttrTypes: map[string]attr.Type{
			"service_tier": types.StringType,
		}},
		"created_at":       types.StringType,
		"description":      types.StringType,
		"email":            types.StringType,
		"hosted_zone_type": types.StringType,
		"links": types.ObjectType{AttrTypes: map[string]attr.Type{
			"self": types.StringType,
		}},
		"id":               types.StringType,
		"masters":          types.ListType{ElemType: types.StringType},
		"name":             types.StringType,
		"pool_id":          types.StringType,
		"private_dns_id":   types.StringType,
		"private_dns_name": types.StringType,
		"project_id":       types.StringType,
		"serial":           types.Int32Type,
		"shared":           types.BoolType,
		"status":           types.StringType,
		"transferred_at":   types.StringType,
		"ttl":              types.Int32Type,
		"type":             types.StringType,
		"updated_at":       types.StringType,
		"version":          types.Int32Type,
	}
}

type HostedZoneResource struct {
	Id               types.String     `tfsdk:"id"`
	Zone             types.Object     `tfsdk:"zone"`
	HostedZoneCreate HostedZoneCreate `tfsdk:"hosted_zone_create"`
	Tags             types.Map        `tfsdk:"tags"`
}

type HostedZoneCreate struct {
	Description  types.String `tfsdk:"description"`
	Email        types.String `tfsdk:"email"`
	Name         types.String `tfsdk:"name"`
	PrivateDnsId types.String `tfsdk:"private_dns_id"`
	Type         types.String `tfsdk:"type"`
}

type RecordDataSource struct { // Resource Group List request 모델을 참고하여 구조체를 구성한다.
	Limit        types.Int32  `tfsdk:"limit"`
	Marker       types.String `tfsdk:"marker"`
	SortDir      types.String `tfsdk:"sort_dir"`
	SortKey      types.String `tfsdk:"sort_key"`
	Name         types.String `tfsdk:"name"`
	ExactName    types.String `tfsdk:"exact_name"`
	Type         types.String `tfsdk:"type"`
	Data         types.String `tfsdk:"data"`
	Status       types.String `tfsdk:"status"`
	Description  types.String `tfsdk:"description"`
	Ttl          types.Int32  `tfsdk:"ttl"`
	HostedZoneId types.String `tfsdk:"hosted_zone_id"`
	Records      []Record     `tfsdk:"records"`
}

type Record struct {
	Action      types.String   `tfsdk:"action"`
	CreatedAt   types.String   `tfsdk:"created_at"`
	Description types.String   `tfsdk:"description"`
	Id          types.String   `tfsdk:"id"`
	Links       *Links         `tfsdk:"links"`
	Name        types.String   `tfsdk:"name"`
	ProjectId   types.String   `tfsdk:"project_id"`
	Records     []types.String `tfsdk:"records"`
	Status      types.String   `tfsdk:"status"`
	Ttl         types.Int32    `tfsdk:"ttl"`
	Type        types.String   `tfsdk:"type"`
	UpdatedAt   types.String   `tfsdk:"updated_at"`
	Version     types.Int32    `tfsdk:"version"`
	ZoneId      types.String   `tfsdk:"zone_id"`
	ZoneName    types.String   `tfsdk:"zone_name"`
}

type RecordDataSourceDetail struct {
	Id           types.String `tfsdk:"id"`
	HostedZoneId types.String `tfsdk:"hosted_zone_id"`
	RecordDetail types.Object `tfsdk:"record_detail"`
}

type RecordDetail struct {
	Action      types.String   `tfsdk:"action"`
	CreatedAt   types.String   `tfsdk:"created_at"`
	Description types.String   `tfsdk:"description"`
	Id          types.String   `tfsdk:"id"`
	Links       *Links         `tfsdk:"links"`
	Name        types.String   `tfsdk:"name"`
	ProjectId   types.String   `tfsdk:"project_id"`
	Records     []types.String `tfsdk:"records"`
	Status      types.String   `tfsdk:"status"`
	Ttl         types.Int32    `tfsdk:"ttl"`
	Type        types.String   `tfsdk:"type"`
	UpdatedAt   types.String   `tfsdk:"updated_at"`
	Version     types.Int32    `tfsdk:"version"`
	ZoneId      types.String   `tfsdk:"zone_id"`
	ZoneName    types.String   `tfsdk:"zone_name"`
}

func (m RecordDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action":      types.StringType,
		"created_at":  types.StringType,
		"description": types.StringType,
		"id":          types.StringType,
		"links": types.ObjectType{AttrTypes: map[string]attr.Type{
			"self": types.StringType,
		}},
		"name":       types.StringType,
		"project_id": types.StringType,
		"records":    types.ListType{ElemType: types.StringType},
		"status":     types.StringType,
		"ttl":        types.Int32Type,
		"type":       types.StringType,
		"updated_at": types.StringType,
		"version":    types.Int32Type,
		"zone_id":    types.StringType,
		"zone_name":  types.StringType,
	}
}

type RecordResource struct {
	Id           types.String `tfsdk:"id"`
	HostedZoneId types.String `tfsdk:"hosted_zone_id"`
	Record       types.Object `tfsdk:"record"`
	RecordCreate RecordCreate `tfsdk:"record_create"`
}

type RecordCreate struct {
	Description types.String   `tfsdk:"description"`
	Name        types.String   `tfsdk:"name"`
	Records     []types.String `tfsdk:"records"`
	Ttl         types.Int32    `tfsdk:"ttl"`
	Type        types.String   `tfsdk:"type"`
}
