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
		Name:      "LP--OCP",
		Operators: []string{},
		// LP--OCP has no Jira component of its own; JIRAComponentID is
		// intentionally null in the generated mapping.  These matchers
		// exist solely to assign suites to the LP--OCP Component
		// Readiness bucket with the "OPP" capability.
		DefaultJiraComponent: "LP--OCP",
		Matchers: []config.ComponentMatcher{
			// BEWARE: FindMatch stops at first match — specific matchers
			// MUST precede broader ones.  ocp-upgrade-healthcheck must
			// come before the ocp-upgrade rule.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--ocp-upgrade-healthcheck(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--ocp-upgrade(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--smoke(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--backup(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--wait-mcp(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--skip-gate(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--preflight(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--acm-clc(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--ocs(--|$)`), Capabilities: []string{"OPP"}},
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
