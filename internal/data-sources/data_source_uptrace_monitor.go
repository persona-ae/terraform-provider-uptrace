package data_sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/persona-ae/terraform-provider-uptrace/internal/models"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

var _ datasource.DataSource = &monitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
	return &monitorDataSource{}
}

type monitorDataSource struct {
	client *uptrace.UptraceClient
}

func (d *monitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *monitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*uptrace.UptraceClient)
	}
}

func (d *monitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing monitor by ID from Uptrace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the monitor.",
			},
			"project_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Project ID the monitor belongs to.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the monitor.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The monitor type.",
			},
			"query": schema.StringAttribute{
				Computed:    true,
				Description: "The monitor's query.",
			},
			"notify_everyone_by_email": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to notify everyone by email.",
			},
			"team_ids": schema.SetAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Description: "List of team ids to be notified by email. Overrides notifyEveryoneByEmail.",
			},
			"channel_ids": schema.SetAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Description: "List of channel ids to send notifications.",
			},
		},
	}
}

func (d *monitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.TFMonitorData

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response uptrace.GetMonitorByIdResponse
	err := d.client.GetMonitorById(ctx, config.ID.ValueInt32(), &response)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read monitor: "+err.Error())
		tflog.Error(ctx, "Uptrace Client error: "+err.Error())
		return
	}

	monitor := response.Monitor
	state := models.TFMonitorData{
		ID:        types.Int32Value(monitor.ID),
		Name:      types.StringValue(monitor.Name),
		ProjectID: types.Int32Value(monitor.ProjectID),
		Status:    types.StringValue(monitor.Status),
		Type:      types.StringValue(monitor.Type),
		Query:     types.StringValue(monitor.Params.Query),
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
