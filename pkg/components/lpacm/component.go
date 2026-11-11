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
			// Legacy flat suite, kept solely because maintainer direction
			// (etirta, PR 908 comment 2026-10-30 22:34Z) named this one entry
			// for restoration; every other legacy LP matcher stays removed.
			// It carries no rows today and cannot acquire any: "ACM-lp-interop"
			// has 0 rows in both committed snapshots (ci_analysis_us,
			// ci_analysis_qe) at this head and at the merge base, and the US
			// ingestion filter never admits the name (it is absent from
			// includeSuites and matches none of includeSuitePatterns), so the
			// matcher is inert here and map-verify is unaffected either way.
			// Declared intent only; retirable when the maintainer says so
			// (INTEROP-9629 lineage).
			{Suite: "ACM-lp-interop"},
			{SuiteRegEx: regexp.MustCompile(`^lp-ocp-compat--ACM--`)},
			{SuiteRegEx: regexp.MustCompile(`^lp-chaos--ACM--`)},
			//  LP Interop Test Suites: Capability is used to group Runtime Env.
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM--ACM\+Virt--`), Capabilities: []string{"ACM+Virt"}},
			{SuiteRegEx: regexp.MustCompile(`^lp-interop--ACM--OPP--`), Capabilities: []string{"OPP"}},
			// Runtime-first emission from interop-opp-observability-odf, which
			// still writes the literal suite name in openshift/release (the
			// canonical matcher above does not cover it). Restored so the live
			// acm-obs-odf rows stay mapped instead of falling to Unknown; drop
			// together with the exclusion narrowing once that step migrates
			// (INTEROP-9629).
			{Suite: "lp-interop--OPP--acm-obs-odf", Capabilities: []string{"OPP"}},
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
