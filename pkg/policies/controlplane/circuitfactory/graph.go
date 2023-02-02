package circuitfactory

import (
	"fmt"
	"strings"

	"github.com/emicklei/dot"
	"google.golang.org/protobuf/types/known/structpb"

	policymonitoringv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/monitoring/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// ToGraphView creates a graph (Currently showing depth=1) representation of a Circuit.
func (circuit *Circuit) ToGraphView() ([]*policymonitoringv1.ComponentView, []*policymonitoringv1.Link) {
	var componentsDTO []*policymonitoringv1.ComponentView
	var links []*policymonitoringv1.Link
	type componentData struct {
		componentID string
		portName    string
	}
	outSignalsIndex := make(map[string][]componentData)
	inSignalsIndex := make(map[string][]componentData)
	for _, ch := range circuit.Tree.Children {
		c := ch.Root
		var inPorts, outPorts []*policymonitoringv1.PortView
		for name, signals := range c.PortMapping.Ins {
			for _, signal := range signals {
				if signal.SignalType() == runtime.SignalTypeNamed {
					signalName := signal.SignalName
					inPorts = append(inPorts, &policymonitoringv1.PortView{
						PortName: name,
						Value:    &policymonitoringv1.PortView_SignalName{SignalName: signalName},
						Looped:   signal.Looped,
					})
					inSignalsIndex[signalName] = append(inSignalsIndex[signalName], componentData{
						componentID: c.ComponentID.String(),
						portName:    name,
					})
				} else if signal.SignalType() == runtime.SignalTypeConstant {
					inPorts = append(inPorts, &policymonitoringv1.PortView{
						PortName: name,
						Value:    &policymonitoringv1.PortView_ConstantValue{ConstantValue: signal.ConstantSignal.Value},
					})
				}
			}
		}
		for name, signals := range c.PortMapping.Outs {
			for _, signal := range signals {
				signalName := signal.SignalName
				outPorts = append(outPorts, &policymonitoringv1.PortView{
					PortName: name,
					Value:    &policymonitoringv1.PortView_SignalName{SignalName: signalName},
					Looped:   signal.Looped,
				})
				outSignalsIndex[signalName] = append(outSignalsIndex[signalName], componentData{
					componentID: c.ComponentID.String(),
					portName:    name,
				})
			}
		}
		componentConfig := c.Config
		componentMap, err := structpb.NewStruct(componentConfig)
		if err != nil {
			log.Error().Err(err).Msg("converting component map")
		}

		componentName := c.Name()
		componentDescription := c.ShortDescription()

		componentsDTO = append(componentsDTO, &policymonitoringv1.ComponentView{
			ComponentId:          c.ComponentID.String(),
			ComponentName:        componentName,
			ComponentDescription: componentDescription,
			ComponentType:        string(c.Type()),
			Component:            componentMap,
			InPorts:              convertPortViews(inPorts),
			OutPorts:             convertPortViews(outPorts),
		})
	}
	// compute links
	for signalName := range outSignalsIndex {
		for _, outComponent := range outSignalsIndex[signalName] {
			for _, inComponent := range inSignalsIndex[signalName] {
				links = append(links, &policymonitoringv1.Link{
					Source: &policymonitoringv1.SourceTarget{
						ComponentId: string(outComponent.componentID),
						PortName:    outComponent.portName,
					},
					Target: &policymonitoringv1.SourceTarget{
						ComponentId: string(inComponent.componentID),
						PortName:    inComponent.portName,
					},
					SignalName: signalName,
				})
			}
		}
	}
	return componentsDTO, links
}

func convertPortViews(ports []*policymonitoringv1.PortView) []*policymonitoringv1.PortView {
	var converted []*policymonitoringv1.PortView
	return append(converted, ports...)
}

// Mermaid returns Components and Links as a mermaid graph.
func Mermaid(components []*policymonitoringv1.ComponentView, links []*policymonitoringv1.Link) string {
	var sb strings.Builder
	// prefix for fake constant components
	fakeConstantPrefix := "FakeConstant"
	// constantID for fake constant components
	var constantID int
	// subgraph for each component

	sb.WriteString("flowchart LR\n")

	renderComponentSubGraph := func(component *policymonitoringv1.ComponentView) string {
		linkExists := func(componentID, portName string, links []*policymonitoringv1.Link) bool {
			for _, link := range links {
				if link.Source.ComponentId == componentID && link.Source.PortName == portName {
					return true
				} else if link.Target.ComponentId == componentID && link.Target.PortName == portName {
					return true
				}
			}
			return false
		}

		var s strings.Builder
		if component.ComponentName == "Variable" {
			// lookup value in component.Component struct
			constSignalObj := component.Component.Fields["constant_signal"].GetStructValue()
			specialValue := constSignalObj.Fields["special_value"].GetStringValue()
			value := constSignalObj.Fields["value"].GetNumberValue()
			outPort := component.OutPorts[0].PortName
			// render constant as a circle with value
			if specialValue != "" {
				s.WriteString(fmt.Sprintf("%s((%s))\n", component.ComponentId+outPort, specialValue))
			} else {
				s.WriteString(fmt.Sprintf("%s((%0.2f))\n", component.ComponentId+outPort, value))
			}

			return s.String()
		}
		name := component.ComponentName
		// only show component description
		if component.ComponentDescription != "" {
			description := component.ComponentDescription
			// truncate description if too long (mermaid limitation)
			if len(description) > 27 {
				description = description[:27] + "..."
			}
			name = fmt.Sprintf("<center>%s<br/>%s</center>", name, description)
		}
		s.WriteString(fmt.Sprintf("subgraph %s[%s]\n", component.ComponentId, name))

		if len(component.InPorts) > 0 {
			// add subgraph for inports
			s.WriteString(fmt.Sprintf("subgraph %s_inports[ ]\n", component.ComponentId))
			// style to make inports invisible
			s.WriteString(fmt.Sprintf("style %s_inports fill:none,stroke:none\n", component.ComponentId))
			// InPorts and OutPorts are nodes in the subgraph
			for _, inputPort := range component.InPorts {
				isConstantInput := false
				if _, ok := inputPort.GetValue().(*policymonitoringv1.PortView_ConstantValue); ok {
					isConstantInput = true
				}
				// check if link exists for this port or it's a constant input
				if linkExists(component.ComponentId, inputPort.PortName, links) || isConstantInput {
					s.WriteString(fmt.Sprintf("%s[%s]\n", component.ComponentId+inputPort.PortName, inputPort.PortName))
				}
			}
			s.WriteString("end\n")
		}

		if len(component.OutPorts) > 0 {
			// add subgraph for outports
			s.WriteString(fmt.Sprintf("subgraph %s_outports[ ]\n", component.ComponentId))
			// style to make outports invisible
			s.WriteString(fmt.Sprintf("style %s_outports fill:none,stroke:none\n", component.ComponentId))
			for _, outputPort := range component.OutPorts {
				if linkExists(component.ComponentId, outputPort.PortName, links) {
					s.WriteString(fmt.Sprintf("%s[%s]\n", component.ComponentId+outputPort.PortName, outputPort.PortName))
				}
			}
			s.WriteString("end\n")
		}

		s.WriteString("end\n")
		// fake nodes for constant value ports
		for _, inPort := range component.InPorts {
			if constValue, ok := inPort.GetValue().(*policymonitoringv1.PortView_ConstantValue); ok {
				// Concatenate fakeConstant prefix to constantComponentID to avoid collision with real component IDs
				constantComponentID := fmt.Sprintf("%s%d", fakeConstantPrefix, constantID)
				constantID++
				sb.WriteString(fmt.Sprintf("%s((%0.2f))\n", constantComponentID, constValue.ConstantValue))
				// link constant to component
				sb.WriteString(fmt.Sprintf("%s --> %s\n", constantComponentID, component.ComponentId+inPort.PortName))
			}
		}
		return s.String()
	}

	for _, c := range components {
		sb.WriteString(renderComponentSubGraph(c))
	}

	// links
	for _, link := range links {
		sb.WriteString(fmt.Sprintf("%s --> |%s| %s\n", link.Source.ComponentId+link.Source.PortName, link.SignalName, link.Target.ComponentId+link.Target.PortName))
	}

	return sb.String()
}

// DOT returns Components and Links as a DOT graph description.
func DOT(components []*policymonitoringv1.ComponentView, links []*policymonitoringv1.Link) string {
	g := dot.NewGraph(dot.Directed)
	g.AttributesMap.Attr("splines", "ortho")
	g.AttributesMap.Attr("rankdir", "LR")

	// indexed by component id
	clusters := make(map[string]*dot.Graph)
	for i := range components {
		name := fmt.Sprintf("%s[%s]", components[i].ComponentName, strings.SplitN(components[i].ComponentId, ".", 1)[0])
		cluster := g.Subgraph(name, dot.ClusterOption{})
		cluster.AttributesMap.Attr("margin", "50.0")
		clusters[components[i].ComponentId] = cluster
		var anyIn, anyOut dot.Node
		for _, inPort := range components[i].InPorts {
			anyIn = cluster.Node(inPort.PortName)
			cluster.AddToSameRank("input", anyIn)
			// fake nodes for constant value ports
			if constValue, ok := inPort.GetValue().(*policymonitoringv1.PortView_ConstantValue); ok {
				// Concatenate fakeConstant prefix to constantComponentID to avoid collision with real component IDs
				fromNode := cluster.Node(fmt.Sprintf("%0.2f", constValue.ConstantValue))
				// link constant to component
				cluster.Edge(fromNode, anyIn)
			}
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
