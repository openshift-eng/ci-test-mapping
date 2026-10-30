package lpcnv

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPcnvComponent = Component{
	Component: &config.Component{
		Name:                 "LP--CNV",
		Operators:            []string{},
		DefaultJiraComponent: "LP--CNV",
		Matchers: []config.ComponentMatcher{
			// Legacy flat suite, restored per maintainer direction (etirta,
			// 2026-10-30): map-verify validates against the periodically
			// refreshed junit.json snapshot, which still carries live rows with
			// this exact suite name. Do not remove until the flat emissions
			// retire (INTEROP-9629 lineage).
			{Suite: "CNV-lp-interop"},
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--CNV--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--CNV--`)},
			//  LP Interop Test Suites: Capability is used to group Runtime Env.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--CNV--ACM\+Virt--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--CNV--OPP--`), Capabilities: []string{"OPP"}},
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
