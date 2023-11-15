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
			{Suite: "DR_Testing"},
			{Suite: "ETCD"},
			{Suite: "etcd related features"},
		},
		TestRenames: map[string]string{
			"[bz-etcd] pathological event should not see excessive RequiredInstallerResourcesMissing secrets": "[bz-etcd] should not see excessive RequiredInstallerResourcesMissing secrets",
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
