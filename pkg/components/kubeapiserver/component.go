package kubeapiserver

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var KubeApiserverComponent = Component{
	Component: &config.Component{
		Name:                 "kube-apiserver",
		Operators:            []string{"kube-apiserver"},
		DefaultJiraComponent: "kube-apiserver",
		Namespaces: []string{
			"default",
			"openshift",
			"kube-system",
			"openshift-config",
			"openshift-config-managed",
			"openshift-kube-apiserver",
			"openshift-kube-apiserver-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				IncludeAny: []string{
					"bz-kube-apiserver",
					"cache-kube-api-",
					"[sig-api-machinery][Feature:APIServer]",
					"should have a status in the CRD schema", // only observed failures are kube-apiserver availability.
					"API LBs follow /readyz of kube-apiserver",
					"clients should not use APIs that are removed in upcoming releases",
					"operators should not create watch channels very often",
					"upgrade should succeed: kube-apiserver",
				},
			},
			{
				IncludeAll: []string{"single second disruption", "kube-api"},
			},
			{
				IncludeAll: []string{"kube-api-", "-connections"},
			},
			{
				SIG:      "sig-api-machinery",
				Priority: -1,
			},
			{Suite: "API_Server"},
			{Suite: "CRD related"},
			{Suite: "Event related scenarios"},
			{Suite: "Roles and RoleBindings tests"},
			{Suite: "ServiceAccount and Policy Managerment"},
			{Suite: "events and logs related"},
			{Suite: "limit range related scenarios:"},
			{Suite: "test master config related steps"},
			{Suite: "REST policy related features"},
			// nolint:misspell
			{Suite: "senarios for checking transfer scheme"},
			{Suite: "Api proxy related cases"},
			{Suite: "REST related features"},
			{Suite: "REST features"},
			{Suite: "KUBE API server related features"},
			{
				SIG:          "sig-arch",
				IncludeAll:   []string{"tls artifacts"},
				Capabilities: []string{"TLS"},
				Priority:     -1,
			},
			{
				SIG:          "sig-node",
				IncludeAll:   []string{"TargetDown should not be at or above info in ns/kube-system"},
				Capabilities: []string{"Alerts"},
				Priority:     -10,
			},
		},
		TestRenames: map[string]string{
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/default":                                                       "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/default",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/kube-system":                                                   "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/kube-system",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift":                                                     "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-config":                                              "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-config",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-config-managed":                                      "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-config-managed",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/default":                                                    "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/default",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/kube-system":                                                "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/kube-system",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift":                                                  "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-config":                                           "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-config",
			"[Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-config-managed":                                   "[bz-Unknown][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-config-managed",
			"[kube-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-apiserver":                               "[bz-kube-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-apiserver",
			"[kube-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-apiserver-operator":                      "[bz-kube-apiserver][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-kube-apiserver-operator",
			"[kube-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-apiserver":                            "[bz-kube-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-apiserver",
			"[kube-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-apiserver-operator":                   "[bz-kube-apiserver][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-kube-apiserver-operator",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-kube-api connection/new":                              "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-kube-api connection/new target=api-int":               "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-kube-api connection/new target=service-network":       "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/service-network/cache-kube-api connection/new":              "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-kube-api connection/reused":                           "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-kube-api connection/reused target=api-int":            "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/cache-kube-api connection/reused target=service-network":    "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be nearly zero single second disruptions for disruption/service-network/cache-kube-api connection/reused":           "[sig-network] there should be nearly zero single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-kube-api connection/new":                           "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-kube-api connection/new target=api-int":            "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-kube-api connection/new target=service-network":    "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/service-network/cache-kube-api connection/new":           "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-kube-api connection/reused":                        "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-kube-api connection/reused target=api-int":         "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/cache-kube-api connection/reused target=service-network": "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for disruption/service-network/cache-kube-api connection/reused":        "[sig-network] there should be reasonably few single second disruptions for cache-kube-api-reused-connections",
			"[sig-arch][Late] collect certificate data [Suite:openshift/conformance/parallel]":                                                              "[sig-arch][Late][Jira:\"kube-apiserver\"] collect certificate data [Suite:openshift/conformance/parallel]",
			"[sig-arch][Late] all tls artifacts must be registered [Suite:openshift/conformance/parallel]":                                                  "[sig-arch][Late][Jira:\"kube-apiserver\"] all tls artifacts must be registered [Suite:openshift/conformance/parallel]",
			"[sig-arch][Late] all registered tls artifacts must have no metadata violation regressions [Suite:openshift/conformance/parallel]":              "[sig-arch][Late][Jira:\"kube-apiserver\"] all registered tls artifacts must have no metadata violation regressions [Suite:openshift/conformance/parallel]",
			"[sig-arch][Late] operators should not create watch channels very often":                                                                        "[sig-arch][Late] operators should not create watch channels very often [apigroup:apiserver.openshift.io] [Suite:openshift/conformance/parallel]",
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
