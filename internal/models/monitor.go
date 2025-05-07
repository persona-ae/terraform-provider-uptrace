package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// - "github.com/hashicorp/terraform-plugin-framework/types"
type TFMonitorData struct {
	ID                    types.String  `tfsdk:"id"`
	Name                  types.String  `tfsdk:"name"`
	ProjectID             types.Int32   `tfsdk:"project_id"`
	Status                types.String  `tfsdk:"status"`
	Type                  types.String  `tfsdk:"type"`
	Query                 types.String  `tfsdk:"query"`
	Metrics               types.List    `tfsdk:"metrics"`
	NotifyEveryoneByEmail types.Bool    `tfsdk:"notify_everyone_by_email"`
	TeamIDs               types.Set     `tfsdk:"team_ids"`
	ChannelIDs            types.Set     `tfsdk:"channel_ids"`
	Column                types.String  `tfsdk:"column"`
	MinAllowedValue       types.Float32 `tfsdk:"min_allowed_value"`
	MaxAllowedValue       types.Float32 `tfsdk:"max_allowed_value"`
}
