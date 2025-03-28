package networkingmultus

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var MultusComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / multus",
		Operators:            []string{},
		DefaultJiraComponent: "Networking / multus",
		Namespaces: []string{
			"openshift-multus",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"multihoming", "multinetworkpolicy", "multus"},
			},

			{Suite: "Multus-CNI ipv6 related scenarios"},
			{Suite: "Multus-CNI related scenarios"},
			{Suite: "SDN multus compoment upgrade testing"},
		},
		TestRenames: map[string]string{
			"[Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-multus":    "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-multus",
			"[Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-multus": "[bz-Networking][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-multus",
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
