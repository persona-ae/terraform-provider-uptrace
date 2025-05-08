package utils

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/persona-ae/terraform-provider-uptrace/internal/models"
	uptrace "github.com/persona-ae/terraform-provider-uptrace/internal/services"
)

func TFMonitorToUptraceMonitor(ctx context.Context, plan models.TFMonitorData, out *uptrace.Monitor) diag.Diagnostics {

	if !plan.TeamIDs.IsUnknown() {
		teamIds, diags := IntListToSlice(ctx, plan.TeamIDs)
		if diags.HasError() {
			return diags
		}
		out.TeamIDs = teamIds
	}
	if !plan.ChannelIDs.IsUnknown() {
		channelIds, diags := IntListToSlice(ctx, plan.ChannelIDs)
		if diags.HasError() {
			return diags
		}
		out.ChannelIDs = channelIds
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

	if !plan.ID.IsUnknown() {
		id, err := strconv.Atoi(plan.ID.ValueString())
		if err != nil {
			return diag.Diagnostics{} // TODO
		}
		out.ID = int32(id)
	}
	if !plan.ProjectID.IsUnknown() {
		out.ProjectID = plan.ProjectID.ValueInt32()
	}
	if !plan.Name.IsUnknown() {
		out.Name = plan.Name.ValueString()
	}
	if !plan.Status.IsUnknown() {
		out.Status = plan.Status.ValueString()
	}
	if !plan.NotifyEveryoneByEmail.IsUnknown() {
		out.NotifyEveryoneByEmail = plan.NotifyEveryoneByEmail.ValueBool()
	}
	if !plan.RepeatInterval.IsUnknown() {
		out.RepeatInterval.Strategy = plan.RepeatInterval.ValueString()
	}
	if !plan.Type.IsUnknown() {
		out.Type = plan.Type.ValueString()
	}

	// params
	if !plan.Query.IsUnknown() {
		out.Params.Query = plan.Query.ValueString()
	}
	if !plan.Column.IsUnknown() {
		out.Params.Column = plan.Column.ValueString()
	}
	if !plan.ColumnUnit.IsUnknown() {
		out.Params.ColumnUnit = plan.ColumnUnit.ValueString()
	}
	if !plan.BoundsSource.IsUnknown() {
		out.Params.BoundsSource = plan.BoundsSource.ValueString()
	}
	if !plan.GroupingInterval.IsUnknown() {
		out.Params.GroupingInterval = plan.GroupingInterval.ValueInt32()
	}
	if !plan.CheckNumPoint.IsUnknown() {
		out.Params.CheckNumPoint = plan.CheckNumPoint.ValueInt32()
	}
	if !plan.NullsMode.IsUnknown() {
		out.Params.NullsMode = plan.NullsMode.ValueString()
	}
	if !plan.TimeOffset.IsUnknown() {
		out.Params.TimeOffset = plan.TimeOffset.ValueInt32()
	}
	if !plan.MinDevValue.IsUnknown() {
		out.Params.MinDevValue = plan.MinDevValue.ValueFloat64()
	}
	if !plan.MinDevFraction.IsUnknown() {
		out.Params.MinDevFraction = plan.MinDevFraction.ValueFloat64()
	}
	if !plan.MinAllowedValue.IsUnknown() {
		out.Params.MinAllowedValue = plan.MinAllowedValue.ValueFloat64Pointer()
	}
	if !plan.MaxAllowedValue.IsUnknown() {
		out.Params.MaxAllowedValue = plan.MaxAllowedValue.ValueFloat64Pointer()
	}
	if !plan.MinAllowedFlappingValue.IsUnknown() {
		out.Params.Flapping.MinAllowedValue = plan.MinAllowedFlappingValue.ValueFloat64Pointer()
	}
	if !plan.MaxAllowedFlappingValue.IsUnknown() {
		out.Params.Flapping.MaxAllowedValue = plan.MaxAllowedFlappingValue.ValueFloat64Pointer()
	}
	if !plan.Tolerance.IsUnknown() {
		out.Params.Tolerance = plan.Tolerance.ValueString()
	}
	if !plan.TrainingPeriod.IsUnknown() {
		out.Params.TrainingPeriod = plan.TrainingPeriod.ValueInt32()
	}

	return nil
}

func OverlayMonitorOnTFMonitorData(ctx context.Context, monitor uptrace.Monitor, data *models.TFMonitorData) diag.Diagnostics {
	var diags diag.Diagnostics

	data.TeamIDs, diags = Int32SliceToList(monitor.TeamIDs)
	if diags.HasError() {
		return diags
	}
	data.ChannelIDs, diags = Int32SliceToList(monitor.ChannelIDs)
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

	idStr := strconv.Itoa(int(monitor.ID))
	data.ID = types.StringValue(idStr)
	data.ProjectID = types.Int32Value(monitor.ProjectID)
	data.Name = types.StringValue(monitor.Name)
	data.Status = types.StringValue(monitor.Status)
	data.NotifyEveryoneByEmail = types.BoolValue(monitor.NotifyEveryoneByEmail)
	data.Type = types.StringValue(monitor.Type)
	data.RepeatInterval = types.StringValue(monitor.RepeatInterval.Strategy)

	data.Tolerance = types.StringValue(monitor.Params.Tolerance)
	data.TrainingPeriod = types.Int32Value(monitor.Params.TrainingPeriod)
	data.Query = types.StringValue(monitor.Params.Query)
	data.Column = types.StringValue(monitor.Params.Column)
	data.ColumnUnit = types.StringValue(monitor.Params.ColumnUnit)
	data.BoundsSource = types.StringValue(monitor.Params.BoundsSource)
	data.GroupingInterval = types.Int32Value(monitor.Params.GroupingInterval)
	data.CheckNumPoint = types.Int32Value(monitor.Params.CheckNumPoint)
	data.NullsMode = types.StringValue(monitor.Params.NullsMode)
	data.TimeOffset = types.Int32Value(monitor.Params.TimeOffset)
	data.MinDevValue = types.Float64Value(monitor.Params.MinDevValue)
	data.MinDevFraction = types.Float64Value(monitor.Params.MinDevFraction)
	data.MinAllowedValue = types.Float64PointerValue(monitor.Params.MinAllowedValue)
	data.MaxAllowedValue = types.Float64PointerValue(monitor.Params.MaxAllowedValue)
	data.MinAllowedFlappingValue = types.Float64PointerValue(monitor.Params.Flapping.MinAllowedValue)
	data.MaxAllowedFlappingValue = types.Float64PointerValue(monitor.Params.Flapping.MaxAllowedValue)

	return nil
}

func IntListToSlice(ctx context.Context, val types.List) ([]int32, diag.Diagnostics) {
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

func Int32SliceToList(ints []int32) (types.List, diag.Diagnostics) {
	values := make([]attr.Value, len(ints))
	for i, v := range ints {
		values[i] = types.Int32Value(v)
	}
	return types.ListValue(types.Int32Type, values)
}
