package circuitfactory

import (
	"fmt"
	"strconv"
	"strings"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policymonitoringv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/monitoring/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	fakeConstantComponentName = "FakeConstant"
)

// Tree is a graph view of a Circuit.
type Tree struct {
	Node     *runtime.ConfiguredComponent
	Children []Tree
}

// RootTree returns Tree to represent Circuit root.
func RootTree(circuitProto *policylangv1.Circuit) (Tree, error) {
	circuitComponent, err := prepareComponent(
		runtime.NewDummyComponent("Circuit", "The Circuit Root", runtime.ComponentTypeStandAlone),
		circuitProto,
		runtime.NewComponentID(runtime.RootComponentID),
		false,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare circuit component")
		return Tree{}, err
	}
	tree := Tree{
		Node: circuitComponent,
	}
	return tree, nil
}

// CreateComponents creates circuit components along with their identifiers and fx options.
//
// Note that number of returned components might be greater than number of
// components in componentsProto, as some components may be composite multi-component stacks or nested circuits.
func (tree *Tree) CreateComponents(
	componentsProto []*policylangv1.Component,
	circuitID runtime.ComponentID,
	policyReadAPI iface.Policy,
) ([]*runtime.ConfiguredComponent, fx.Option, error) {
	var (
		leafComponents []*runtime.ConfiguredComponent
		options        []fx.Option
	)

	for compIndex, componentProto := range componentsProto {
		subTree, leafComps, compOption, err := NewComponentAndOptions(
			componentProto,
			circuitID.ChildID(strconv.Itoa(compIndex)),
			policyReadAPI,
		)
		if err != nil {
			return nil, nil, err
		}
		options = append(options, compOption)

		// Append subTree to tree.Children
		tree.Children = append(tree.Children, subTree)

		// Add subComponents to configuredComponents
		leafComponents = append(leafComponents, leafComps...)
	}

	return leafComponents, fx.Options(options...), nil
}

// TreeGraph walks the tree and gets graph representation of the links amongst the children of each node.
func (tree *Tree) TreeGraph() (*policymonitoringv1.Tree, error) {
	treeMsg, err := treeGraph(tree, *tree)
	if err != nil {
		return nil, err
	}
	actuators := collectActuators(*tree)
	treeMsg.Actuators = actuators
	return treeMsg, nil
}

func collectActuators(current Tree) []*policymonitoringv1.ComponentView {
	var actuators []*policymonitoringv1.ComponentView

	if current.Node.Component.IsActuator() {
		actuators = append(actuators, componentViewFromConfiguredComponent(current.Node))
	}

	for _, child := range current.Children {
		childActuators := collectActuators(child)
		actuators = append(actuators, childActuators...)
	}

	return actuators
}

func treeGraph(root *Tree, current Tree) (*policymonitoringv1.Tree, error) {
	graph, err := root.GetSubGraph(current.Node.ComponentID, 1)
	if err != nil {
		log.Errorf("Error getting subgraph: %v", err)
		return nil, err
	}
	treeMsg := &policymonitoringv1.Tree{
		Node:  componentViewFromConfiguredComponent(current.Node),
		Graph: graph,
	}
	for _, child := range current.Children {
		childTreeMsg, childErr := treeGraph(root, child)
		if childErr != nil {
			log.Errorf("Error getting child tree: %v", childErr)
			return nil, childErr
		}
		treeMsg.Children = append(treeMsg.Children, childTreeMsg)
	}

	return treeMsg, nil
}

// GetSubGraph returns a subgraph of the tree.
func (tree *Tree) GetSubGraph(componentID runtime.ComponentID, depth int) (*policymonitoringv1.Graph, error) {
	internalComponents, externalComponents, err := tree.ExpandSubCircuit(componentID, depth)
	if err != nil {
		return nil, err
	}

	internalLinks, externalLinks := computeLinks(internalComponents, externalComponents)

	internalComponentViews := make([]*policymonitoringv1.ComponentView, len(internalComponents))
	for i, c := range internalComponents {
		internalComponentViews[i] = componentViewFromConfiguredComponent(c)
	}

	// limit external components to ones which are linked to internal components
	// this is to avoid showing components which are not connected to the circuit
	// but are part of the circuit
	externalComponents = filterExternalComponents(externalComponents, externalLinks)
	externalComponentViews := make([]*policymonitoringv1.ComponentView, len(externalComponents))
	for i, c := range externalComponents {
		externalComponentViews[i] = componentViewFromConfiguredComponent(c)
	}

	// Insert fake constant nodes
	fakeConstantNodes, fakeConstantLinks := addFakeConstants(internalComponents)
	internalComponentViews = append(internalComponentViews, fakeConstantNodes...)
	internalLinks = append(internalLinks, fakeConstantLinks...)

	graph := &policymonitoringv1.Graph{
		InternalComponents: internalComponentViews,
		ExternalComponents: externalComponentViews,
		InternalLinks:      internalLinks,
		ExternalLinks:      externalLinks,
	}

	return graph, nil
}

