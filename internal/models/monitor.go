package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// - "github.com/hashicorp/terraform-plugin-framework/types"
type TFMonitorData struct {
	// required

	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Query   types.String `tfsdk:"query"`
	Metrics types.List   `tfsdk:"metrics"`

	// optional

	ProjectID               types.Int32   `tfsdk:"project_id"`
	Status                  types.String  `tfsdk:"status"`
	NotifyEveryoneByEmail   types.Bool    `tfsdk:"notify_everyone_by_email"`
	RepeatInterval          types.String  `tfsdk:"repeat_interval"`
	Column                  types.String  `tfsdk:"column"`
	ColumnUnit              types.String  `tfsdk:"column_unit"`
	BoundsSource            types.String  `tfsdk:"bounds_source"`
	GroupingInterval        types.Int32   `tfsdk:"grouping_interval"`
	CheckNumPoint           types.Int32   `tfsdk:"check_num_point"`
	NullsMode               types.String  `tfsdk:"nulls_mode"`
	TimeOffset              types.Int32   `tfsdk:"time_offset"`
	MinDevValue             types.Float64 `tfsdk:"min_dev_value"`
	MinDevFraction          types.Float64 `tfsdk:"min_dev_fraction"`
	MinAllowedValue         types.Float64 `tfsdk:"min_allowed_value"`
	MaxAllowedValue         types.Float64 `tfsdk:"max_allowed_value"`
	MinAllowedFlappingValue types.Float64 `tfsdk:"min_allowed_flapping_value"`
	MaxAllowedFlappingValue types.Float64 `tfsdk:"max_allowed_flapping_value"`
	Tolerance               types.String  `tfsdk:"tolerance"`
	TrainingPeriod          types.Int32   `tfsdk:"training_period"`
	TeamIDs                 types.List    `tfsdk:"team_ids"`
	ChannelIDs              types.List    `tfsdk:"channel_ids"`
}
