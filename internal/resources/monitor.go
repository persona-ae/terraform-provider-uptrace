package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/persona-ae/terraform-provider-uptrace/internal/models"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
	"github.com/persona-ae/terraform-provider-uptrace/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &monitorResource{}
	_ resource.ResourceWithConfigure   = &monitorResource{}
	_ resource.ResourceWithImportState = &monitorResource{}
)

func NewMonitorResource() resource.Resource {
	return &monitorResource{}
}

// monitorResource is the resource implementation.
type monitorResource struct {
	// this client is set by the provider
	client *uptrace.UptraceClient
}

// Metadata returns the resource type name.
func (r *monitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Debug(ctx, "monitorResource.Metadata", map[string]any{"req": req, "resp": resp})

	resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Schema defines the schema for the resource.
func (r *monitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Debug(ctx, "monitorResource.Schema", map[string]any{"req": req, "resp": resp})

	resp.Schema = schema.Schema{
		Description: "Manages a monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitor.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The monitor type ('metric' or 'error').",
			},
			"query": schema.StringAttribute{
				Required:    true,
				Description: "The monitor's query eg. \"perMin(sum($spans)) as spans\".",
			},
			"metrics": schema.ListAttribute{
				Required:    true,
				Description: "List of metrics to monitor eg. [{\"name\": \"uptrace_tracing_spans\", \"alias\": \"spans\"}].",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":  types.StringType,
						"alias": types.StringType,
					},
				},
			},
			// begin optionals
			"repeat_interval": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Notification repeat interval",
				MarkdownDescription: `Notification repeat interval
By default, Uptrace uses adaptive interval to wait before sending a notification again.

The interval starts from 15 minutes and doubles every 3 notifications, e.g. 15m, 15m, 15m, 30m, 30m, 30m, 1h...

The max interval is 24 hours.
`,
			},
			"column_unit": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The unit of the metric in the selected column",
			},
			"nulls_mode": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Nulls handling mode: allow, forbid, convert. The default is allow.",
			},
			"tolerance": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The tolerance of the automaticly triggered monitor (low, medium, or high).",
				MarkdownDescription: `The tolerance of the automaticly triggered monitor (low, medium, or high).
To reduce the number of alers, pick higher tolerance.
`,
			},
			"notify_everyone_by_email": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Whether to notify everyone by email.",
			},
			"min_dev_value": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Min deviation value",
			},
			"min_dev_fraction": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Min deviation fraction",
			},
			"min_allowed_value": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Inclusive. Values lower than this are reported (At least min_allowed_value or max_allowed_value is required).",
			},
			"max_allowed_value": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Inclusive. Values greater than this are reported (At least min_allowed_value or max_allowed_value is required).",
			},
			"min_allowed_flapping_value": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Min allowed number",
				MarkdownDescription: `Min allowed number
Flapping occures when the monitor triggers the same alert for a short period of time because the monitored value changes back and forth around the trigger point. To reduce the noise, you can configure additional conditions required to close the alert.
For example, the filesystem utilization monitor may fluctuate from 0.89 to 0.9, causing the alert status to change constantly. By configuring the maximum allowed value to 0.85, the alert won't be closed until the value changes from 0.9 to 0.85.
`,
			},
			"max_allowed_flapping_value": schema.Float64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Max allowed number (trigger value: 500)",
				MarkdownDescription: `Max allowed number (trigger value: 500)
Flapping occures when the monitor triggers the same alert for a short period of time because the monitored value changes back and forth around the trigger point. To reduce the noise, you can configure additional conditions required to close the alert.
For example, the filesystem utilization monitor may fluctuate from 0.89 to 0.9, causing the alert status to change constantly. By configuring the maximum allowed value to 0.85, the alert won't be closed until the value changes from 0.9 to 0.85.
`,
			},
			"training_period": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Training period",
				MarkdownDescription: `Training period
Use smaller training periods for volatile values such as CPU usage.
`,
			},
			"time_offset": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Time offset in milliseconds, e.g. 60000 delays check by 1 minute.",
			},
			"grouping_interval": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Grouping interval in milliseconds. The default 60000 (1 minute).",
			},
			"check_num_point": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Number of points to check. The default is 5.",
			},
			"team_ids": schema.ListAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Optional:    true,
				Description: "List of team ids to be notified by email. Overrides notifyEveryoneByEmail.",
			},
			"channel_ids": schema.ListAttribute{
				ElementType: types.Int32Type,
				Computed:    true,
				Optional:    true,
				Description: "List of channel ids to send notifications.",
			},
			// begin computed
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the monitor.",
			},
			"column": schema.StringAttribute{
				Computed:    true,
				Description: "Column name to monitor, eg. spans.",
			},
			"project_id": schema.Int32Attribute{
				Computed:    true,
				Description: "The ID of the project this monitor is associated with.",
			},
			"bounds_source": schema.StringAttribute{
				Computed:    true,
				Description: "Bounds trigger source (manual or auto).",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *monitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "monitorResource.Configure", map[string]any{"req": req, "resp": resp})
	if req.ProviderData == nil {
		return
	}

	// extract the client from the provider data
	client, ok := req.ProviderData.(*uptrace.UptraceClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *uptrace.UptraceClient, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create a new resource.
func (r *monitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "monitorResource.Create", map[string]any{"req": req, "resp": resp})
	var plan models.TFMonitorData

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "planned query", map[string]any{"query": plan.Query.ValueString()})

	// Generate API request body from plan
	monitor := uptrace.MakeMonitorWithDefaults()
	diags = utils.TFMonitorToUptraceMonitor(ctx, plan, &monitor)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating monitor", map[string]any{"monitor": monitor, "query": monitor.Params.Query})

	// Create new monitor
	var response uptrace.MonitorResponse
	err := r.client.CreateMonitor(ctx, monitor, &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create monitor",
			fmt.Sprintf("Failed to create monitor: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "CreateMonitor OK: %s", map[string]any{"response": response})

	// Save data into Terraform state
	diags = utils.OverlayMonitorOnTFMonitorData(ctx, response.Monitor, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read resource information.
func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "monitorResource.Read", map[string]any{"req": req, "resp": resp})

	// Get current state
	// Read data from Terraform state
	var state models.TFMonitorData
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get fresh state from uptrace
	// Generate API request body from plan
	var response uptrace.MonitorResponse
	err := r.client.GetMonitorById(ctx, state.ID.ValueString(), &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get monitor",
			fmt.Sprintf("Failed to get monitor: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "GetMonitorById OK", map[string]any{"response": response})

	// Set refreshed state
	diags := utils.OverlayMonitorOnTFMonitorData(ctx, response.Monitor, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource information.
func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "monitorResource.Update", map[string]any{"req": req, "resp": resp})

	var plan models.TFMonitorData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()

	monitor := uptrace.MakeMonitorWithDefaults()
	diags := utils.TFMonitorToUptraceMonitor(ctx, plan, &monitor)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response uptrace.MonitorResponse
	err := r.client.UpdateMonitor(ctx, id, monitor, &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update monitor",
			fmt.Sprintf("Failed to update monitor: %s", err),
		)
		return
	}

	// log the response
	tflog.Info(ctx, "UpdateMonitor OK", map[string]any{"response": response})

	diags = utils.OverlayMonitorOnTFMonitorData(ctx, response.Monitor, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete resource information.
func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "monitorResource.Delete", map[string]any{"req": req, "resp": resp})

	var state models.TFMonitorData
	// Read Terraform plan data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	err := r.client.DeleteMonitor(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete monitor",
			fmt.Sprintf("Failed to delete monitor: %s", err),
		)
		return
	}

	// log the response
	var response any
	tflog.Info(ctx, "DeleteMonitor OK", map[string]any{"response": response})
}

func (r *monitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "monitorResource.ImportState", map[string]any{"req": req, "resp": resp})

	// Get fresh state from Uptrace
	id := req.ID
	var response uptrace.MonitorResponse
	err := r.client.GetMonitorById(ctx, id, &response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get monitor",
			fmt.Sprintf("Failed to get monitor: %s", err),
		)
		return
	}

	// Save data into Terraform state
	var state models.TFMonitorData
	diags := utils.OverlayMonitorOnTFMonitorData(ctx, response.Monitor, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
