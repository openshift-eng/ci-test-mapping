package testframework

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var TestFrameworkComponent = Component{
	Component: &config.Component{
		Name:                 "Test Framework",
		Operators:            []string{},
		DefaultJiraComponent: "Test Framework",
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"Undiagnosed panic"},
				Priority:   1,
			},
			{
				SIG: "sig-trt",
			},
			{
				SIG: "sig-ci",
			},
			{
				IncludeAll:   []string{"ci-cluster-network-liveness-", "-connections"},
				Capabilities: []string{"Build Clusters"},
			},

			// TRT is owner of last resort for these
			{
				SIG:          "sig-arch",
				IncludeAll:   []string{"events should not repeat"},
				Capabilities: []string{"Pathological Events"},
				Priority:     -10,
			},
			{
				SIG:          "sig-arch",
				IncludeAll:   []string{"Alerts alert/"},
				Capabilities: []string{"Alerts"},
				Priority:     -10,
			},

			{
				IncludeAll:   []string{"bz-Unknown", "alert/KubePodNotReady"},
				Priority:     -10,
				Capabilities: []string{"KubePodNotReady - Other"},
			},
		},
		TestRenames: map[string]string{
			"[sig-arch] pathological event should not see excessive pull back-off on registry.redhat.io":                                                                                                                                                                                              "[sig-arch] should not see excessive pull back-off on registry.redhat.io",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in all the other namespaces":                                                                                                                                                                                   "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in all the other namespaces",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in all the other namespaces":                                                                                                                                                                                "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in all the other namespaces",
			"[Jira:\"Test Framework\"] Prometheus [apigroup:config.openshift.io] when installed on the cluster shouldn't have failing rules evaluation [Suite:openshift/conformance/parallel]":                                                                                                        "[sig-instrumentation] Prometheus [apigroup:image.openshift.io] when installed on the cluster shouldn't have failing rules evaluation [Suite:openshift/conformance/parallel]",
			"[Jira:\"Test Framework\"] Prometheus [apigroup:coreos.monitoring.com] when installed on the cluster shouldn't report any alerts in firing state apart from Watchdog and AlertmanagerReceiversNotConfigured [Early][apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]": "[sig-instrumentation] Prometheus [apigroup:image.openshift.io] when installed on the cluster should have a AlertmanagerReceiversNotConfigured alert in firing state [Suite:openshift/conformance/parallel]",
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
