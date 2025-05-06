package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/persona-ae/terraform-provider-uptrace/internal/models"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

func TFMonitorToUptraceMonitor(ctx context.Context, plan models.TFMonitorData, out *uptrace.Monitor) diag.Diagnostics {
	teamIds, diags := IntSetToSlice(ctx, plan.TeamIDs)
	if diags.HasError() {
		return diags
	}
	channelIds, diags := IntSetToSlice(ctx, plan.ChannelIDs)
	if diags.HasError() {
		return diags
	}

	if !plan.Name.IsUnknown() {
		out.Name = plan.Name.ValueString()
	}
	if !plan.Type.IsUnknown() {
		out.Type = plan.Type.ValueString()
	}
	if !plan.NotifyEveryoneByEmail.IsUnknown() {
		out.NotifyEveryoneByEmail = plan.NotifyEveryoneByEmail.ValueBool()
	}
	if !plan.TeamIDs.IsUnknown() {
		out.TeamIDs = teamIds
	}
	if !plan.ChannelIDs.IsUnknown() {
		out.ChannelIDs = channelIds
	}

	out.Params = uptrace.Params{}

	return nil
}

func OverlayMonitorOnTFMonitorData(ctx context.Context, monitor uptrace.MonitorResponse, data *models.TFMonitorData) diag.Diagnostics {
	// first the required types
	data.ID = types.Int32Value(monitor.ID)
	data.Name = types.StringValue(monitor.Name)
	data.ProjectID = types.Int32Value(monitor.ProjectID)
	data.Status = types.StringValue(monitor.Status)
	data.Type = types.StringValue(monitor.Type)
	data.Query = types.StringValue(monitor.Params.Query)
	data.NotifyEveryoneByEmail = types.BoolValue(monitor.NotifyEveryoneByEmail)

	var diags diag.Diagnostics
	data.TeamIDs, diags = Int32SliceToSet(monitor.TeamIDs)
	if diags.HasError() {
		return diags
	}
	data.ChannelIDs, diags = Int32SliceToSet(monitor.ChannelIDs)
	if diags.HasError() {
		return diags
	}

	return nil
}

func IntSetToSlice(ctx context.Context, val types.Set) ([]int32, diag.Diagnostics) {
	if val.IsNull() || val.IsUnknown() {
		return nil, nil
	}

	var tfInts []types.Int32
	diags := val.ElementsAs(ctx, &tfInts, false)
	if diags.HasError() {
		return nil, diags
	}
	return convertElementsToInts(tfInts)
}

func convertElementsToInts(tfInts []types.Int32) ([]int32, diag.Diagnostics) {
	ints := make([]int32, len(tfInts))
	for i, s := range tfInts {
		ints[i] = s.ValueInt32()
	}
	return ints, nil
}

func Int32SliceToSet(ints []int32) (types.Set, diag.Diagnostics) {
	values := make([]attr.Value, len(ints))
	for i, v := range ints {
		values[i] = types.Int32Value(v)
	}
	return types.SetValue(types.Int32Type, values)
}
