package monitoring

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var MonitoringComponent = Component{
	Component: &config.Component{
		Name:                 "Monitoring",
		Operators:            []string{"monitoring"},
		DefaultJiraComponent: "Monitoring",
		Namespaces: []string{
			"openshift-monitoring",
			"openshift-user-workload-monitoring",
		},
		TestRenames: map[string]string{
			"[sig-instrumentation][Late] Alerts shouldn't report any alerts in firing or pending state apart from Watchdog and AlertmanagerReceiversNotConfigured and have no gaps in Watchdog firing [Skipped:Disconnected] [Suite:openshift/conformance/parallel]":  "[sig-instrumentation][Late] Alerts shouldn't report any alerts in firing state apart from Watchdog and AlertmanagerReceiversNotConfigured [Suite:openshift/conformance/parallel]",
			"[sig-instrumentation][sig-builds][Feature:Builds] Prometheus when installed on the cluster should start and expose a secured proxy and verify build metrics [apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]": "[sig-instrumentation][sig-builds][Feature:Builds] Prometheus when installed on the cluster should start and expose a secured proxy and verify build metrics [apigroup:config.openshift.io][apigroup:build.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
			"[sig-instrumentation][Late] Alerts shouldn't exceed the series limit of total series sent via telemetry from each cluster [Suite:openshift/conformance/parallel]":                                                                                        "[sig-instrumentation][Late] Alerts shouldn't exceed the 650 series limit of total series sent via telemetry from each cluster [Suite:openshift/conformance/parallel]",
			"[Monitoring][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-monitoring":                                                                                                                                                 "[bz-Monitoring][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-monitoring",
			"[Monitoring][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-monitoring":                                                                                                                                              "[bz-Monitoring][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-monitoring",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-user-workload-monitoring":                                                                                                                                      "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-user-workload-monitoring",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-user-workload-monitoring":                                                                                                                                   "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-user-workload-monitoring",

			"[sig-scheduling][Early] The openshift-monitoring prometheus-adapter pods [apigroup:monitoring.coreos.com] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-monitoring prometheus-adapter pods [apigroup:monitoring.coreos.com] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
			"[sig-scheduling][Early] The openshift-monitoring thanos-querier pods [apigroup:monitoring.coreos.com] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]":     "[sig-scheduling][Early] The openshift-monitoring thanos-querier pods [apigroup:monitoring.coreos.com] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-Monitoring"},
			},
			{
				IncludeAll: []string{"bz-monitoring"},
			},
			{
				SIG:      "sig-instrumentation",
				Priority: 1,
			},
			{
				IncludeAll: []string{"cluster-monitoring-operator"},
			},
			{Suite: "Cluster_Observability", Priority: 1},
			{Suite: "Symptom Detection"},
			{Suite: "cluster monitoring related upgrade check"},
			{Suite: "servicemonitor"},
			{Suite: "serviceXXXitor"},
			{Suite: "alerts browser"},
			{Suite: "query browser"},
			{Suite: "Install and configuration related scenarios"},
			{Suite: "Monitoring dashboards related features"},
			{
				IncludeAny: []string{
					"upgrade should succeed: monitoring",
				},
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
