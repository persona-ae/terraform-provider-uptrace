package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// - "github.com/hashicorp/terraform-plugin-framework/types"
type TFMonitorData struct {
	ID                    types.Int32  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ProjectID             types.Int32  `tfsdk:"project_id"`
	Status                types.String `tfsdk:"status"`
	Type                  types.String `tfsdk:"type"`
	Query                 types.String `tfsdk:"query"`
	NotifyEveryoneByEmail types.Bool   `tfsdk:"notify_everyone_by_email"`
	TeamIDs               types.Set    `tfsdk:"team_ids"`
	ChannelIDs            types.Set    `tfsdk:"channel_ids"`
}
