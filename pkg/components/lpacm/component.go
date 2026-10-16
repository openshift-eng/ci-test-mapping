package lpacm

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var LPacmComponent = Component{
	Component: &config.Component{
		Name:                 "LP--ACM",
		Operators:            []string{},
		DefaultJiraComponent: "LP--ACM",
		Matchers: []config.ComponentMatcher{
			{Suite: "ACM-lp-interop"},
			// OPP-specific matchers — specific OPP-first patterns MUST precede
			// the broad ACM matchers below (FindMatch stops at first match).
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--Observability(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^acm-opp-app$`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--acm-opp-app(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`(?i)^lp-interop--(OPP--)?acm-obs-odf`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--acm-alc(--|$)`), Capabilities: []string{"OPP"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--acm-grc(--|$)`), Capabilities: []string{"OPP"}},
			// ACM+Virt test patterns
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM-Virt--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--acm-virt`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--ACM--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--ACM--`), Capabilities: []string{"ACM+Virt"}},
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
