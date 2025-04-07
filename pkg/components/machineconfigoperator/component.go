package machineconfigoperator

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var MachineConfigOperatorComponent = Component{
	Component: &config.Component{
		Name:                 "Machine Config Operator",
		Operators:            []string{"machine-config"},
		DefaultJiraComponent: "Machine Config Operator",
		Namespaces: []string{
			"openshift-machine-config-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Machine Config Operator"},
			},
			{
				IncludeAll: []string{"bz-machine config operator"},
			},
			{
				IncludeAll: []string{"machine-config-operator"},
			},
			{
				IncludeAny: []string{
					"machine-config-operator",
					"node count should match or exceed machine count",
					"OSUpdateStarted event should be recorded for nodes that reach OSUpdateStaged",
					"upgrade should succeed: machine-config",
				},
			},
			{
				SIG:          "sig-cluster-lifecycle",
				IncludeAll:   []string{"Pods cannot access the /config"},
				Capabilities: []string{"Config"},
			},
			{Suite: "MCO"},
			{
				SIG: "sig-mco",
			},
		},
		TestRenames: map[string]string{
			"[Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-machine-config-operator":    "[bz-Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-machine-config-operator",
			"[Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-machine-config-operator": "[bz-Machine Config Operator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-machine-config-operator",
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
