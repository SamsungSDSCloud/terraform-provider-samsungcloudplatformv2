package certificatemanager

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-certificatemanager"

//------------ CertificateManager -------------------//

type CertificateManagerDetailDataSource struct {
	Id          types.String `tfsdk:"id"`
	Certificate types.Object `tfsdk:"certificate"`
}
type CertificateManagerDataSource struct {
	// Input
	Id           types.String   `tfsdk:"id"`      // VPC Peering ID
	Size         types.Int32    `tfsdk:"size"`    // size
	Page         types.Int32    `tfsdk:"page"`    // Page
	Sort         types.String   `tfsdk:"sort"`    // Sort
	Name         types.String   `tfsdk:"name"`    // VPC Peering Name
	IsMine       types.Bool     `tfsdk:"is_mine"` // VPC Peering Name
	Cn           types.String   `tfsdk:"cn"`      // Source VPC ID
	State        []types.String `tfsdk:"state"`   // State
	Certificates []Certificate  `tfsdk:"certificates"`
}

type CertificateManagerResource struct {
	Id          types.String `tfsdk:"id"`
	CertBody    types.String `tfsdk:"cert_body"`
	CertChain   types.String `tfsdk:"cert_chain"`
	Name        types.String `tfsdk:"name"`
	PrivateKey  types.String `tfsdk:"private_key"`
	Timezone    types.String `tfsdk:"timezone"`
	Region      types.String `tfsdk:"region"`
	Tags        types.Map    `tfsdk:"tags"`
	Recipients  []types.Map  `tfsdk:"recipients"`
	Certificate types.Object `tfsdk:"certificate"`
}

type CertificateManagerSelfSignResource struct {
	Id           types.String `tfsdk:"id"`
	Cn           types.String `tfsdk:"cn"`
	Name         types.String `tfsdk:"name"`
	Organization types.String `tfsdk:"organization"`
	Timezone     types.String `tfsdk:"timezone"`
	Region       types.String `tfsdk:"region"`
	NotAfterDt   types.String `tfsdk:"not_after_dt"`
	NotBeforeDt  types.String `tfsdk:"not_before_dt"`
	Tags         types.Map    `tfsdk:"tags"`
	Recipients   []types.Map  `tfsdk:"recipients"`
	Certificate  types.Object `tfsdk:"certificate"`
}

type Certificate struct {
	CertKind    types.String `tfsdk:"cert_kind"`
	Cn          types.String `tfsdk:"cn"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	NotAfterDt  types.String `tfsdk:"not_after_dt"`
	NotBeforeDt types.String `tfsdk:"not_before_dt"`
	State       types.String `tfsdk:"state"`
}

func (m Certificate) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cert_kind":     types.StringType,
		"cn":            types.StringType,
		"name":          types.StringType,
		"not_after_dt":  types.StringType,
		"id":            types.StringType,
		"not_before_dt": types.StringType,
		"state":         types.StringType,
	}
}
