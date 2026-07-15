package configinspection

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceType = "scp-configinspection"

// #region List of Config Inspect Object
type ConfigInspectionDataSources struct {
	// Input
	WithCount            types.String   `tfsdk:"with_count"`             // With count
	Limit                types.Int32    `tfsdk:"limit"`                  // Limit
	Marker               types.String   `tfsdk:"marker"`                 // Marker
	Sort                 types.String   `tfsdk:"sort"`                   // Sort
	IsMine               types.Bool     `tfsdk:"is_mine"`                // My Config Inspection
	DiagnosisID          types.String   `tfsdk:"diagnosis_id"`           // Id of diagnosis
	DiagnosisName        types.String   `tfsdk:"diagnosis_name"`         // Name of diagnosis
	CSPType              types.String   `tfsdk:"csp_type"`               // Type of cloud service provider
	DiagnosisAccountID   types.String   `tfsdk:"diagnosis_account_id"`   // Id of diagnosis
	RecentDiagnosisState []types.String `tfsdk:"recent_diagnosis_state"` // Recent diagnosis state
	StartDate            types.String   `tfsdk:"start_date"`             // Start date
	EndDate              types.String   `tfsdk:"end_date"`               // End date

	// Output
	TotalCount       types.Int32       `tfsdk:"total_count"`       // Total count
	Links            []Link            `tfsdk:"links"`             // Links
	SummaryResponses []SummaryResponse `tfsdk:"summary_responses"` //
}

type SummaryResponse struct {
	CreatedAt            types.String `tfsdk:"created_at"`             // Created date
	CspType              types.String `tfsdk:"csp_type"`               // Type of cloud service provider
	DiagnosisAccountId   types.String `tfsdk:"diagnosis_account_id"`   // Id of diagnosis
	DiagnosisCheckType   types.String `tfsdk:"diagnosis_check_type"`   // Check type of diagnosis
	DiagnosisId          types.String `tfsdk:"diagnosis_id"`           // Id of diagnosis
	DiagnosisName        types.String `tfsdk:"diagnosis_name"`         // Name of diagnosis
	DiagnosisType        types.String `tfsdk:"diagnosis_type"`         // diagnosis Type
	PlanType             types.String `tfsdk:"plan_type"`              // Plan Type
	ErrorState           types.String `tfsdk:"error_state"`            // State Error Type
	RecentDiagnosisAt    types.String `tfsdk:"recent_diagnosis_at"`    // Last diagnosis time
	RecentDiagnosisState types.String `tfsdk:"recent_diagnosis_state"` // Last diagnosis state
}

type Link struct {
	Href types.String `tfsdk:"href"`
	Rel  types.String `tfsdk:"rel"`
}

// #endregion

// #region Detail of Config Inspect Object
type ConfigInspectionDataSource struct {
	// Input
	DiagnosisID types.String `tfsdk:"diagnosis_id"` // Id of diagnosis

	// Output
	AuthKeyResponses *AuthKeyResponse           `tfsdk:"auth_key_responses"`
	ScheduleResponse *DiagnosisScheduleResponse `tfsdk:"schedule_response"`
	SummaryResponses *SummaryResponse           `tfsdk:"summary_responses"`
}

type AuthKeyResponse struct {
	AuthKeyCreatedAt types.String `tfsdk:"auth_key_created_at"` // created date of authkey
	AuthKeyExpiredAt types.String `tfsdk:"auth_key_expired_at"` // expired date of authkey
	AuthKeyId        types.String `tfsdk:"auth_key_id"`         // Id of auth key
	AuthKeyState     types.String `tfsdk:"auth_key_state"`      // state of auth key
	UserId           types.String `tfsdk:"user_id"`             // user Id
}

type DiagnosisScheduleResponse struct {
	DiagnosisId               types.String `tfsdk:"diagnosis_id"`                 // Id of diagnosis
	DiagnosisStartTimePattern types.String `tfsdk:"diagnosis_start_time_pattern"` // Start time( 5-minute increments, 00 to 23 hours, 00 to 55 minutes )
	FrequencyType             types.String `tfsdk:"frequency_type"`               // Schedule type( monthly, weekly, daily)
	FrequencyValue            types.String `tfsdk:"frequency_value"`              // Schedule value (01~31, MONDAY~SUNDAY, everyDay)
	UseDiagnosisCheckTypeBp   types.String `tfsdk:"use_diagnosis_check_type_bp"`  // Checklist Best Practice Use
	UseDiagnosisCheckTypeSsi  types.String `tfsdk:"use_diagnosis_check_type_ssi"` // Checklist SSI usage
}

