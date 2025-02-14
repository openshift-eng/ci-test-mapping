package networkingclusternetworkoperator

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ClusterNetworkOperatorComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / cluster-network-operator",
		Operators:            []string{"networking", "network"},
		DefaultJiraComponent: "Networking / cluster-network-operator",
		Namespaces: []string{
			"openshift-host-network",
			"openshift-network-diagnostics",
			"openshift-network-operator",
			"openshift-kni-infra",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Networking"},
				Priority:   -1,
			},
			{
				IncludeAll: []string{"bz-networking"},
			},
			{
				SIG: "sig-network",
				// Everything not already matched goes here
				Priority: -1,
			},
			{
				IncludeAll: []string{"cluster-network-operator"},
			},
			{
				IncludeAll: []string{"CNO"},
			},
			{Suite: "Operator related networking scenarios"},
		},
		TestRenames: map[string]string{
			"[Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-host-network":           "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-host-network",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-network-diagnostics":    "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-network-diagnostics",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-network-operator":       "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-network-operator",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-host-network":        "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-host-network",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-network-diagnostics": "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-network-diagnostics",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-network-operator":    "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-network-operator",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kni-infra":                 "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kni-infra",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kni-infra":              "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kni-infra",
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
