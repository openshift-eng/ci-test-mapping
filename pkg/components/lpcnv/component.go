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
			{Suite: "CNV-lp-interop"},
			// OPP-specific matchers — OPP-first patterns MUST precede
			// the broad CNV matchers below (FindMatch stops at first match).
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--cnv(--|$)`), Capabilities: []string{"OPP"}},
			// ACM+Virt test patterns
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--CNV--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--cnv--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--CNV--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--CNV--`), Capabilities: []string{"ACM+Virt"}},
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
