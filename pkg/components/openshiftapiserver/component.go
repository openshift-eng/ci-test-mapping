package openshiftapiserver

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OpenshiftApiserverComponent = Component{
	Component: &config.Component{
		Name:                 "openshift-apiserver",
		Operators:            []string{"openshift-apiserver", "openshift-lifecycle-manager", "openshift-lifecycle-manager-catalog", "openshift-lifecycle-manager-packageserver", "platform-operators-aggregated"},
		DefaultJiraComponent: "openshift-apiserver",
		Namespaces: []string{
			"openshift-apiserver",
			"openshift-apiserver-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-openshift-apiserver"},
			},
			{
				IncludeAll: []string{"openshift-api-", "-connections"},
			},
			{Suite: "Projects"},
			{Suite: "project list tests"},
			{Suite: "project permissions"},
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
