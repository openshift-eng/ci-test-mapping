package lpcomplpinterop

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPCompLpInteropComponent = Component{
	Component: &config.Component{
		Name:                 "LPComp-lp-interop",
		Operators:            []string{},
		DefaultJiraComponent: "LPComp",
		Matchers: []config.ComponentMatcher{
			{Suite: "LPComp-lp-interop"},
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--LPComp--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--LPComp--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--LPComp--`)},
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