// #endregion

// #region Config inspection diagnosis resource
type ConfigInspectionDiagnosisResource struct {
	// Input fields
	AccountId          types.String              `tfsdk:"account_id"`           // account Id
	AuthKeyRequest     *AuthKeyRequest           `tfsdk:"auth_key_request"`     // Auth key request
	CspType            types.String              `tfsdk:"csp_type"`             // Type of cloud service provider
	DiagnosisAccountId types.String              `tfsdk:"diagnosis_account_id"` // Id of diagnosis
	DiagnosisCheckType types.String              `tfsdk:"diagnosis_check_type"` // Check type of diagnosis
	DiagnosisId        types.String              `tfsdk:"diagnosis_id"`         // Id of diagnosis
	DiagnosisName      types.String              `tfsdk:"diagnosis_name"`       // Name of diagnosis
	DiagnosisType      types.String              `tfsdk:"diagnosis_type"`       // diagnosis Type
	PlanType           types.String              `tfsdk:"plan_type"`            // plan Type
	ScheduleRequest    *DiagnosisScheduleRequest `tfsdk:"schedule_request"`     // Schedule request
	Tags               types.Map                 `tfsdk:"tags"`                 // Tag List

	// Output fields
	Result             types.Bool   `tfsdk:"result"`               // Result
	CreatedDiagnosisId types.String `tfsdk:"created_diagnosis_id"` // Result
}

type AuthKeyRequest struct {
	AuthKeyCreatedAt types.String `tfsdk:"auth_key_created_at"` // created date of authkey
	AuthKeyExpiredAt types.String `tfsdk:"auth_key_expired_at"` // expired date of authkey
	AuthKeyId        types.String `tfsdk:"auth_key_id"`         // Id of auth key
	DiagnosisId      types.String `tfsdk:"diagnosis_id"`        // Id of diagnosis
}

type DiagnosisScheduleRequest struct {
	DiagnosisId               types.String `tfsdk:"diagnosis_id"`                 // Id of diagnosis
	DiagnosisStartTimePattern types.String `tfsdk:"diagnosis_start_time_pattern"` // Start time( 5-minute increments, 00 to 23 hours, 00 to 55 minutes )
	FrequencyType             types.String `tfsdk:"frequency_type"`               // Schedule type( monthly, weekly, daily)
	FrequencyValue            types.String `tfsdk:"frequency_value"`              // Schedule value (01~31, MONDAY~SUNDAY, everyDay)
	UseDiagnosisCheckTypeBp   types.String `tfsdk:"use_diagnosis_check_type_bp"`  // Checklist Best Practice Use
	UseDiagnosisCheckTypeSsi  types.String `tfsdk:"use_diagnosis_check_type_ssi"` // Checklist SSI usage
}

// #endregion

// # region Config inspection diagnosis results
type ConfigInspectionDiagnosisResultListDataSources struct {
	// Input
	WithCount      types.String `tfsdk:"with_count"`      // With count
	Limit          types.Int32  `tfsdk:"limit"`           // Limit
	Marker         types.String `tfsdk:"marker"`          // Marker
	Sort           types.String `tfsdk:"sort"`            // Sort
	AccountId      types.String `tfsdk:"account_id"`      // Account Id
	DiagnosisID    types.String `tfsdk:"diagnosis_id"`    // Id of diagnosis
	DiagnosisName  types.String `tfsdk:"diagnosis_name"`  // Name of diagnosis
	CSPType        types.String `tfsdk:"csp_type"`        // Type of cloud service provider
	DiagnosisState types.String `tfsdk:"diagnosis_state"` // diagnosis state
	StartDate      types.String `tfsdk:"start_date"`      // Start date
	EndDate        types.String `tfsdk:"end_date"`        // End date
	UserId         types.String `tfsdk:"user_id"`         // User id

	// Output
	TotalCount               types.Int32               `tfsdk:"total_count"`                // Count
	DiagnosisResultResponses []DiagnosisResultResponse `tfsdk:"diagnosis_result_responses"` // Diagnosis result responses
	Links                    []Link                    `tfsdk:"links"`                      // Links

}

