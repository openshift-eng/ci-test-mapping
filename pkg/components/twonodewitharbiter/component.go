package twonodewitharbiter

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var TwoNodeWithArbiterComponent = Component{
	Component: &config.Component{
		Name:                 "Two Node with Arbiter",
		Operators:            []string{},
		DefaultJiraComponent: "Two Node with Arbiter",
		Matchers:             []config.ComponentMatcher{},
		TestRenames: map[string]string{
			"[sig-apps][apigroup:apps.openshift.io][OCPFeatureGate:HighlyAvailableArbiter] Deployments on HighlyAvailableArbiterMode topology should be created on master nodes when no node selected [Suite:openshift/conformance/parallel]":                    "[sig-apps][apigroup:apps.openshift.io] Deployments on HighlyAvailableArbiterMode topology should be created on master nodes when no node selected [Suite:openshift/conformance/parallel]",
			"[sig-apps][apigroup:apps.openshift.io][OCPFeatureGate:HighlyAvailableArbiter] Evaluate DaemonSet placement in HighlyAvailableArbiterMode topology should not create a DaemonSet on the Arbiter node [Suite:openshift/conformance/parallel]":         "[sig-apps][apigroup:apps.openshift.io] Evaluate DaemonSet placement in HighlyAvailableArbiterMode topology should not create a DaemonSet on the Arbiter node [Suite:openshift/conformance/parallel]",
			"[sig-apps][apigroup:apps.openshift.io][OCPFeatureGate:HighlyAvailableArbiter] Deployments on HighlyAvailableArbiterMode topology should be created on arbiter nodes when arbiter node is selected [Suite:openshift/conformance/parallel]":           "[sig-apps][apigroup:apps.openshift.io] Deployments on HighlyAvailableArbiterMode topology should be created on arbiter nodes when arbiter node is selected [Suite:openshift/conformance/parallel]",
			"[sig-etcd][apigroup:config.openshift.io][OCPFeatureGate:HighlyAvailableArbiter] Ensure etcd health and quorum in HighlyAvailableArbiterMode should have all etcd pods running and quorum met [Suite:openshift/conformance/parallel]":                "[sig-etcd][apigroup:config.openshift.io] Ensure etcd health and quorum in HighlyAvailableArbiterMode should have all etcd pods running and quorum met [Suite:openshift/conformance/parallel]",
			"[sig-node][apigroup:config.openshift.io][OCPFeatureGate:HighlyAvailableArbiter] expected Master and Arbiter node counts Should validate that there are Master and Arbiter nodes as specified in the cluster [Suite:openshift/conformance/parallel]": "[sig-node][apigroup:config.openshift.io] expected Master and Arbiter node counts Should validate that there are Master and Arbiter nodes as specified in the cluster [Suite:openshift/conformance/parallel]",
			"[sig-node][apigroup:config.openshift.io][OCPFeatureGate:HighlyAvailableArbiter] required pods on the Arbiter node Should verify that the correct number of pods are running on the Arbiter node [Suite:openshift/conformance/parallel]":             "[sig-node][apigroup:config.openshift.io] required pods on the Arbiter node Should verify that the correct number of pods are running on the Arbiter node [Suite:openshift/conformance/parallel]",
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
