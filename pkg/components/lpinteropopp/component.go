package lpinteropopp

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPinteropOPPComponent = Component{
	Component: &config.Component{
		Name:                 "lp-interop--OPP",
		Operators:            []string{},
		DefaultJiraComponent: "lp-interop--OPP",
		Matchers: []config.ComponentMatcher{
			{Suite: "lp-interop--OPP"},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--interop-opp-`)},

			// PR-B: invisible suite renames
			{Suite: "lp-interop--OPP--smoke"},
			{Suite: "lp-interop--OPP--skip-gate"},
			{Suite: "lp-interop--OPP--acm-app"},

			// PR-C: working suite renames
			{Suite: "lp-interop--OPP--odf-health"},
			{Suite: "lp-interop--OPP--acm-obs-odf"},

			// PR-D: JUnit wrappers
			{Suite: "lp-interop--OPP--install-operators"},
			{Suite: "lp-interop--OPP--interop-opp-backup"},
			{Suite: "lp-interop--OPP--interop-opp-product-upgrade-acm"},
			{Suite: "lp-interop--OPP--interop-opp-product-upgrade-acs"},
			{Suite: "lp-interop--OPP--interop-opp-product-upgrade-odf"},
			{Suite: "lp-interop--OPP--interop-opp-product-upgrade-quay"},
			{Suite: "lp-interop--OPP--interop-opp-wait-mcp"},

			// PR-E: sed post-process renames
			{Suite: "lp-interop--OPP--acs-smoke"},
			{Suite: "lp-interop--OPP--acm-clc"},
			{Suite: "lp-interop--OPP--ocs"},
			{Suite: "lp-interop--OPP--acm-upgrade"},
			{Suite: "lp-interop--OPP--acs-upgrade"},

			// PR-F: per-policy breakout
			{Suite: "lp-interop--OPP--acm-policies"},
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
