package networkingrouter

import (
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
)

type Component struct {
	*config.Component
}

var RouterComponent = Component{
	Component: &config.Component{
		Name:                 "Networking / router",
		Operators:            []string{"ingress"},
		DefaultJiraComponent: "Networking / router",
		Namespaces: []string{
			"openshift-ingress",
			"openshift-ingress-canary",
			"openshift-ingress-operator",
		},
		Matchers: []config.ComponentMatcher{
			{
				SIG:        "sig-network",
				IncludeAll: []string{"Feature:Router"},
			},
			{
				SIG:        "sig-network-edge",
				IncludeAll: []string{"Feature:Router"},
			},
			{
				IncludeAll: []string{"single second disruption", "ingress-to-"},
				Priority:   20,
			},
			{
				IncludeAny: []string{
					"bz-Routing",
					"openshift-ingress",
					"via cluster ingress",
					"Cluster frontend ingress", "[sig-arch] Managed cluster should [apigroup:apps.openshift.io] should expose cluster services outside the cluster [apigroup:route.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
					"[sig-arch] Managed cluster should expose cluster services outside the cluster [apigroup:route.openshift.io] [Skipped:Disconnected] [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Idling with a single service and ReplicationController should idle the service and ReplicationController properly [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many UDP senders (by continuing to drop all packets on the floor) [Serial] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should work with TCP (when fully idled) [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should work with TCP (when fully idled) [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should work with UDP [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling should handle many UDP senders (by continuing to drop all packets on the floor) [Serial] [Suite:openshift/conformance/serial]",
					"[sig-network-edge][Feature:Idling] Unidling should work with TCP (when fully idled) [Skipped:Network/OVNKubernetes] [Suite:openshift/conformance/parallel]",
					"[sig-network-edge][Feature:Idling] Unidling should work with UDP [Suite:openshift/conformance/parallel]",
					"[sig-networking] should not see excessive FailedToUpdateEndpointSlices Error updating Endpoint Slices",
					"upgrade should succeed: ingress",
				},
			},
			{Suite: "Network_Edge"},
			{Suite: "Routing and DNS related scenarios"},
			{Suite: "Test Ingress API logging options"},
			{Suite: "Testing HTTP Headers related scenarios"},
			{Suite: "Testing Ingress Operator related scenarios"},
			{Suite: "Testing haproxy rate limit related features"},
			{Suite: "Testing haproxy router"},
			{Suite: "Testing ingress to route object"},
			{Suite: "Testing route"},
			{Suite: "Testing timeout route"},
			{Suite: "Testing websocket features"},
			{Suite: "Testing wildcard routes"},
			{Suite: "route related"},
			{Suite: "Testing abrouting"},
		},
		TestRenames: map[string]string{
			"[Routing][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ingress":                                                             "[bz-Routing][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ingress",
			"[Routing][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ingress-canary":                                                      "[bz-Routing][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ingress-canary",
			"[Routing][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ingress-operator":                                                    "[bz-Routing][invariant] alert/KubePodNotReady should not be at or above info in ns/openshift-ingress-operator",
			"[Routing][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ingress":                                                          "[bz-Routing][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ingress",
			"[Routing][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ingress-canary":                                                   "[bz-Routing][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ingress-canary",
			"[Routing][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ingress-operator":                                                 "[bz-Routing][invariant] alert/KubePodNotReady should not be at or above pending in ns/openshift-ingress-operator",
			"[sig-network-edge] disruption/service-load-balancer-with-pdb connection/new should be available throughout the test":                                           "[sig-network-edge] Application behind service load balancer with PDB remains available using new connections",
			"[sig-network-edge] disruption/service-load-balancer-with-pdb connection/reused should be available throughout the test":                                        "[sig-network-edge] Application behind service load balancer with PDB remains available using reused connections",
			"[sig-network] there should be nearly zero single second disruptions for ns/openshift-console route/console disruption/ingress-to-console connection/new":       "[sig-network] there should be nearly zero single second disruptions for ingress-to-console-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for ns/openshift-console route/console disruption/ingress-to-console connection/reused":    "[sig-network] there should be nearly zero single second disruptions for ingress-to-console-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for ns/openshift-console route/console disruption/ingress-to-console connection/new":    "[sig-network] there should be reasonably few single second disruptions for ingress-to-console-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for ns/openshift-console route/console disruption/ingress-to-console connection/reused": "[sig-network] there should be reasonably few single second disruptions for ingress-to-console-reused-connections",

			"[sig-network] there should be nearly zero single second disruptions for ns/openshift-authentication route/oauth-openshift disruption/ingress-to-oauth-server connection/new":       "[sig-network] there should be nearly zero single second disruptions for ingress-to-oauth-server-new-connections",
			"[sig-network] there should be nearly zero single second disruptions for ns/openshift-authentication route/oauth-openshift disruption/ingress-to-oauth-server connection/reused":    "[sig-network] there should be nearly zero single second disruptions for ingress-to-oauth-server-reused-connections",
			"[sig-network] there should be reasonably few single second disruptions for ns/openshift-authentication route/oauth-openshift disruption/ingress-to-oauth-server connection/new":    "[sig-network] there should be reasonably few single second disruptions for ingress-to-oauth-server-new-connections",
			"[sig-network] there should be reasonably few single second disruptions for ns/openshift-authentication route/oauth-openshift disruption/ingress-to-oauth-server connection/reused": "[sig-network] there should be reasonably few single second disruptions for ingress-to-oauth-server-reused-connections",

			"[sig-scheduling][Early] The HAProxy router pods [apigroup:route.openshift.io] should be scheduled on different nodes [Skipped:SingleReplicaTopology] [Suite:openshift/conformance/parallel]": "[sig-scheduling][Early] The HAProxy router pods [apigroup:route.openshift.io] should be scheduled on different nodes [Suite:openshift/conformance/parallel]",
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
