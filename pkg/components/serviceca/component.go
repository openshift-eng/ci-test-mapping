package serviceca

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ServiceCaComponent = Component{
	Component: &config.Component{
		Name:                 "service-ca",
		Operators:            []string{"service-ca"},
		DefaultJiraComponent: "service-ca",
		Namespaces: []string{
			"openshift-service-ca",
			"openshift-service-ca-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"[bz-service-ca]"},
			},
			{
				IncludeAny: []string{
					"upgrade should succeed: service-ca",
				},
			},
		},
		TestRenames: map[string]string{
			"[service-ca][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-service-ca":             "[bz-service-ca][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-service-ca",
			"[service-ca][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-service-ca-operator":    "[bz-service-ca][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-service-ca-operator",
			"[service-ca][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-service-ca":          "[bz-service-ca][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-service-ca",
			"[service-ca][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-service-ca-operator": "[bz-service-ca][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-service-ca-operator",
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
