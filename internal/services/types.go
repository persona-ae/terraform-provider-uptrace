package uptrace

// start response models

type GetMonitorsResponse struct {
	Count    int       `json:"count"`
	Monitors []Monitor `json:"monitors"`
}

type GetMonitorByIdResponse Monitor

// start response-model vocabulary

type Monitor struct {
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

	Params         Params         `json:"params"`
	CheckedAt      int64          `json:"checkedAt"`
	CreatedAt      float64        `json:"createdAt"`
	Error          string         `json:"error"`
	ID             int            `json:"id"`
	NextCheckTime  int64          `json:"nextCheckTime"`
	ProjectID      int            `json:"projectId"`
	RepeatInterval RepeatInterval `json:"repeatInterval"`
	Status         string         `json:"status"`
	UpdatedAt      float64        `json:"updatedAt"`
}

type Params struct {
	Metrics []Metric `json:"metrics"`
	Query   string   `json:"query"`
}

type Metric struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type RepeatInterval struct {
	Strategy string `json:"strategy"`
}
