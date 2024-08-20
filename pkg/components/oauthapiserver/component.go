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
		Namespaces: []string{
			"openshift-oauth-apiserver",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAll: []string{"bz-oauth-apiserver"},
			},
			{
				IncludeAll: []string{"oauth-api", "-connections"},
			},
			{
				SIG:      "sig-auth",
				Priority: -1,
			},
		},
		TestRenames: map[string]string{
			"[sig-auth][Feature:Authentication] TestFrontProxy should succeed [Suite:openshift/conformance/parallel]":                                                                                 "[sig-auth][Feature:Authentication] TestFrontProxy should succeed [apigroup:config.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-auth][Feature:OAuthServer] ClientSecretWithPlus should create oauthclient [apigroup:oauth.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]":         "[sig-auth][Feature:OAuthServer] ClientSecretWithPlus should create oauthclient [apigroup:config.openshift.io][apigroup:oauth.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-auth][Feature:OAuthServer] OAuth server has the correct token and certificate fallback semantics [apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]":               "[sig-auth][Feature:OAuthServer] OAuth server has the correct token and certificate fallback semantics [apigroup:config.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/parallel]",
			"[sig-auth][Feature:OAuthServer] well-known endpoint should be reachable [apigroup:route.openshift.io] [Suite:openshift/conformance/parallel]":                                            "[sig-auth][Feature:OAuthServer] well-known endpoint should be reachable [apigroup:config.openshift.io][apigroup:route.openshift.io] [Suite:openshift/conformance/parallel]",
			"[oauth-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-oauth-apiserver":                                                                       "[bz-oauth-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-oauth-apiserver",
			"[oauth-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-oauth-apiserver":                                                                    "[bz-oauth-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-oauth-apiserver",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-oauth-api connection/new":                                                                       "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-oauth-api connection/new target=api-int":                                                        "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-oauth-api connection/new target=service-network":                                                "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/service-network/cache-oauth-api connection/new":                                                       "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-oauth-api connection/reused":                                                                    "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-oauth-api connection/reused target=api-int":                                                     "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-oauth-api connection/reused target=service-network":                                             "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/service-network/cache-oauth-api connection/reused":                                                    "[sig-network] there should be nearly zero single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-oauth-api connection/new":                                                                    "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-oauth-api connection/new target=api-int":                                                     "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-oauth-api connection/new target=service-network":                                             "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/service-network/cache-oauth-api connection/new":                                                    "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-oauth-api connection/reused":                                                                 "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-oauth-api connection/reused target=api-int":                                                  "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-oauth-api connection/reused target=service-network":                                          "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/service-network/cache-oauth-api connection/reused":                                                 "[sig-network] there should be reasonably few single second disruptions for cache-oauth-api-reused-connections",
			"[Serial] [sig-auth][Feature:OAuthServer] [RequestHeaders] [IdP] test RequestHeaders IdP [apigroup:config.openshift.io][apigroup:user.openshift.io] [Suite:openshift/conformance/serial]": "[Serial] [sig-auth][Feature:OAuthServer] [RequestHeaders] [IdP] test RequestHeaders IdP [apigroup:config.openshift.io][apigroup:user.openshift.io][apigroup:apps.openshift.io] [Suite:openshift/conformance/serial]",
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
