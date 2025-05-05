package uptrace

// start response models

type GetMonitorsResponse struct {
	Count    int       `json:"count"`
	Monitors []Monitor `json:"monitors"`
}

type GetMonitorByIdResponse Monitor

// start response-model vocabulary

type Monitor struct {
	ChannelIDs            []int          `json:"channelIds"`
	CheckedAt             int64          `json:"checkedAt"`
	CreatedAt             float64        `json:"createdAt"`
	Error                 string         `json:"error"`
	ID                    int            `json:"id"`
	Name                  string         `json:"name"`
	NextCheckTime         int64          `json:"nextCheckTime"`
	NotifyEveryoneByEmail bool           `json:"notifyEveryoneByEmail"`
	Params                Params         `json:"params"`
	ProjectID             int            `json:"projectId"`
	RepeatInterval        RepeatInterval `json:"repeatInterval"`
	Status                string         `json:"status"`
	TeamIDs               []int          `json:"teamIds"`
	Type                  string         `json:"type"`
	UpdatedAt             float64        `json:"updatedAt"`
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
