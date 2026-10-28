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
			// OPP component-first platform-infra suites (ADR 0007): tie-free
			// OCP class; no generic ^lp-interop--OPP-- parent matcher exists
			// because equal-priority multi-claims hard-error conflict resolution.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--OCP--OPP--`), Capabilities: []string{"opp"}},
			// TODO(INTEROP-9629): transitional bridges for currently-unprefixed
			// infra emissions (these exact suite names are emitted only by OPP
			// interop lane steps); remove when the producers emit canonical
			// lp-interop--OCP--OPP--* names.
			{Suite: "opp-restore", Capabilities: []string{"opp"}},
			{Suite: "preflight", Capabilities: []string{"opp"}},
			{Suite: "ocp-upgrade", Capabilities: []string{"opp"}},
			{Suite: "opp-post-upgrade-policy", Capabilities: []string{"opp"}},
			{Suite: "opp-pre-upgrade-checks", Capabilities: []string{"opp"}},
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
