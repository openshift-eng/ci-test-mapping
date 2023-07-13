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
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-network",
				// Tests that skip a network other than OVN are assumed to belong to us.
				IncludeAll: []string{"Skipped:Network/"},
				ExcludeAny: []string{"Skipped:Network/OVNKubernetes", "Skipped:Network/OVNKuberenetes"},
			},
			{
				Namespaces: []string{
					"openshift-ovn-kubernetes",
				},
				Priority: 1,
			},
			{
				IncludeAll: []string{"ovn-kubernetes"},
				Priority:   1,
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
