package config

import (
	"reflect"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
)

func TestComponent_FindMatch(t *testing.T) {
	tests := []struct {
		name    string
		matcher ComponentMatcher
		test    v1.TestInfo
		matches bool
	}{
		{
			name: "sig matches",
			matcher: ComponentMatcher{
				SIG: "sig-network-edge",
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: true,
		},
		{
			name: "sig does not match",
			matcher: ComponentMatcher{
				SIG: "sig-auth",
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: false,
		},
		{
			name: "sig and suite and include any matches",
			matcher: ComponentMatcher{
				SIG:        "sig-network-edge",
				Suite:      "openshift-tests",
				IncludeAny: []string{"Unidling", "Ingress"},
			},
			test: v1.TestInfo{
				Name:  "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
				Suite: "openshift-tests",
			},
			matches: true,
		},
		{
			name: "include all matches all",
			matcher: ComponentMatcher{
				IncludeAll: []string{"Unidling", "apigroup:route.openshift.io"},
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: true,
		},
		{
			name: "include all does not match when all are not present",
			matcher: ComponentMatcher{
				IncludeAll: []string{"Unidling", "OpenShiftSDN"},
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: false,
		},
		{
			name: "exclude all matches all",
			matcher: ComponentMatcher{
				ExcludeAll: []string{"Unidling", "Skipped:Network/OVNKubernetes"},
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: false,
		},
		{
			name: "exclude all does not match when all are not present",
			matcher: ComponentMatcher{
				ExcludeAll: []string{"Unidling", "UDP"},
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: true,
		},
		{
			name: "exclude any matches any value",
			matcher: ComponentMatcher{
				ExcludeAny: []string{"Unidling", "Ingress"},
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: false,
		},
		{
			name: "exclude any matches no value",
			matcher: ComponentMatcher{
				ExcludeAny: []string{"Ingress", "Routing"},
			},
			test: v1.TestInfo{
				Name: "[sig-network-edge][Feature:Idling] Unidling [apigroup:apps.openshift.io][apigroup:route.openshift.io] should handle many TCP connections by possibly dropping those over a certain bound [Serial] [Skipped:Network/OVNKubernetes]",
			},
			matches: true,
		},
		{
			name: "Suite matches regex",
			matcher: ComponentMatcher{
				Suite: "^CNV-lp-interop.*",
			},
			test: v1.TestInfo{
				Name: "[sig-compute] CNV-lp-interop:VirtualMachine A valid VirtualMachine given Using RunStrategyAlways [test_id:3165]",
			},
			matches: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{
				Matchers: []ComponentMatcher{tt.matcher},
			}
			got := c.FindMatch(&tt.test)
			if !tt.matches && got != nil {
				t.Errorf("FindMatch() matched, but we didn't want one")
			}

			if tt.matches && (got == nil || !reflect.DeepEqual(*got, tt.matcher)) {
				t.Errorf("FindMatch() did not match, but we wanted a match")
			}
		})
	}
}

func TestIHateRegexes(t *testing.T) {
	actual := ExtractNamespaceFromTestName("[sig-arch][bz-Unknown][Late] Alerts [apigroup:monitoring.coreos.com] alert/KubePodNotReady should not be at or above info in ns/openshift [Suite:openshift/conformance/parallel]")
	if actual != "openshift" {
		t.Fatal(actual)
	}
	actual = ExtractNamespaceFromTestName("[sig-arch][bz-openshift-apiserver][Late] Alerts [apigroup:monitoring.coreos.com] alert/KubePodNotReady should not be at or above info in ns/openshift-apiserver [Suite:openshift/conformance/parallel]")
	if actual != "openshift-apiserver" {
		t.Fatal(actual)
	}
}
