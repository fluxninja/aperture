package controlplane

import (
	"fmt"
	"strings"

	"github.com/emicklei/dot"
	"google.golang.org/protobuf/types/known/structpb"

	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
)

// ComponentDTO takes a CompiledCircuit and returns its graph representation.
func ComponentDTO(circuit CompiledCircuit) ([]*languagev1.ComponentView, []*languagev1.Link) {
	var componentsDTO []*languagev1.ComponentView
	var links []*languagev1.Link
	type componentData struct {
		componentID string
		portName    string
	}
	outSignalsIndex := make(map[string][]componentData)
	inSignalsIndex := make(map[string][]componentData)
	for _, c := range circuit {
		var inPorts, outPorts []*languagev1.PortView
		for name, signals := range c.InPortToSignalsMap {
			for _, signal := range signals {
				signalName := signal.Name
				inPorts = append(inPorts, &languagev1.PortView{
					PortName:   name,
					SignalName: signalName,
					Looped:     signal.Looped,
				})
				inSignalsIndex[signalName] = append(inSignalsIndex[signalName], componentData{
					componentID: c.ComponentID,
					portName:    name,
				})
			}
		}
		for name, signals := range c.OutPortToSignalsMap {
			for _, signal := range signals {
				signalName := signal.Name
				outPorts = append(outPorts, &languagev1.PortView{
					PortName:   name,
					SignalName: signalName,
					Looped:     signal.Looped,
				})
				outSignalsIndex[signalName] = append(outSignalsIndex[signalName], componentData{
					componentID: c.ComponentID,
					portName:    name,
				})
			}
		}
		componentMap, err := structpb.NewStruct(c.CompiledComponent.MapStruct)
		if err != nil {
			log.Trace().Msgf("converting component map: %v", err)
		}
		componentsDTO = append(componentsDTO, &languagev1.ComponentView{
			ComponentId:       c.ComponentID,
			ComponentName:     c.CompiledComponent.Name,
			ComponentType:     string(c.CompiledComponent.ComponentType),
			Component:         componentMap,
			InPorts:           convertPortViews(inPorts),
			OutPorts:          convertPortViews(outPorts),
			ParentComponentId: c.ParentComponentID,
		})
	}
	// compute links
	for signalName := range outSignalsIndex {
		for _, outComponent := range outSignalsIndex[signalName] {
			for _, inComponent := range inSignalsIndex[signalName] {
				links = append(links, &languagev1.Link{
					Source: &languagev1.SourceTarget{
						ComponentId: outComponent.componentID,
						PortName:    outComponent.portName,
					},
					Target: &languagev1.SourceTarget{
						ComponentId: inComponent.componentID,
						PortName:    inComponent.portName,
					},
					SignalName: signalName,
				})
			}
		}
	}
	return componentsDTO, links
}

func convertPortViews(ports []*languagev1.PortView) []*languagev1.PortView {
	var converted []*languagev1.PortView
	for i := range ports {
		converted = append(converted, ports[i])
	}
	return converted
}

// Mermaid returns Components and Links as a mermaid graph.
func Mermaid(components []*languagev1.ComponentView, links []*languagev1.Link) string {
	var sb strings.Builder
	sb.WriteString("flowchart LR\n")

	parentComponents := make(map[string][]*languagev1.ComponentView)
	for i, c := range components {
		if c.ParentComponentId == "" {
			continue
		}
		parentComponents[c.ParentComponentId] = append(parentComponents[c.ParentComponentId], c)
		// remove this element from the slice
		components[i] = components[len(components)-1]
		components = components[:len(components)-1]
	}

	renderComponentSubGraph := func(component *languagev1.ComponentView) string {
		var s strings.Builder
		if component.ComponentName == "Constant" {
			// lookup value in component.Component struct
			value := component.Component.Fields["value"].GetNumberValue()
			outPort := component.OutPorts[0].PortName
			// render constant as a circle with value
			s.WriteString(fmt.Sprintf("%s((%0.2f))\n", component.ComponentId+outPort, value))
			return s.String()
		}
		s.WriteString(fmt.Sprintf("subgraph %s[%s]\n", component.ComponentId, component.ComponentName))
		// InPorts and OutPorts are nodes in the subgraph
		for _, inputPort := range component.InPorts {
			s.WriteString(fmt.Sprintf("%s[%s]\n", component.ComponentId+inputPort.PortName, inputPort.PortName))
		}
		for _, outputPort := range component.OutPorts {
			s.WriteString(fmt.Sprintf("%s[%s]\n", component.ComponentId+outputPort.PortName, outputPort.PortName))
		}
		s.WriteString("end\n")
		return s.String()
	}

	// subgraph for each component
	for _, c := range components {
		// if it's a parent component then render the subgraph for each component
		if _, ok := parentComponents[c.ComponentId]; ok {
			for _, childComponent := range parentComponents[c.ComponentId] {
				childComponent.ComponentName = c.ComponentName + "/" + childComponent.ComponentName
				sb.WriteString(renderComponentSubGraph(childComponent))
			}
		} else {
			sb.WriteString(renderComponentSubGraph(c))
		}
	}

	// links
	for _, link := range links {
		sb.WriteString(fmt.Sprintf("%s --> |%s| %s\n", link.Source.ComponentId+link.Source.PortName, link.SignalName, link.Target.ComponentId+link.Target.PortName))
	}

	return sb.String()
}

// DOT returns Components and Links as a DOT graph description.
func DOT(components []*languagev1.ComponentView, links []*languagev1.Link) string {
	g := dot.NewGraph(dot.Directed)
	g.AttributesMap.Attr("splines", "ortho")
	g.AttributesMap.Attr("rankdir", "LR")

	// indexed by component id
	clusters := make(map[string]*dot.Graph)
	for i := range components {
		var sg *dot.Graph
		if components[i].ParentComponentId == "" {
			sg = g
		} else {
			sg = clusters[components[i].ParentComponentId]
		}
		name := fmt.Sprintf("%s[%s]", components[i].ComponentName, strings.SplitN(components[i].ComponentId, ".", 1)[0])
		cluster := sg.Subgraph(name, dot.ClusterOption{})
		cluster.AttributesMap.Attr("margin", "50.0")
		clusters[components[i].ComponentId] = cluster
		var anyIn, anyOut dot.Node
		for j := range components[i].InPorts {
			anyIn = cluster.Node(components[i].InPorts[j].PortName)
			cluster.AddToSameRank("input", anyIn)
		}
		for j := range components[i].OutPorts {
			anyOut = cluster.Node(components[i].OutPorts[j].PortName)
			cluster.AddToSameRank("output", anyOut)
		}
		if len(components[i].InPorts) > 0 && len(components[i].OutPorts) > 0 {
			cluster.Edge(anyIn, anyOut).Attr("style", "invis")
		}
	}
	for i := range links {
		g.Edge(clusters[links[i].Source.ComponentId].Node(links[i].Source.PortName),
			clusters[links[i].Target.ComponentId].Node(links[i].Target.PortName)).Attr("label", links[i].SignalName)
	}
	return g.String()
}