func addFakeConstants(internalComponents []*runtime.ConfiguredComponent) (fakeConstantNodes []*policymonitoringv1.ComponentView, fakeConstantLinks []*policymonitoringv1.Link) {
	for _, c := range internalComponents {
		// Insert fake constant nodes
		for port, signals := range c.PortMapping.Ins {
			for _, signal := range signals {
				if signal.SignalType() == runtime.SignalTypeConstant {
					fakeNodeID := fmt.Sprintf("%s_%s_FakeConstant", c.ComponentID.String(), port)
					fakeConstantNodes = append(fakeConstantNodes, &policymonitoringv1.ComponentView{
						ComponentId:          fakeNodeID,
						ComponentName:        fakeConstantComponentName,
						ComponentDescription: signal.ConstantSignal.Description(),
						ComponentType:        string(runtime.ComponentTypeSource),
						OutPorts: []*policymonitoringv1.PortView{
							{
								PortName: "out",
								Value:    &policymonitoringv1.PortView_ConstantValue{ConstantValue: signal.ConstantSignal.Value},
							},
						},
					})
					fakeConstantLinks = append(fakeConstantLinks, &policymonitoringv1.Link{
						Source: &policymonitoringv1.SourceTarget{
							ComponentId: fakeNodeID,
							PortName:    "out",
						},
						Target: &policymonitoringv1.SourceTarget{
							ComponentId: c.ComponentID.String(),
							PortName:    port,
						},
						Value: &policymonitoringv1.Link_ConstantValue{ConstantValue: signal.ConstantSignal.Value},
					})
				}
			}
		}
	}
	return
}

func filterExternalComponents(externalComponents []*runtime.ConfiguredComponent, externalLinks []*policymonitoringv1.Link) []*runtime.ConfiguredComponent {
	// Create a map to store the filtered external components
	filteredComponents := make(map[string]*runtime.ConfiguredComponent)

	// Check if external components are linked to internal components through external links
	for _, link := range externalLinks {
		for _, component := range externalComponents {
			// Check if the link source or target component ID matches the external component ID
			if link.Source.ComponentId == component.ComponentID.String() || link.Target.ComponentId == component.ComponentID.String() {
				filteredComponents[component.ComponentID.String()] = component
			}
		}
	}

	// Convert the filtered components map to a slice
	filteredComponentSlice := make([]*runtime.ConfiguredComponent, 0, len(filteredComponents))
	for _, component := range filteredComponents {
		filteredComponentSlice = append(filteredComponentSlice, component)
	}

	return filteredComponentSlice
}

// ExpandSubCircuit returns a list of ConfiguredComponents in the circuit with the component at componentID expanded up to given depth.
func (tree *Tree) ExpandSubCircuit(componentID runtime.ComponentID, depth int) ([]*runtime.ConfiguredComponent, []*runtime.ConfiguredComponent, error) {
	var internalComponents, externalComponents []*runtime.ConfiguredComponent

	if componentID.String() == tree.Node.ComponentID.String() {
		internalComponents = tree.ExpandCircuit(depth)
		return internalComponents, nil, nil
	}

	for _, child := range tree.Children {
		if componentID.String() == child.Node.ComponentID.String() {
			internalComponents = child.ExpandCircuit(depth)
		} else if strings.HasPrefix(componentID.String(), child.Node.ComponentID.String()+runtime.NestedComponentDelimiter) {
			internalComponentsFromSubCircuit, externalComponentsFromSubCircuit, err := child.ExpandSubCircuit(componentID, depth)
			if err != nil {
				return nil, nil, err
			}
			externalComponents = append(externalComponents, externalComponentsFromSubCircuit...)
			internalComponents = internalComponentsFromSubCircuit
		} else {
			externalComponents = append(externalComponents, child.Node)
		}
	}
	return internalComponents, externalComponents, nil
}

// ExpandCircuit returns a list of ConfiguredComponents in the circuit expanded up to given depth.
func (tree *Tree) ExpandCircuit(depth int) []*runtime.ConfiguredComponent {
	return tree.collectComponents(depth, 0)
}

func (tree *Tree) collectComponents(maxDepth int, currentDepth int) []*runtime.ConfiguredComponent {
	// If the current depth is greater than or equal to the maximum depth and maxDepth is not -1 end the recursion.
	if maxDepth != -1 && currentDepth >= maxDepth {
		return []*runtime.ConfiguredComponent{tree.Node}
	}

	components := []*runtime.ConfiguredComponent{}
	// If the tree has children, recurse into them.
	if len(tree.Children) > 0 {
		for _, child := range tree.Children {
			childComponents := child.collectComponents(maxDepth, currentDepth+1)
			components = append(components, childComponents...)
		}
	} else {
		// If the tree doesn't have children, add its root component to the list.
		components = append(components, tree.Node)
	}
	return components
}

