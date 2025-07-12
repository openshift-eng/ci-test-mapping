package cloudcomputeunknown

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var UnknownComponent = Component{
	Component: &config.Component{
		Name:                 "Cloud Compute / Unknown",
		Operators:            []string{"cluster-api", "machine-approver", "machine-api", "control-plane-machine-set"},
		DefaultJiraComponent: "Cloud Compute / Unknown",
		Namespaces: []string{
			"openshift-cluster-machine-approver",
			"openshift-machine-api",
		},
		TestRenames: map[string]string{
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-machine-approver":    "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cluster-machine-approver",
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-machine-api":                 "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-machine-api",
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-machine-approver": "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cluster-machine-approver",
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-machine-api":              "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-machine-api",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"bz-Cloud Compute",
					"bz-cluster-api",
					"bz-control-plane-machine-set",
					"Managed cluster should have same number of Machines and Nodes",
					"control plane machine set operator",
					"upgrade should succeed: control-plane-machine-set",
					"upgrade should succeed: cluster-api",
				},
				Priority: -1,
			},
			{
				IncludeAll: []string{"service-load-balancer-", "disruption"},
			},
			{
				SIG:        "sig-network-edge",
				IncludeAll: []string{"Application behind service load balancer"},
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
