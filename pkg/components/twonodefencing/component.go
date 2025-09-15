package twonodefencing

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var TwoNodeFencingComponent = Component{
	Component: &config.Component{
		Name:                 "Two Node Fencing",
		Operators:            []string{},
		DefaultJiraComponent: "Two Node Fencing",
		Matchers:             []config.ComponentMatcher{},
		TestRenames: map[string]string{
			"[sig-etcd][apigroup:config.openshift.io][OCPFeatureGate:DualReplica] Two Node with Fencing should have etcd pods and containers configured correctly":              "[sig-etcd][apigroup:config.openshift.io][OCPFeatureGate:DualReplica][Suite:openshift/two-node] Two Node with Fencing pods and podman containers Should validate the number of etcd pods and containers as configured",
			"[sig-etcd][apigroup:config.openshift.io][OCPFeatureGate:DualReplica] Two Node with Fencing should have podman etcd containers running on each node":                "[sig-etcd][apigroup:config.openshift.io][OCPFeatureGate:DualReplica][Suite:openshift/two-node] Two Node with Fencing pods and podman containers Should verify the number of podman-etcd containers as configured",
			"[sig-node][apigroup:config.openshift.io][OCPFeatureGate:DualReplica] Two Node with Fencing topology should only have two control plane nodes and no arbiter nodes": "[sig-node][apigroup:config.openshift.io][OCPFeatureGate:DualReplica][Suite:openshift/two-node] Two Node with Fencing topology Should validate the number of control-planes, arbiters as configured",
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
