package etcd

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var EtcdComponent = Component{
	Component: &config.Component{
		Name:                 "Etcd",
		Operators:            []string{"etcd"},
		DefaultJiraComponent: "Etcd",
		Namespaces: []string{
			"openshift-etcd",
			"openshift-etcd-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG: "sig-etcd",
			},
			{
				IncludeAll: []string{"bz-etcd"},
			},
			{
				IncludeAll: []string{"bz-Etcd"},
			},
			{
				IncludeAll: []string{"cluster-etcd-operator"},
			},
			{
				IncludeAny: []string{
					":Etcd ",
					"etcd-operator and cluster works well after upgrade",
					"upgrade should succeed: etcd",
				},
				Priority: 1,
			},
			{Suite: "DR_Testing"},
			{Suite: "ETCD"},
			{Suite: "etcd related features"},
		},
		TestRenames: map[string]string{
			"[bz-etcd] pathological event should not see excessive RequiredInstallerResourcesMissing secrets":         "[bz-etcd] should not see excessive RequiredInstallerResourcesMissing secrets",
			"[Etcd][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-etcd":             "[bz-Etcd][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-etcd",
			"[Etcd][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-etcd-operator":    "[bz-Etcd][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-etcd-operator",
			"[Etcd][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-etcd":          "[bz-Etcd][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-etcd",
			"[Etcd][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-etcd-operator": "[bz-Etcd][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-etcd-operator",

			"[sig-scheduling][Early] The openshift-etcd pods [apigroup:operator.openshift.io] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The openshift-etcd pods [apigroup:operator.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
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
