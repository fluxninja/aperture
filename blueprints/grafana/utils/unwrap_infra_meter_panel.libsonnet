local infraMeterPanelLibrary = import '../infra_meter_panel_library.libsonnet';

function(receiverName, policyName, infraMeter, datasource, extraFilters) {

  // Take the first part of the split as string to look for the receiver name
  local receiverKey = std.split(receiverName, '/')[0],

  // If the receiverKey exists in the infraMeterPanelLibrary, generate panels, otherwise return an empty array
  local generatedPanels = if std.objectHas(infraMeterPanelLibrary, receiverKey)
  then infraMeterPanelLibrary[receiverKey](policyName, infraMeter, datasource, extraFilters)
  else { panel: [] },  // Return an object with an empty 'panel' key if receiverKey does not exist

  panel:
    if std.objectHas(generatedPanels, 'panels')
    then generatedPanels.panels
    else [generatedPanels.panel],
}
