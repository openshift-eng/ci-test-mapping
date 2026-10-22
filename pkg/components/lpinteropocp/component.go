package lpinteropocp

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPinteropOCPComponent = Component{
	Component: &config.Component{
		Name:                 "LP--OCP",
		Operators:            []string{},
		DefaultJiraComponent: "LP--OCP",
		Matchers: []config.ComponentMatcher{
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--smoke`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--backup`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--wait-mcp`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--ocp-upgrade`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--ocp-upgrade-healthcheck`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--skip-gate`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--preflight`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--acm-clc`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--ocs`), Capabilities: []string{"opp"}},
			// Legacy OPP-prefixed patterns
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--smoke`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--backup`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--wait-mcp`), Capabilities: []string{"opp"}},
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
