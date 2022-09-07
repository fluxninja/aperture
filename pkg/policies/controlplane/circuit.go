package controlplane

import (
	"fmt"
	"strings"

	"github.com/emicklei/dot"
)

// ComponentDTO takes a CompiledCircuit and returns its graph representation.
func ComponentDTO(circuit CompiledCircuit) ([]Component, []Link) {
	var componentsDTO []Component
	var links []Link
	type componentData struct {
		componentID string
		portName    string
	}
	outSignalsIndex := make(map[string][]componentData)
	inSignalsIndex := make(map[string][]componentData)
	for _, c := range circuit {
		var inPorts, outPorts []Port
		for name, signals := range c.InPortToSignalsMap {
			signalName := signals[0].Name
			inPorts = append(inPorts, Port{
				Name:   name,
				Signal: signalName,
				Looped: signals[0].Looped,
			})
			inSignalsIndex[signalName] = append(inSignalsIndex[signalName], componentData{
				componentID: c.ComponentID,
				portName:    name,
			})
		}
		for name, signals := range c.OutPortToSignalsMap {
			signalName := signals[0].Name
			outPorts = append(outPorts, Port{
				Name:   name,
				Signal: signalName,
				Looped: signals[0].Looped,
			})
			outSignalsIndex[signalName] = append(outSignalsIndex[signalName], componentData{
				componentID: c.ComponentID,
				portName:    name,
			})
		}
		componentsDTO = append(componentsDTO, Component{
			ComponentID:       c.ComponentID,
			ComponentName:     c.CompiledComponent.Name,
			ComponentType:     string(c.CompiledComponent.ComponentType),
			Component:         c.CompiledComponent.MapStruct,
			InPorts:           inPorts,
			OutPorts:          outPorts,
			ParentComponentID: c.ParentComponentID,
		})
	}
	// compute links
	for signalName := range outSignalsIndex {
		for _, outComponent := range outSignalsIndex[signalName] {
			for _, inComponent := range inSignalsIndex[signalName] {
				links = append(links, Link{
					Source: SourceTarget{
						ComponentID: outComponent.componentID,
						PortName:    outComponent.portName,
					},
					Target: SourceTarget{
						ComponentID: inComponent.componentID,
						PortName:    inComponent.portName,
					},
					SignalName: signalName,
				})
			}
		}
	}
	return componentsDTO, links
}

// DOT returns Components and Links as a DOT graph description.
func DOT(components []Component, links []Link) string {
	g := dot.NewGraph(dot.Directed)
	g.AttributesMap.Attr("splines", "ortho")
	g.AttributesMap.Attr("rankdir", "LR")

	// indexed by component id
	clusters := make(map[string]*dot.Graph)
	for _, c := range components {
		cluster := g.Subgraph(fmt.Sprintf("%s (%s)", c.ComponentName, strings.SplitN(c.ComponentID, ".", 1)[0]), dot.ClusterOption{})
		cluster.AttributesMap.Attr("margin", "50.0")
		clusters[c.ComponentID] = cluster
		var anyIn, anyOut dot.Node
		for _, p := range c.InPorts {
			anyIn = cluster.Node(p.Name)
			cluster.AddToSameRank("input", anyIn)
		}
		for _, p := range c.OutPorts {
			anyOut = cluster.Node(p.Name)
			cluster.AddToSameRank("output", anyOut)
		}
		if len(c.InPorts) > 0 && len(c.OutPorts) > 0 {
			cluster.Edge(anyIn, anyOut).Attr("style", "invis")
		}
	}
	for _, l := range links {
		g.Edge(clusters[l.Source.ComponentID].Node(l.Source.PortName),
			clusters[l.Target.ComponentID].Node(l.Target.PortName)).Attr("label", l.SignalName)
	}
	return g.String()
}

// Port is enables Component connection.
type Port struct {
	Name   string `json:"portName"`
	Signal string `json:"signalName"`
	Looped bool   `json:"looped"`
}

type jsonb map[string]any

// Component is a computational block that forms the circuit.
type Component struct {
	ComponentID       string `json:"componentID"`
	ComponentName     string `json:"componentName"`
	ComponentType     string `json:"componentType"`
	Component         jsonb  `json:"component,omitempty"`
	InPorts           []Port `json:"inPorts"`
	OutPorts          []Port `json:"outPorts"`
	ParentComponentID string `json:"parentComponentID,omitempty"`
}

// SourceTarget describes a link attachment to a component.
type SourceTarget struct {
	ComponentID string `json:"componentID"`
	PortName    string `json:"portName"`
}

// Link is a connection between Components.
type Link struct {
	Source     SourceTarget `json:"source"`
	Target     SourceTarget `json:"target"`
	SignalName string       `json:"signalName"`
}
