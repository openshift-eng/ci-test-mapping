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
			// OPP component-first platform-infra suites (ADR 0007): canonical
			// OCP-class suites route to this component because no LP--OCP
			// component is registered; if one is ever created, MOVE this matcher
			// here (never copy it). Do not add a generic ^lp-interop--OPP-- parent
			// matcher: it would tie at Priority 0 with lpodf's
			// ^lp-interop--OPP--pytest (28 live rows) and lpacm's
			// ^lp-interop--OPP--acm, and equal-priority multi-claims hard-error
			// conflict resolution (pkg/components/component.go getHighestPriority).
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--OPP--`), Capabilities: []string{"opp"}},
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
