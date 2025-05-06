package uptrace

// start response models

// get

type GetMonitorsResponse struct {
	Count    int               `json:"count"`
	Monitors []MonitorResponse `json:"monitors"`
}

type GetMonitorByIdResponse struct {
	Monitor MonitorResponse `json:"monitor"`
}

// metric monitor

type CreateMonitorRequest struct {
	monitorRequest
}

type CreateMonitorResponse struct {
	monitorIdResponse
}

type UpdateMonitorRequest struct {
	monitorRequest
}

type UpdateMonitorResponse struct {
	monitorIdResponse
}

// error monitor

type CreateErrorMonitorRequest struct {
	errorMonitorRequest
}

type CreateErrorMonitorResponse struct {
	monitorIdResponse
}

type UpdateErrorMonitorRequest struct {
	errorMonitorRequest
}

type UpdateErrorMonitorResponse struct {
	monitorIdResponse
}

// end response models

// start response-model vocabulary

type monitorBase struct {
	// Monitor name.
	Name string `json:"name"`
	// Must be set to metric.
	Type string `json:"type"`
	// Whether to notify everyone by email.
	NotifyEveryoneByEmail bool `json:"notifyEveryoneByEmail"`
	// List of team ids to be notified by email. Overrides notifyEveryoneByEmail.
	TeamIDs []int `json:"teamIds"`
	// List of channel ids to send notifications.
	ChannelIDs []int `json:"channelIds"`
}

type MonitorResponse struct {
	monitorBase

	Params         ParamsResponse `json:"params"`
	ID             int            `json:"id"`
	ProjectID      int            `json:"projectId"`
	Status         string         `json:"status"`
	UpdatedAt      float64        `json:"updatedAt"`
	CheckedAt      int64          `json:"checkedAt"`
	CreatedAt      float64        `json:"createdAt"`
	Error          string         `json:"error"`
	NextCheckTime  int64          `json:"nextCheckTime"`
	RepeatInterval RepeatInterval `json:"repeatInterval"`
}

type monitorRequest struct {
	monitorBase

	Params ParamsRequest `json:"params"`
}

type errorMonitorRequest struct {
	monitorBase

	Params ErrorParamsRequest `json:"params"`
}

type paramsBase struct {
	Metrics []Metric `json:"metrics"`
	Query   string   `json:"query"`
}

type ParamsRequest struct {
	paramsBase

	Column          string  `json:"column"`
	MinAllowedValue float32 `json:"minAllowedValue"`
	MaxAllowedValue float32 `json:"maxAllowedValue"`

	GroupingInterval *float32 `json:"groupingInterval,omitempty"`
	CheckNumPoint    *int     `json:"checkNumPoint,omitempty"`
	NullsMode        *string  `json:"nullsMode,omitempty"`
	TimeOffset       *float32 `json:"timeOffset,omitempty"`
}

type ErrorParamsRequest struct {
	paramsBase
}

type ParamsResponse struct {
	paramsBase
}

type Metric struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type RepeatInterval struct {
	Strategy string `json:"strategy"`
}

type monitorIdResponse struct {
	Monitor MonitorId `json:"monitor"`
}

type MonitorId struct {
	Id string `json:"id"`
}
