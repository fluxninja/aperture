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
func ComponentDTO(circuit CompiledCircuit) ([]languagev1.ComponentView, []languagev1.Link) {
	var componentsDTO []languagev1.ComponentView
	var links []languagev1.Link
	type componentData struct {
		componentID string
		portName    string
	}
	outSignalsIndex := make(map[string][]componentData)
	inSignalsIndex := make(map[string][]componentData)
	for _, c := range circuit {
		var inPorts, outPorts []languagev1.PortView
		for name, signals := range c.InPortToSignalsMap {
			signalName := signals[0].Name
			inPorts = append(inPorts, languagev1.PortView{
				PortName:   name,
				SignalName: signalName,
				Looped:     signals[0].Looped,
			})
			inSignalsIndex[signalName] = append(inSignalsIndex[signalName], componentData{
				componentID: c.ComponentID,
				portName:    name,
			})
		}
		for name, signals := range c.OutPortToSignalsMap {
			signalName := signals[0].Name
			outPorts = append(outPorts, languagev1.PortView{
				PortName:   name,
				SignalName: signalName,
				Looped:     signals[0].Looped,
			})
			outSignalsIndex[signalName] = append(outSignalsIndex[signalName], componentData{
				componentID: c.ComponentID,
				portName:    name,
			})
		}
		componentMap, err := structpb.NewStruct(c.CompiledComponent.MapStruct)
		if err != nil {
			log.Trace().Msgf("converting component map: %w", err)
		}
		componentsDTO = append(componentsDTO, languagev1.ComponentView{
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
				links = append(links, languagev1.Link{
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

func convertPortViews(ports []languagev1.PortView) []*languagev1.PortView {
	converted := make([]*languagev1.PortView, len(ports))
	for _, p := range ports {
		converted = append(converted, &p)
	}
	return converted
}

// DOT returns Components and Links as a DOT graph description.
func DOT(components []languagev1.ComponentView, links []languagev1.Link) string {
	g := dot.NewGraph(dot.Directed)
	g.AttributesMap.Attr("splines", "ortho")
	g.AttributesMap.Attr("rankdir", "LR")

	// indexed by component id
	clusters := make(map[string]*dot.Graph)
	for _, c := range components {
		cluster := g.Subgraph(fmt.Sprintf("%s (%s)", c.ComponentName, strings.SplitN(c.ComponentId, ".", 1)[0]), dot.ClusterOption{})
		cluster.AttributesMap.Attr("margin", "50.0")
		clusters[c.ComponentId] = cluster
		var anyIn, anyOut dot.Node
		for _, p := range c.InPorts {
			anyIn = cluster.Node(p.PortName)
			cluster.AddToSameRank("input", anyIn)
		}
		for _, p := range c.OutPorts {
			anyOut = cluster.Node(p.PortName)
			cluster.AddToSameRank("output", anyOut)
		}
		if len(c.InPorts) > 0 && len(c.OutPorts) > 0 {
			cluster.Edge(anyIn, anyOut).Attr("style", "invis")
		}
	}
	for _, l := range links {
		g.Edge(clusters[l.Source.ComponentId].Node(l.Source.PortName),
			clusters[l.Target.ComponentId].Node(l.Target.PortName)).Attr("label", l.SignalName)
	}
	return g.String()
}

type jsonb map[string]any
