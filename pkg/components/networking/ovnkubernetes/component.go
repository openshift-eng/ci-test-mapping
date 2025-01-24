package networkingovnkubernetes

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OvnKubernetesComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / ovn-kubernetes",
		Operators:            []string{},
		DefaultJiraComponent: "Networking / ovn-kubernetes",
		Namespaces: []string{
			"openshift-ovn-kubernetes",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-network",
				// Tests that skip a network other than OVN are assumed to belong to us.
				IncludeAll: []string{"Skipped:Network/"},
				ExcludeAny: []string{"Skipped:Network/OVNKubernetes", "Skipped:Network/OVNKuberenetes"},
			},
			{
				IncludeAll: []string{"ovn-kubernetes"},
				Priority:   1,
			},
			{Suite: "OVN related networking scenarios"},
			{Suite: "OVNKubernetes IPsec related networking scenarios"},
			{Suite: "OVNKubernetes Windows Container related networking scenarios"},
			{Suite: "SDN/OVN metrics related networking scenarios"},
			{Suite: "ipv6 dual stack cluster test scenarios"},
			{Suite: "sdn2ovn migration testing"},
		},
		TestRenames: map[string]string{
			"[Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ovn-kubernetes":    "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ovn-kubernetes",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ovn-kubernetes": "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ovn-kubernetes",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes does not mirror EndpointSlices in namespaces not using user defined primary networks L2 dualstack primary UDN [Suite:openshift/conformance/parallel]":                                               "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions does not mirror EndpointSlices in namespaces not using user defined primary networks L2 dualstack primary UDN [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes does not mirror EndpointSlices in namespaces not using user defined primary networks L3 dualstack primary UDN [Suite:openshift/conformance/parallel]":                                               "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions does not mirror EndpointSlices in namespaces not using user defined primary networks L3 dualstack primary UDN [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L2 dualstack primary UDN, cluster-networked pods [Suite:openshift/conformance/parallel]": "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L2 primary UDN, cluster-networked pods [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L2 dualstack primary UDN, host-networked pods [Suite:openshift/conformance/parallel]":    "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L2 primary UDN, host-networked pods [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L3 dualstack primary UDN, cluster-networked pods [Suite:openshift/conformance/parallel]": "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L3 primary UDN, cluster-networked pods [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L3 dualstack primary UDN, host-networked pods [Suite:openshift/conformance/parallel]":    "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] EndpointSlices mirroring when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions mirrors EndpointSlices managed by the default controller for namespaces with user defined primary networks L3 primary UDN, host-networked pods [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes can perform east/west traffic between nodes for two pods connected over a L2 dualstack primary UDN [Suite:openshift/conformance/parallel]":                                                                                   "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions can perform east/west traffic between nodes for two pods connected over a L2 primary UDN [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes can perform east/west traffic between nodes two pods connected over a L3 dualstack primary UDN [Suite:openshift/conformance/parallel]":                                                                                       "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions can perform east/west traffic between nodes two pods connected over a L3 primary UDN [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes is isolated from the default network with L3 dualstack primary UDN [Suite:openshift/conformance/parallel]":                                                                                                                   "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions is isolated from the default network with L3 primary UDN [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes is isolated from the default network with L2 dualstack primary UDN [Suite:openshift/conformance/parallel]":                                                                                                                   "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes created using NetworkAttachmentDefinitions is isolated from the default network with L2 primary UDN [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes UserDefinedNetwork should delete NetworkAttachmentDefinition when UserDefinedNetwork is deleted [Suite:openshift/conformance/parallel]":                                                                                      "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes UserDefinedNetwork CRD controller should delete NetworkAttachmentDefinition when UserDefinedNetwork is deleted [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes UserDefinedNetwork should create NetworkAttachmentDefinition according to spec [Suite:openshift/conformance/parallel]":                                                                                                       "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes UserDefinedNetwork CRD controller should create NetworkAttachmentDefinition according to spec [Suite:openshift/conformance/parallel]",
			"[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes UserDefinedNetwork pod connected to UserDefinedNetwork cannot be deleted when being used [Suite:openshift/conformance/parallel]":                                                                                             "[sig-network][OCPFeatureGate:NetworkSegmentation][Feature:UserDefinedPrimaryNetworks] when using openshift ovn-kubernetes UserDefinedNetwork CRD controller pod connected to UserDefinedNetwork cannot be deleted when being used [Suite:openshift/conformance/parallel]",
		},
	},
}

func (c *Component) IdentifyTest(test *v1.TestInfo) (*v1.TestOwnership, error) {
	if matcher := c.FindMatch(test); matcher != nil {
		jira := matcher.JiraComponent
		if jira == "" {
			jira = c.DefaultJiraComponent
		}
		return &v1.TestOwnership{
			Name:          test.Name,
			Component:     c.Name,
			JIRAComponent: jira,
			Priority:      matcher.Priority,
			Capabilities:  append(matcher.Capabilities, identifyCapabilities(test)...),
		}, nil
	}

	return nil, nil
}

func (c *Component) StableID(test *v1.TestInfo) string {
	// Look up the stable name for our test in our renamed tests map.
	if stableName, ok := c.TestRenames[test.Name]; ok {
		return stableName
	}
	return test.Name
}

func (c *Component) JiraComponents() (components []string) {
	components = []string{c.DefaultJiraComponent}
	for _, m := range c.Matchers {
		components = append(components, m.JiraComponent)
	}

	return components
}
