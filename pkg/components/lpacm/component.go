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
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--ACM--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--ACM--`)},
			//  LP Interop Test Suites: Capability is used to group Runtime Env.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM--ACM\+Virt--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM-Virt--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM--OPP--`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM--`), Capabilities: []string{"ACM+Virt"}},
			// OPP test patterns with opp capability
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--Observability`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^acm-opp-app$`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OPP--acm`), Capabilities: []string{"opp"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM-OBS-ODF`), Capabilities: []string{"opp"}},
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
