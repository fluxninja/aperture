local panelLibrary = import '../panel_library.libsonnet';

function(componentName, componentBody, config) {
  local newConfig =
    if componentName == 'query'
    then config { component_body: componentBody }
    else config,

  local generatedPanels = panelLibrary[std.toString(componentName)](newConfig),

  // this can be either a group of panels or single panel - we have to unwrap
  panel:
    if std.objectHas(generatedPanels, 'panels')
    then
      generatedPanels.panels
    else
      [generatedPanels.panel],
}