type signalsIndex struct {
	inSignalsIndex  signalToComponentIndex
	outSignalsIndex signalToComponentIndex
}

type signalToComponentIndex map[runtime.SignalID][]componentData

type componentData struct {
	componentID string
	portName    string
}

func computeLinks(internalComponents, externalComponents []*runtime.ConfiguredComponent) (internalLinks, externalLinks []*policymonitoringv1.Link) {
	createLinks := func(outIndex, inIndex signalToComponentIndex, linkList *[]*policymonitoringv1.Link) {
		for signalID := range outIndex {
			for _, outComponent := range outIndex[signalID] {
				for _, inComponent := range inIndex[signalID] {
					*linkList = append(*linkList, &policymonitoringv1.Link{
						Source: &policymonitoringv1.SourceTarget{
							ComponentId: string(outComponent.componentID),
							PortName:    outComponent.portName,
						},
						Target: &policymonitoringv1.SourceTarget{
							ComponentId: string(inComponent.componentID),
							PortName:    inComponent.portName,
						},
						Value:        &policymonitoringv1.Link_SignalName{SignalName: signalID.SignalName},
						SubCircuitId: signalID.SubCircuitID,
					})
				}
			}
		}
	}

	internalIndex := buildSignalToComponentIndex(internalComponents)
	externalIndex := buildSignalToComponentIndex(externalComponents)

	// Compute internal links.
	createLinks(internalIndex.outSignalsIndex, internalIndex.inSignalsIndex, &internalLinks)

	// Compute incoming external links.
	createLinks(externalIndex.outSignalsIndex, internalIndex.inSignalsIndex, &externalLinks)

	// Compute outgoing external links.
	createLinks(internalIndex.outSignalsIndex, externalIndex.inSignalsIndex, &externalLinks)

	return internalLinks, externalLinks
}

func buildSignalToComponentIndex(components []*runtime.ConfiguredComponent) *signalsIndex {
	index := &signalsIndex{
		inSignalsIndex:  make(signalToComponentIndex),
		outSignalsIndex: make(signalToComponentIndex),
	}

	for _, component := range components {
		for portName, signals := range component.PortMapping.Ins {
			for _, signal := range signals {
				// skip constant signals
				if signal.SignalName == "" {
					continue
				}
				sigID := signal.SignalID()
				componentInfo := componentData{
					componentID: component.ComponentID.String(),
					portName:    portName,
				}
				index.inSignalsIndex[sigID] = append(index.inSignalsIndex[sigID], componentInfo)
			}
		}
		for portName, signals := range component.PortMapping.Outs {
			for _, signal := range signals {
				// skip constant signals
				if signal.SignalName == "" {
					continue
				}
				sigID := signal.SignalID()
				componentInfo := componentData{
					componentID: component.ComponentID.String(),
					portName:    portName,
				}
				index.outSignalsIndex[sigID] = append(index.outSignalsIndex[sigID], componentInfo)
			}
		}
	}

	return index
}

func componentViewFromConfiguredComponent(component *runtime.ConfiguredComponent) *policymonitoringv1.ComponentView {
	var inPorts, outPorts []*policymonitoringv1.PortView
	for name, signals := range component.PortMapping.Ins {
		for _, signal := range signals {
			if signal.SignalType() == runtime.SignalTypeNamed {
				signalName := signal.SignalName
				inPorts = append(inPorts, &policymonitoringv1.PortView{
					PortName:     name,
					Value:        &policymonitoringv1.PortView_SignalName{SignalName: signalName},
					Looped:       signal.Looped,
					SubCircuitId: signal.SubCircuitID,
				})
			} else if signal.SignalType() == runtime.SignalTypeConstant {
				inPorts = append(inPorts, &policymonitoringv1.PortView{
					PortName: name,
					Value:    &policymonitoringv1.PortView_ConstantValue{ConstantValue: signal.ConstantSignal.Value},
				})
			}
		}
	}
	for name, signals := range component.PortMapping.Outs {
		for _, signal := range signals {
			signalName := signal.SignalName
			outPorts = append(outPorts, &policymonitoringv1.PortView{
				PortName:     name,
				Value:        &policymonitoringv1.PortView_SignalName{SignalName: signalName},
				Looped:       signal.Looped,
				SubCircuitId: signal.SubCircuitID,
			})
		}
	}
	componentConfig := component.Config
	componentMap, err := structpb.NewStruct(componentConfig)
	if err != nil {
		log.Error().Err(err).Msg("converting component map")
	}

	componentName := component.Name()
	componentDescription := component.ShortDescription()
	cv := policymonitoringv1.ComponentView{
		ComponentId:          component.ComponentID.String(),
		ComponentName:        componentName,
		ComponentDescription: componentDescription,
		ComponentType:        string(component.Type()),
		Component:            componentMap,
		InPorts:              inPorts,
		OutPorts:             outPorts,
	}
	return &cv
}