type DiagnosisResultResponse struct {
	CountCheck               types.Int32  `tfsdk:"count_check"`                // Check count
	CountError               types.Int32  `tfsdk:"count_error"`                // Error count
	CountFail                types.Int32  `tfsdk:"count_fail"`                 // Fail count
	CountNa                  types.Int32  `tfsdk:"count_na"`                   // Na count
	CountPass                types.Int32  `tfsdk:"count_pass"`                 // Pass count
	CspType                  types.String `tfsdk:"csp_type"`                   // Type of cloud service provider
	DiagnosisAccountId       types.String `tfsdk:"diagnosis_account_id"`       // Id of diagnosis
	DiagnosisCheckType       types.String `tfsdk:"diagnosis_check_type"`       // Check type of diagnosis
	DiagnosisId              types.String `tfsdk:"diagnosis_id"`               // Id of diagnosis
	DiagnosisName            types.String `tfsdk:"diagnosis_name"`             // Name of diagnosis
	DiagnosisRequestSequence types.String `tfsdk:"diagnosis_request_sequence"` // Sequence of diagnosis request
	DiagnosisResult          types.String `tfsdk:"diagnosis_result"`           // Diagnosis Result
	DiagnosisTotalCount      types.Int32  `tfsdk:"diagnosis_total_count"`      // Diagnosis Total Count
	ProceedDate              types.String `tfsdk:"proceed_date"`               // Proceed Date
	Total                    types.Int32  `tfsdk:"total"`                      // Total count
}

// #endregion

// #region Config inspection diagnossed result detail
type ConfigInspectionDiagnosisResultDetailDataSource struct {
	// Input
	DiagnosisId              types.String `tfsdk:"diagnosis_id"`               // Id of diagnosis
	DiagnosisRequestSequence types.String `tfsdk:"diagnosis_request_sequence"` // Sequence of diagnosis request
	WithCount                types.String `tfsdk:"with_count"`                 // With count
	Limit                    types.Int32  `tfsdk:"limit"`                      // Limit
	Marker                   types.String `tfsdk:"marker"`                     // Marker
	Sort                     types.String `tfsdk:"sort"`                       // Sort

	// Output
	ChecklistName      types.String            `tfsdk:"checklist_name"`       // Checklist Name
	TotalCount         types.Int32             `tfsdk:"total_count"`          // Count
	DiagnosisAccountId types.String            `tfsdk:"diagnosis_account_id"` // Id of diagnosis
	DiagnosisCheckType types.String            `tfsdk:"diagnosis_check_type"` // Check type of diagnosis
	DiagnosisName      types.String            `tfsdk:"diagnosis_name"`       // Name of diagnosis
	Links              []Link                  `tfsdk:"links"`                // Links
	ProceedDate        types.String            `tfsdk:"proceed_date"`         // Proceed Date
	ResultDetailList   []DiagnosisResultDetail `tfsdk:"result_detail_list"`   // Result detail list
	Total              types.Int32             `tfsdk:"total"`                // Total count
}

type DiagnosisResultDetail struct {
	ActionGuide        string       `tfsdk:"action_guide"`         // Measure guide description
	Changed            types.Bool   `tfsdk:"changed"`              // Is changed?
	DiagnosisCheckType string       `tfsdk:"diagnosis_check_type"` // Check type of diagnosis
	DiagnosisCriteria  string       `tfsdk:"diagnosis_criteria"`   // Decision standard description
	DiagnosisItem      string       `tfsdk:"diagnosis_item"`       // Sub category description
	DiagnosisLayer     string       `tfsdk:"diagnosis_layer"`      // Inspector item category description
	DiagnosisMethod    string       `tfsdk:"diagnosis_method"`     // Inspector method description
	DiagnosisResult    string       `tfsdk:"diagnosis_result"`     // Verify result state
	ResultContents     string       `tfsdk:"result_contents"`      // Result Contents
	SubCategory        types.String `tfsdk:"sub_category"`         // Sub category
}

// #endregion

// # region Config inspection diagnosis request
type ConfigInspectionDiagnosisRequestDataSource struct {
	// Input fields
	AccessKey          types.String `tfsdk:"access_key"`           // Access Key
	DiagnosisCheckType types.String `tfsdk:"diagnosis_check_type"` // Check type of diagnosis
	DiagnosisId        types.String `tfsdk:"diagnosis_id"`         // Id of diagnosis
	SecretKey          types.String `tfsdk:"secret_key"`           // Secret Key
	TenantId           types.String `tfsdk:"tenant_id"`            // Tenant ID

	// Output fields
	Result types.Bool `tfsdk:"result"` // Result of diagnosis request

}

// #endregion
