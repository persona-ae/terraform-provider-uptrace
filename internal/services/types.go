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

type MonitorIdResponse struct {
	Monitor monitorId `json:"monitor"`
}

type MonitorRequest struct {
	monitorBase
}

type MonitorResponse struct {
	monitorBase

	ID             int            `json:"id"`
	ProjectID      int            `json:"projectId"`
	Status         string         `json:"status"`
	UpdatedAt      float64        `json:"updatedAt"`
	CheckedAt      int64          `json:"checkedAt"`
	CreatedAt      float64        `json:"createdAt"`
	Error          string         `json:"error"`
	NextCheckTime  int64          `json:"nextCheckTime"`
	RepeatInterval repeatInterval `json:"repeatInterval"`
}

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

	Params Params `json:"params"`
}

type Params struct {
	Metrics []Metric `json:"metrics"`
	Query   string   `json:"query"`

	// optional fields below
	Column           *string  `json:"column"`
	MinAllowedValue  *float32 `json:"minAllowedValue"`
	MaxAllowedValue  *float32 `json:"maxAllowedValue"`
	GroupingInterval *float32 `json:"groupingInterval,omitempty"`
	CheckNumPoint    *int     `json:"checkNumPoint,omitempty"`
	NullsMode        *string  `json:"nullsMode,omitempty"`
	TimeOffset       *float32 `json:"timeOffset,omitempty"`
}

type Metric struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type repeatInterval struct {
	Strategy string `json:"strategy"`
}

type monitorId struct {
	Id string `json:"id"`
}
