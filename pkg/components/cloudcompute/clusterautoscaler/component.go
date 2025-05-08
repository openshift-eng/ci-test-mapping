package cloudcomputeclusterautoscaler

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ClusterAutoscalerComponent = Component{
	Component: &config.Component{
		Name:                 "Cloud Compute / Cluster Autoscaler",
		Operators:            []string{"cluster-autoscaler"},
		DefaultJiraComponent: "Cluster Autoscaler",
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"[sig-cluster-lifecycle][Feature:Machines] Managed cluster should have machine resources [apigroup:machine.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-cluster-lifecycle][Feature:Machines] Managed cluster should have machine resources [Suite:openshift/conformance/parallel]",
					"[sig-cluster-lifecycle][Feature:Machines][Serial] Managed cluster should grow and decrease when scaling different machineSets simultaneously [Suite:openshift/conformance/serial]",
					"[sig-cluster-lifecycle][Feature:Machines][Serial] Managed cluster should grow and decrease when scaling different machineSets simultaneously [Timeout:20m] [Suite:openshift/conformance/serial]",
					"[sig-cluster-lifecycle][Feature:Machines][Serial] Managed cluster should grow and decrease when scaling different machineSets simultaneously [Timeout:30m][apigroup:machine.openshift.io] [Suite:openshift/conformance/serial]",
					"[sig-cluster-lifecycle][Feature:Machines][Serial] Managed cluster should grow and decrease when scaling different machineSets simultaneously [Timeout:30m] [Suite:openshift/conformance/serial]",
				},
			},
			{
				IncludeAny: []string{
					"[sig-cluster-lifecycle] Cluster_Infrastructure CAS",
					"upgrade should succeed: cluster-autoscaler",
					"OCP-22612",
				},
				Priority: 2,
			},
			{Suite: "Cluster Autoscaler Tests"},
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
