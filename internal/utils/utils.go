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
	if !plan.Query.IsUnknown() {
		out.Params.Query = plan.Query.ValueString()
	}
	if !plan.Metrics.IsUnknown() && !plan.Metrics.IsNull() {
		metrics := []uptrace.Metric{}
		metricsList := plan.Metrics.Elements()

		for _, m := range metricsList {
			objVal := m.(types.Object)

			var name string
			var alias string

			if nameAttr, ok := objVal.Attributes()["name"]; ok && !nameAttr.IsNull() {
				name = nameAttr.(types.String).ValueString()
			}
			if aliasAttr, ok := objVal.Attributes()["alias"]; ok && !aliasAttr.IsNull() {
				alias = aliasAttr.(types.String).ValueString()
			}

			metrics = append(metrics, uptrace.Metric{
				Name:  name,
				Alias: alias,
			})
		}
		out.Params.Metrics = metrics
	}

	return nil
}

func OverlayMonitorOnTFMonitorData(ctx context.Context, monitor uptrace.MonitorResponse, data *models.TFMonitorData) diag.Diagnostics {
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

	metrics := make([]attr.Value, 0, len(monitor.Params.Metrics))
	for _, m := range monitor.Params.Metrics {
		obj, diags := types.ObjectValue(map[string]attr.Type{
			"name":  types.StringType,
			"alias": types.StringType,
		}, map[string]attr.Value{
			"name":  types.StringValue(m.Name),
			"alias": types.StringValue(m.Alias),
		})
		if diags.HasError() {
			return diags
		}
		metrics = append(metrics, obj)
	}

	data.Metrics, _ = types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"alias": types.StringType,
		},
	}, metrics)

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
