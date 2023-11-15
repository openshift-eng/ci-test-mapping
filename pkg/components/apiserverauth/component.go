package apiserverauth

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var ApiserverAuthComponent = Component{
	Component: &config.Component{
		Name:                 "apiserver-auth",
		Operators:            []string{},
		DefaultJiraComponent: "apiserver-auth",
		Namespaces: []string{
			"openshift-authentication",
			"openshift-authentication-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-apiserver-auth"},
			},
			{Suite: "Authentication"},
			{Suite: "Group sync related scenarios"},
			{Suite: "SCC policy related scenarios"},
			{Suite: "Seccomp part of SCC policy should be kept and working after upgrade"},
			{Suite: "Security_and_Compliance"},
			{Suite: "User management related"},
			{Suite: "apiserver and auth related upgrade check"},
			{Suite: "auth prometheus metrics feature"},
			{Suite: "change the policy of user/service account"},
			{Suite: "groups and users related features"},
			{Suite: "idp feature"},
			{Suite: "permission related test"},
			{Suite: "token feature"},
			{Suite: "certificates related scenarios"},
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
