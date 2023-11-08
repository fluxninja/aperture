local base = import '../../dashboard/base.libsonnet';
local defaultConfig = import '../../dashboard/default-config.libsonnet';
local panelLibrary = import 'panel-library.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(graphObj, policyName, datasource, extraFilters={})
  local dashboard = base('Policy Summary - %s' % policyName, defaultConfig.refresh_interval);
  local internalComponents = graphObj.internal_components;

  local updateYPosition(panelRows) =
    local updateRowY(row, initialY) =
      local rowHeight = if std.length(row) > 0 then row[0].gridPos.h else 0;
      local updatedRow = std.map(
        function(panel) std.mergePatch(panel, { gridPos: { y: initialY } }),
        row
      );
      {
        updatedRow: updatedRow,
        nextY: initialY + rowHeight,
      };

    std.foldl(
      function(acc, row)
        if std.length(row) > 0
        then
          local result = updateRowY(row, acc.initialY);
          {
            updatedPanels: acc.updatedPanels + result.updatedRow,
            initialY: result.nextY,
          }
        else
          acc,
      panelRows,
      { updatedPanels: [], initialY: 0 }
    ).updatedPanels;

  local panelRows = std.flattenArrays(std.filter(function(x) x != null, [
    if std.objectHas(panelLibrary, c.component_name)
    then
      panelLibrary[c.component_name](datasource, policyName, c, extraFilters)
    for c in internalComponents
  ]));
  local panels = updateYPosition(panelRows);

  local retVal = dashboard.baseDashboard + g.dashboard.withPanels(panels);
  retVal
