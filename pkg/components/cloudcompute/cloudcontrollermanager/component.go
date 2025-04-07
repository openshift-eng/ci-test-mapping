package cloudcomputecloudcontrollermanager

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var CloudControllerManagerComponent = Component{
	Component: &config.Component{
		Name:                 "Cloud Compute / Cloud Controller Manager",
		Operators:            []string{"cloud-controller-manager"},
		DefaultJiraComponent: "Cloud Compute / Cloud Controller Manager",
		Namespaces: []string{
			"openshift-cloud-controller-manager",
			"openshift-cloud-controller-manager-operator",
		},
		TestRenames: map[string]string{
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cloud-controller-manager":             "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cloud-controller-manager",
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cloud-controller-manager-operator":    "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-cloud-controller-manager-operator",
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cloud-controller-manager":          "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cloud-controller-manager",
			"[Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cloud-controller-manager-operator": "[bz-Cloud Compute][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-cloud-controller-manager-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-cloud-controller-manager"},
			},
			{
				IncludeAny: []string{
					"[sig-cloud-provider][Feature:OpenShiftCloudControllerManager][Late] Deploy an external cloud provider [apigroup:config.openshift.io][apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-cloud-provider][Feature:OpenShiftCloudControllerManager][Late] Deploy an external cloud provider [apigroup:machineconfiguration.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-cloud-provider][Feature:OpenShiftCloudControllerManager][Late] Deploy an external cloud provider [Suite:openshift/conformance/parallel]",
					"[sig-cluster-lifecycle] CSRs from machines that are not recognized by the cloud provider are not approved [apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]",
					"[sig-cluster-lifecycle] CSRs from machines that are not recognized by the cloud provider are not approved [Suite:openshift/conformance/parallel]",
				},
			},
			{
				IncludeAny: []string{
					"[sig-cluster-lifecycle] Cluster_Infrastructure CCM",
					"upgrade should succeed: cloud-controller-manager",
					"OCP-41804",
					"OCP-64657",
				},
				Priority: 2,
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
