package networkingdns

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var DNSComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / DNS",
		Operators:            []string{"dns"},
		DefaultJiraComponent: "Networking / DNS",
		Namespaces: []string{
			"openshift-dns",
			"openshift-dns-operator",
		},
		TestRenames: map[string]string{
			"[DNS][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-dns":             "[bz-DNS][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-dns",
			"[DNS][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-dns-operator":    "[bz-DNS][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-dns-operator",
			"[DNS][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-dns":          "[bz-DNS][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-dns",
			"[DNS][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-dns-operator": "[bz-DNS][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-dns-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG:        "sig-network",
				IncludeAll: []string{"DNS"},
			},
			{
				IncludeAll: []string{"cluster-dns-operator"},
			},
			{
				IncludeAll: []string{"bz-DNS"},
			},
			{
				SIG:        "sig-network-edge",
				IncludeAll: []string{"DNS"},
				Priority:   1,
			},
			{Suite: "Testing DNS features"},
			{
				IncludeAny: []string{
					"upgrade should succeed: dns",
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
