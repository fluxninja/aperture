package circuitfactory

import (
	"fmt"
	"strings"

	"github.com/emicklei/dot"
	policymonitoringv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/monitoring/v1"
)

// MermaidGraph returns Components and Links as a mermaid graph.
func MermaidGraph(graph *policymonitoringv1.Graph) string {
	var sb strings.Builder
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
		if component.ComponentName == "Variable" || component.ComponentName == fakeConstantComponentName {
			outPort := component.OutPorts[0].PortName
			// render constant as a circle with value
			s.WriteString(fmt.Sprintf("%s((%s))\n", component.ComponentId+outPort, component.ComponentDescription))
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
				// check if link exists for this port
				if linkExists(component.ComponentId, inputPort.PortName, graph.InternalLinks) {
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
				if linkExists(component.ComponentId, outputPort.PortName, graph.InternalLinks) {
					s.WriteString(fmt.Sprintf("%s[%s]\n", component.ComponentId+outputPort.PortName, outputPort.PortName))
				}
			}
			s.WriteString("end\n")
		}

		s.WriteString("end\n")
		return s.String()
	}

	for _, c := range graph.InternalComponents {
		sb.WriteString(renderComponentSubGraph(c))
	}

	// links
	for _, link := range graph.InternalLinks {
		if link.GetSignalName() != "" {
			// Signal Link
			sb.WriteString(fmt.Sprintf("%s --> |%s| %s\n", link.Source.ComponentId+link.Source.PortName, link.GetSignalName(), link.Target.ComponentId+link.Target.PortName))
		} else {
			// Constant Link
			sb.WriteString(fmt.Sprintf("%s --> %s\n", link.Source.ComponentId+link.Source.PortName, link.Target.ComponentId+link.Target.PortName))
		}
	}

	return sb.String()
}

// DOTGraph returns Components and Links as a DOT graph description.
func DOTGraph(graph *policymonitoringv1.Graph) string {
	g := dot.NewGraph(dot.Directed)
	g.AttributesMap.Attr("splines", "ortho")
	g.AttributesMap.Attr("rankdir", "LR")

	components := graph.InternalComponents
	links := graph.InternalLinks

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
		edge := g.Edge(clusters[links[i].Source.ComponentId].Node(links[i].Source.PortName),
			clusters[links[i].Target.ComponentId].Node(links[i].Target.PortName))
		if links[i].GetSignalName() != "" {
			// Signal Link
			edge.Attr("label", links[i].GetSignalName())
		}
	}
	return g.String()
}
