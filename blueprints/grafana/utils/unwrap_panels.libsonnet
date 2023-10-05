local panelLibrary = import '../panel_library.libsonnet';

function(datasourceName, policyName, component, extra_filters={}) {
  local generatedPanels = panelLibrary[component.component_name](datasourceName, policyName, component, extra_filters),

  // this can be either a group of panels or single panel - we have to unwrap
  panel:
    if std.objectHas(generatedPanels, 'panels')
    then
      generatedPanels.panels
    else
      [generatedPanels.panel],
}
