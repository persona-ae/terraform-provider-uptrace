package data_sources

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

var _ datasource.DataSource = &monitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
	return &monitorDataSource{}
}

type monitorDataSource struct {
	client *uptrace.UptraceClient
}

func (d *monitorDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "uptrace_monitor"
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
				Required:    true,
				Description: "The ID of the monitor.",
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
			/*
				"metrics": schema.ListNestedAttribute{
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "The name of the metric.",
							},
							"alias": schema.StringAttribute{
								Computed:    true,
								Description: "The metric's alias.",
							},
						},
					},
					Computed:    true,
					Description: "The monitor's metrics.",
				},*/
			// TODO: Add more attributes as needed
		},
	}
}

func (d *monitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := d.client.GetMonitorById(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read monitor: "+err.Error())
		return
	}

	resp.State.Set(ctx, map[string]any{
		"id":         strconv.Itoa(monitor.ID),
		"name":       monitor.Name,
		"project_id": int64(monitor.ProjectID),
		"status":     monitor.Status,
		"type":       monitor.Type,
		"query":      monitor.Params.Query,
		// TODO: validate this soon... "metrics":    monitor.Params.Metrics,
	})
}
