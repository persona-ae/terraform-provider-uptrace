package uptrace

type GetMonitorsResponse struct {
	Count    int       `json:"count"`
	Monitors []Monitor `json:"monitors"`
}

type MonitorResponse struct {
	Monitor Monitor `json:"monitor"`
}

type Monitor struct {
	ID                    int32          `json:"id"`
	ProjectID             int32          `json:"projectId"`
	Name                  string         `json:"name"`
	Status                string         `json:"status"`
	Error                 string         `json:"error"`
	NotifyEveryoneByEmail bool           `json:"notifyEveryoneByEmail"`
	RepeatInterval        RepeatInterval `json:"repeatInterval"`
	Type                  string         `json:"type"`
	TeamIDs               []int32        `json:"teamIds"`
	ChannelIDs            []int32        `json:"channelIds"`
	CreatedAt             float64        `json:"createdAt"`
	UpdatedAt             float64        `json:"updatedAt"`
	CheckedAt             float64        `json:"checkedAt"`
	Params                Params         `json:"params"`
}

type RepeatInterval struct {
	Strategy string `json:"strategy"`
}

type Params struct {
	Metrics          []Metric `json:"metrics"`
	Query            string   `json:"query"`
	Column           string   `json:"column"`
	ColumnUnit       string   `json:"columnUnit"`
	BoundsSource     string   `json:"boundsSource"`
	GroupingInterval int32    `json:"groupingInterval"`
	CheckNumPoint    int32    `json:"checkNumPoint"`
	NullsMode        string   `json:"nullsMode"`
	TimeOffset       int32    `json:"timeOffset"`
	MinAllowedValue  *float64 `json:"minAllowedValue"`
	MaxAllowedValue  *float64 `json:"maxAllowedValue"`
	Flapping         Flapping `json:"flapping"`
	Tolerance        string   `json:"tolerance"`
	TrainingPeriod   int32    `json:"trainingPeriod"`
	MinDevFraction   float64  `json:"minDevFraction"`
	MinDevValue      float64  `json:"minDevValue"`
}

type Metric struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type Flapping struct {
	MinAllowedValue *float64 `json:"minAllowedValue"`
	MaxAllowedValue *float64 `json:"maxAllowedValue"`
}

func MakeMonitorWithDefaults() Monitor {
	minAllowedValue := float64(0)
	return Monitor{
		Status:         "active",
		RepeatInterval: RepeatInterval{Strategy: "default"},
		Params: Params{
			ColumnUnit:       "1",
			BoundsSource:     "manual",
			GroupingInterval: 60000,
			CheckNumPoint:    5,
			NullsMode:        "allow",
			MinDevFraction:   0.2,
			MinAllowedValue:  &minAllowedValue,
			Flapping:         Flapping{},
			Tolerance:        "medium",
			TrainingPeriod:   86400000,
		},
	}
}
