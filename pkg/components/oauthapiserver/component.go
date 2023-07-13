package oauthapiserver

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var OauthApiserverComponent = Component{
	Component: &config.Component{
		Name:                 "oauth-apiserver",
		Operators:            []string{"authentication"},
		DefaultJiraComponent: "oauth-apiserver",
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-oauth-apiserver"},
			},
			{
				Namespaces: []string{
					"openshift-oauth-apiserver",
				},
			},
			{
				IncludeAll: []string{"oauth-api", "-connections"},
			},
			{
				SIG: "sig-auth",
			},
		},
		TestRenames: map[string]string{
			"[sig-auth][Feature:Authentication] TestFrontProxy should succeed [Suite:openshift/conformance/parallel]":                                                                         "[sig-auth][Feature:Authentication] TestFrontProxy should succeed [apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-auth][Feature:OAuthServer] ClientSecretWithPlus should create oauthclient [apigroup:oauth.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]": "[sig-auth][Feature:OAuthServer] ClientSecretWithPlus should create oauthclient [apigroup:config.openshift.io][apigroup:oauth.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-auth][Feature:OAuthServer] OAuth server has the correct token and certificate fallback semantics [apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]":       "[sig-auth][Feature:OAuthServer] OAuth server has the correct token and certificate fallback semantics [apigroup:config.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-auth][Feature:OAuthServer] well-known endpoint should be reachable [apigroup:route.openshift.io] [Suite:openshift/conformance/parallel]":                                    "[sig-auth][Feature:OAuthServer] well-known endpoint should be reachable [apigroup:config.openshift.io][apigroup:route.openshift.io] [Suite:openshift/conformance/parallel]",
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
