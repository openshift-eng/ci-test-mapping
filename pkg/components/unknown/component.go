package unknown

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var UnknownComponent = Component{
	Component: &config.Component{
		Name:                 "Unknown",
		Operators:            []string{},
		DefaultJiraComponent: "Unknown",
		Matchers:             []config.ComponentMatcher{},
		TestRenames: map[string]string{
			"[sig-cluster-lifecycle] pathological event should not see excessive Back-off restarting failed containers":                   "[sig-cluster-lifecycle] should not see excessive Back-off restarting failed containers",
			"[sig-cluster-lifecycle] pathological event should not see excessive Back-off restarting failed containers in e2e namespaces": "[sig-cluster-lifecycle] should not see excessive Back-off restarting failed containers in e2e namespaces",
			"verify the cluster readiness and stability": "verify the cluster can be connected",
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
