package openshiftcontrollermanagercontrollermanager

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ControllerManagerComponent = Component{
	Component: &config.Component{
		Name:                 "openshift-controller-manager / controller-manager",
		Operators:            []string{"openshift-controller-manager"},
		DefaultJiraComponent: "openshift-controller-manager / controller-manager",
		Namespaces: []string{
			"openshift-controller-manager",
			"openshift-controller-manager-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-openshift-controller-manager"},
			},
			{
				IncludeAny: []string{
					"upgrade should succeed: openshift-controller-manager",
				},
			},
			{Suite: "Check rollout restart and retry in Deployment/DC"},
		},
		TestRenames: map[string]string{
			"[openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-controller-manager":             "[bz-openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-controller-manager",
			"[openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-controller-manager-operator":    "[bz-openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-controller-manager-operator",
			"[openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-controller-manager":          "[bz-openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-controller-manager",
			"[openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-controller-manager-operator": "[bz-openshift-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-controller-manager-operator",
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
