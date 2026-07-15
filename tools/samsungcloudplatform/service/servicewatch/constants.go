package servicewatch

// ----------------------------
// Alert type / level / enums
// ----------------------------
const (
	// --- Alert Types -------------------------------------------------
	AlertTypeMetric    = "METRIC_ALERT"
	AlertTypeService   = "SERVICE_ALERT"
	AlertTypeComposite = "COMPOSITE_ALERT"

	// --- Alert Levels ------------------------------------------------
	AlertLevelHigh   = "HIGH"
	AlertLevelMiddle = "MIDDLE"
	AlertLevelLow    = "LOW"

	// --- Statistic ---------------------------------------------------
	StatSum = "SUM"
	StatAvg = "AVG"
	StatMax = "MAX"
	StatMin = "MIN"

	// --- Operators ---------------------------------------------------
	OpEQ    = "EQ"
	OpNotEQ = "NOT_EQ"
	OpGT    = "GT"
	OpGTE   = "GTE"
	OpLT    = "LT"
	OpLTE   = "LTE"
	OpRange = "RANGE"

	// --- Missing‑Data options ----------------------------------------
	MissingDataMissing      = "MISSING"
	MissingDataBreaching    = "BREACHING"
	MissingDataNotBreaching = "NOT_BREACHING"
	MissingDataIgnore       = "IGNORE"

	// --- Boolean Y/N -------------------------------------------------
	YnYes = "Y"
	YnNo  = "N"
)

// ----------------------------
// Time & formatting
// ----------------------------
const (
	TimeFormatDisplay   = "2006-01-02 15:04:05" // human‑readable (used in convertFromAlertDetailResponse)
)

// ----------------------------
// Error & log message templates
// ----------------------------
const (
	// ---- Generic ----
	ErrUnexpectedConfigure = "Unexpected Data Source Configure Type"

	// ---- Alert CRUD ----
	ErrCreateAlert            = "Error creating Alert"
	ErrReadAlert              = "Error Reading Alert"
	ErrUpdateAlert            = "Error Updating Alert"
	ErrUpdateDescriptionAlert = "Error Updating Alert Description"
	ErrUpdateActivatedAlert   = "Error Updating Alert Activated"
	ErrDeleteAlert            = "Error Deleting Alert"
	ErrGetMetrics             = "Error Get Metrics"
	ErrReadMetrics            = "Error Reading Metrics"
	ErrCreateDashboard        = "Error creating Dashboard"
	ErrReadDashboard          = "Error Reading Dashboard"
	ErrUpdateDashboard        = "Error Updating Dashboard"
	ErrDeleteDashboard        = "Error Deleting Dashboard"

	// ---- Template strings (used with fmt.Sprintf) ----
	ErrUnexpectedConfigureFmt    = "Expected *client.Instance, got: %T. Please report this issue to the provider developers."
	ErrCreateAlertFmt            = "Could not create Alert, unexpected error: %s\nReason: %s"
	ErrReadAlertFmt              = "Could not read Alert ID %s: %s\nReason: %s"
	ErrUpdateAlertFmt            = "Could not update Alert, unexpected error: %s\nReason: %s"
	ErrUpdateDescriptionAlertFmt = "Could not update alert description, unexpected error: %s\nReason: %s"
	ErrUpdateActivatedAlertFmt   = "Could not update alert activated, unexpected error: %s\nReason: %s"
	ErrDeleteAlertFmt            = "Could not delete Alert, unexpected error: %s\nReason: %s"
	ErrGetMetricsFmt             = "Could not get Metrics, unexpected error: %s\nReason: %s\n   namespace_name: %s\n   metricName: %s\n   dimensionKeys: %v"
	ErrReadMetricsFmt            = "Cloud not read Metrics.\nReason: 404 Not Found\n   namespace_name: %s\n   metricName: %s\n   dimensionKeys: %v"
	ErrCreateDashboardFmt        = "Could not create Dashboard, unexpected error: %s\nReason: %s"
	ErrUpdateDashboardFmt        = "Could not update Dashboard, unexpected error: %s\nReason: %s"
	ErrDeleteDashboardFmt        = "Could not delete Dashboard, unexpected error: %s\nReason: %s"
	ErrReadDashboardFmt          = "Could not read Dashboard ID %s: %s\nReason: %s"
)
