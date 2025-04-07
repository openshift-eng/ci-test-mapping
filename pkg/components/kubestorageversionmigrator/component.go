package kubestorageversionmigrator

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var KubeStorageVersionMigratorComponent = Component{
	Component: &config.Component{
		Name:                 "kube-storage-version-migrator",
		Operators:            []string{"kube-storage-version-migrator"},
		DefaultJiraComponent: "kube-storage-version-migrator",
		Namespaces: []string{
			"openshift-kube-storage-version-migrator",
			"openshift-kube-storage-version-migrator-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-kube-storage-version-migrator"},
			},
			{
				IncludeAny: []string{
					"upgrade should succeed: kube-storage-version-migrator",
				},
			},
		},
		TestRenames: map[string]string{
			"[kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-storage-version-migrator":             "[bz-kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-storage-version-migrator",
			"[kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-storage-version-migrator-operator":    "[bz-kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-storage-version-migrator-operator",
			"[kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-storage-version-migrator":          "[bz-kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-storage-version-migrator",
			"[kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-storage-version-migrator-operator": "[bz-kube-storage-version-migrator][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-storage-version-migrator-operator",
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
