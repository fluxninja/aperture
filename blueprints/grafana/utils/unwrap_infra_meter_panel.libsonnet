local infraMeterPanelLibrary = import '../infra_meter_panel_library.libsonnet';

function(receiverName, policyName, infraMeter, datasource, extraFilters) {

  // Take the first part of the split as string to look for the receiver name

  local generatedPanels = infraMeterPanelLibrary[std.split(receiverName, '/')[0]](policyName, infraMeter, datasource, extraFilters),

  panel:
    if std.objectHas(generatedPanels, 'panels')
    then
      generatedPanels.panels
    else
      [generatedPanels.panel],
}
