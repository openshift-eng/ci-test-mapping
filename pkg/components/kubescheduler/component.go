package kubescheduler

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var KubeSchedulerComponent = Component{
	Component: &config.Component{
		Name:                 "kube-scheduler",
		Operators:            []string{"kube-scheduler"},
		DefaultJiraComponent: "kube-scheduler",
		Namespaces: []string{
			"openshift-kube-scheduler",
			"openshift-kube-scheduler-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-kube-scheduler"},
			},
			{
				SIG:      "sig-scheduling",
				Priority: -20,
			},
			{
				IncludeAny: []string{
					"upgrade should succeed: kube-scheduler",
				},
			},
			{Suite: "Scheduler predicates and priority test suites"},
			{Suite: "Scheduler related scenarios"},
			{Suite: "Testing Scheduler Operator related scenarios"},
			{Suite: "nodeAffinity"},
			{Suite: "podAffinity"},
			{Suite: "resouces related scenarios"},
			{Suite: "taint toleration related scenarios"},
			{Suite: "Scheduler alert related features"},
		},
		TestRenames: map[string]string{
			"[kube-scheduler][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-scheduler":             "[bz-kube-scheduler][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-scheduler",
			"[kube-scheduler][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-scheduler-operator":    "[bz-kube-scheduler][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-scheduler-operator",
			"[kube-scheduler][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-scheduler":          "[bz-kube-scheduler][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-scheduler",
			"[kube-scheduler][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-scheduler-operator": "[bz-kube-scheduler][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-scheduler-operator",
			"[sig-scheduling][Early] The openshift-apiserver pods [apigroup:authorization.openshift.io][apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:project.openshift.io][apigroup:quota.openshift.io][apigroup:route.openshift.io][apigroup:security.openshift.io][apigroup:template.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-apiserver pods [apigroup:apps.openshift.io][apigroup:authorization.openshift.io][apigroup:build.openshift.io][apigroup:image.openshift.io][apigroup:project.openshift.io][apigroup:quota.openshift.io][apigroup:route.openshift.io][apigroup:security.openshift.io][apigroup:template.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
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
