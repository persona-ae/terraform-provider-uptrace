package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// - "github.com/hashicorp/terraform-plugin-framework/types"
type TFMonitorData struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	ProjectID types.Int64  `tfsdk:"project_id"`
	Status    types.String `tfsdk:"status"`
	Type      types.String `tfsdk:"type"`
	Query     types.String `tfsdk:"query"`
}
