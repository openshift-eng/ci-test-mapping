package tracinguiplugin

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ClusterObservabilityOperatorComponent = Component{
	Component: &config.Component{
		Name:                 "tracing-uiplugin",
		Operators:            []string{"cluster-observability-operator"},
		DefaultJiraComponent: "tracing-uiplugin",
		Matchers: []config.ComponentMatcher{
			{
				Suite: "tracing-uiplugin",
			},
		},
		TestRenames: map[string]string{
			"tracing-uiplugin [Capability:UIPlugin][Capability:TraceVisualization][Capability:TimeRange] Test Distributed Traces Cutoffbox functionality":                                                                                             "OpenShift Distributed Tracing UI Plugin tests Test Distributed Traces Cutoffbox functionality",
			"tracing-uiplugin [Capability:UIPlugin][Capability:EmptyState] Test Distributed Tracing UI plugin page without any Tempo instances":                                                                                                       "OpenShift Distributed Tracing UI Plugin tests Test Distributed Tracing UI plugin page without any Tempo instances",
			"tracing-uiplugin [Capability:UIPlugin][Capability:TraceLimits] Test trace limit functionality":                                                                                                                                           "OpenShift Distributed Tracing UI Plugin tests Test trace limit functionality",
			"tracing-uiplugin [Capability:UIPlugin][Capability:TraceVisualization][Capability:RBAC] Test Distributed Tracing UI plugin with Tempo instances and verify traces using user having cluster-admin role":                                   "OpenShift Distributed Tracing UI Plugin tests Test Distributed Tracing UI plugin with Tempo instances and verify traces using user having cluster-admin role",
			"tracing-uiplugin [Capability:UIPlugin][Capability:TraceVisualization][Capability:SpanLinks][Capability:RBAC] Test Distributed Tracing UI plugin with Tempo instances and verify traces, span links using user having cluster-admin role": "OpenShift Distributed Tracing UI Plugin tests Test Distributed Tracing UI plugin with Tempo instances and verify traces, span links using user having cluster-admin role",
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

func (c *Component) ListNamespaces() []string {
	return c.Namespaces
}
