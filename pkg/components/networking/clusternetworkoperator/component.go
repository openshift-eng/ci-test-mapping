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
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Networking"},
			},
			{
				IncludeAll: []string{"bz-networking"},
			},
			{
				Namespaces: []string{
					"openshift-host-network",
					"openshift-network-diagnostics",
					"openshift-network-operator",
					"openshift-kni-infra",
				},
				Priority: 1,
			},
			{
				SIG: "sig-network",
				// Everything not already matched goes here
				Priority: -1,
			},
			{
				IncludeAll: []string{"cluster-network-operator"},
			},
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
