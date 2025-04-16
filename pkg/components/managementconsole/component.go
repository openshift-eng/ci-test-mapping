package managementconsole

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ManagementConsoleComponent = Component{
	Component: &config.Component{
		Name:                 "Management Console",
		Operators:            []string{"console-operator", "console"},
		DefaultJiraComponent: "Management Console",
		Namespaces: []string{
			"openshift-console",
			"openshift-console-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Management Console"},
			},
			{
				IncludeAll: []string{"UserInterface"},
				Priority:   -1,
			},
			{
				IncludeAny: []string{
					"upgrade should succeed: console",
				},
			},
		},
		TestRenames: map[string]string{
			"[Management Console][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-console":             "[bz-Management Console][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-console",
			"[Management Console][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-console-operator":    "[bz-Management Console][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-console-operator",
			"[Management Console][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-console":          "[bz-Management Console][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-console",
			"[Management Console][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-console-operator": "[bz-Management Console][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-console-operator",

			"[sig-scheduling][Early] The openshift-console console pods [apigroup:console.openshift.io] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]":   "[sig-scheduling][Early] The openshift-console console pods [apigroup:console.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
			"[sig-scheduling][Early] The openshift-console downloads pods [apigroup:console.openshift.io] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-console downloads pods [apigroup:console.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
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
