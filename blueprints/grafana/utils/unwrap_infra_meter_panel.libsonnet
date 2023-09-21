local panelLibrary = import '../panel_library.libsonnet';

function(receiverName, policyName, infraMeterName, datasource, extraFilters) {

  local generatedPanels = panelLibrary[receiverName](policyName, infraMeterName, datasource, extraFilters),

  panel:
    if std.objectHas(generatedPanels, 'panels')
    then
      generatedPanels.panels
    else
      [generatedPanels.panel],
}
