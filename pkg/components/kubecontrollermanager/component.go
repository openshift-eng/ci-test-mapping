package kubecontrollermanager

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var KubeControllerManagerComponent = Component{
	Component: &config.Component{
		Name:                 "kube-controller-manager",
		Operators:            []string{"kube-controller-manager"},
		DefaultJiraComponent: "kube-controller-manager",
		Namespaces: []string{
			"openshift-kube-controller-manager",
			"openshift-kube-controller-manager-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-kube-controller-manager"},
			},
			{
				IncludeAny: []string{
					"Feature:ClusterResourceQuota",
					"ResourceQuota",
					"upgrade should succeed: kube-controller-manager",
				},
				Priority: 1, // quota is owned by KCM more strongly than apimachinery
			},
			{Suite: "Features of daemonset"},
			{Suite: "Quota related scenarios"},
			{Suite: "Testing kube-controller-manager-operator"},
			{Suite: "cronjob related"},
			{Suite: "deployment related features"},
			{Suite: "deployment/dc related features via web"},
			{Suite: "job.feature"},
			{Suite: "quota related feature"},
			{Suite: "Features of daeXXXset"},
			{Suite: "deployment related steps"},
		},
		TestRenames: map[string]string{
			"[sig-api-machinery][Feature:ClusterResourceQuota] Cluster resource quota should control resource limits across namespaces [apigroup:quota.openshift.io][apigroup:image.openshift.io] [Suite:openshift/conformance/parallel]": "[sig-api-machinery][Feature:ClusterResourceQuota] Cluster resource quota should control resource limits across namespaces [apigroup:quota.openshift.io][apigroup:image.openshift.io][apigroup:monitoring.coreos.com][apigroup:template.openshift.io] [Suite:openshift/conformance/parallel]",
			"[kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-controller-manager":                                                                                           "[bz-kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-controller-manager",
			"[kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-controller-manager-operator":                                                                                  "[bz-kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-controller-manager-operator",
			"[kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-controller-manager":                                                                                        "[bz-kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-controller-manager",
			"[kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-controller-manager-operator":                                                                               "[bz-kube-controller-manager][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-controller-manager-operator",
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
