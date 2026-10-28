package components

import (
	"reflect"
	"strings"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpacm"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpacs"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpcnv"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpinteropopp"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpodf"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpquay"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
)

func TestLPCapabilities(t *testing.T) {
	t.Parallel()
	tests := []struct {
		suite      string
		component  string
		capability string
	}{
		{"lp-interop--ACM--ACM+Virt--GRC", "LP--ACM", "ACM+Virt"},
		{"lp-interop--ACM--ACM+Virt--observability", "LP--ACM", "ACM+Virt"},
		{"lp-interop--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", "ACM+Virt"},
		{"lp-interop--ACM-Virt--GRC", "LP--ACM", "ACM+Virt"},
		{"lp-interop--ACM-Virt--observability", "LP--ACM", "ACM+Virt"},
		{"lp-interop--ACM--ACMVirt--GRC", "LP--ACM", "ACM+Virt"},
		{"lp-interop--CNV--ACMVirt--Testing VM workload.", "LP--CNV", "ACM+Virt"},
		{"lp-interop--ACM--ACM+Virtual--GRC", "LP--ACM", "ACM+Virt"},
		{"lp-interop--CNV--ACM+Virtual--Testing VM workload.", "LP--CNV", "ACM+Virt"},
		{"lp-interop--ACM--ACM+Virt", "LP--ACM", "ACM+Virt"},
		{"lp-interop--CNV--ACM+Virt", "LP--CNV", "ACM+Virt"},
		{"lp-interop--ACM--ACM+Virt-GRC", "LP--ACM", "ACM+Virt"},
		{"lp-interop--CNV--ACM+Virt-Testing VM workload.", "LP--CNV", "ACM+Virt"},
		{"lp-interop--ACM--standalone", "LP--ACM", "ACM+Virt"},
		{"lp-interop--CNV--standalone", "LP--CNV", "ACM+Virt"},
		{"lp-chaos--ACM--ACM+Virt--GRC", "LP--ACM", "ACM+Virt"},
		{"lp-chaos--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", "ACM+Virt"},
		{"lp-ocp-compat--ACM--ACM+Virt--GRC", "LP--ACM", "ACM+Virt"},
		{"lp-ocp-compat--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", "ACM+Virt"},
		{"lp-interop--OPP--Observability", "LP--ACM", "opp"},
		{"acm-opp-app", "LP--ACM", "opp"},
		{"lp-interop--OPP--acm", "LP--ACM", "opp"},
		{"lp-interop--ACM-OBS-ODF", "LP--ACM", "opp"},
		// Canonical component-first OPP suite names (ADR 0007) route with the
		// opp capability; for ACM this is also the ACM+Virt negative.
		{"lp-interop--ACM--OPP--acm/policy/all", "LP--ACM", "opp"},
		{"lp-interop--ACM--OPP--acm/policy/smoke", "LP--ACM", "opp"},
		{"lp-interop--ACM--OPP--acm/lifecycle/all", "LP--ACM", "opp"},
		{"lp-interop--ACM--OPP--acm/lifecycle/smoke", "LP--ACM", "opp"},
		{"lp-interop--ACS--OPP--acs/all", "LP--ACS", "opp"},
		{"lp-interop--ODF--OPP--odf/all", "LP--ODF", "opp"},
		{"lp-interop--Quay--OPP--quay/smoke", "LP--Quay", "opp"},
		{"lp-interop--OCP--OPP--wait-mcp", "lp-interop--OPP", "opp"},
		// TODO(INTEROP-9629): transitional bridges for currently-unprefixed
		// emissions; remove with the producers' canonical rename.
		{"opp-restore", "lp-interop--OPP", "opp"},
		{"preflight", "lp-interop--OPP", "opp"},
		{"ocp-upgrade", "lp-interop--OPP", "opp"},
		{"opp-post-upgrade-policy", "lp-interop--OPP", "opp"},
		{"opp-pre-upgrade-checks", "lp-interop--OPP", "opp"},
		{"interop-opp-product-upgrade-acm", "LP--ACM", "opp"},
		{"interop-tests-opp-quay-smoke", "LP--Quay", "opp"},
		{"ACM-lp-interop", "LP--ACM", ""},
		{"CNV-lp-interop", "LP--CNV", "ACM+Virt"},
	}
	componentRegistry := registry.NewComponentRegistry()
	identifier := NewTestIdentifier(componentRegistry, nil)
	for _, testCase := range tests {
		t.Run(testCase.suite, func(t *testing.T) {
			t.Parallel()
			testInfo := &v1.TestInfo{Name: "should preserve VM workload availability", Suite: testCase.suite}
			component := componentRegistry.Components[testCase.component]
			ownership, err := component.IdentifyTest(testInfo)
			if err != nil || ownership == nil {
				t.Fatalf("IdentifyTest() = %v, %v", ownership, err)
			}
			var wantCapabilities []string
			if testCase.capability != "" {
				wantCapabilities = []string{testCase.capability}
			}
			if !reflect.DeepEqual(ownership.Capabilities, wantCapabilities) {
				t.Errorf("IdentifyTest() capabilities = %v, want %v", ownership.Capabilities, wantCapabilities)
			}
			if component.StableID(testInfo) != testInfo.Name {
				t.Errorf("StableID() changed test name %q", testInfo.Name)
			}
			ownership, err = identifier.Identify(testInfo)
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if testCase.capability == "" {
				wantCapabilities = []string{DefaultCapability}
			}
			if !reflect.DeepEqual(ownership.Capabilities, wantCapabilities) {
				t.Errorf("Identify() capabilities = %v, want %v", ownership.Capabilities, wantCapabilities)
			}
			if ownership.Component != testCase.component || ownership.JIRAComponent != testCase.component || ownership.Priority != 0 {
				t.Errorf("Identify() ownership = %+v, want component/Jira %s and priority 0", ownership, testCase.component)
			}
			if ownership.Name != testInfo.Name || ownership.ID != util.StableID(testInfo, testInfo.Name) {
				t.Errorf("Identify() changed test identity: %+v", ownership)
			}
		})
	}
	for _, component := range []string{"ACM", "CNV"} {
		t.Run(component+" feature capability is unique", func(t *testing.T) {
			t.Parallel()
			testInfo := &v1.TestInfo{
				Name:  "[Feature:ACM+Virt] should preserve VM workload availability",
				Suite: "lp-interop--" + component + "--ACM+Virt--Testing VM workload.",
			}
			ownership, err := identifier.Identify(testInfo)
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if !reflect.DeepEqual(ownership.Capabilities, []string{"ACM+Virt"}) {
				t.Errorf("Identify() capabilities = %v, want exactly [ACM+Virt]", ownership.Capabilities)
			}
		})
	}
}

func TestLPMatcherPrecedence(t *testing.T) {
	t.Parallel()
	components := []struct {
		name      string
		component *config.Component
	}{
		{"ACM", lpacm.LPacmComponent.Component},
		{"CNV", lpcnv.LPcnvComponent.Component},
	}
	for _, component := range components {
		prefix := "lp-interop--" + component.name + "--"
		tests := []struct {
			suffix       string
			wantSpecific bool
		}{
			{"ACM+Virt--GRC", true},
			{"ACM+Virt--Testing VM workload.", true},
			{"ACMVirt--GRC", false},
			{"ACMMVirt--GRC", false},
			{"ACM+Virtual--GRC", false},
			{"ACM+Virt", false},
			{"ACM+Virt-GRC", false},
			{"standalone", false},
		}
		for _, testCase := range tests {
			t.Run(component.name+"/"+testCase.suffix, func(t *testing.T) {
				t.Parallel()
				wantPattern := "^" + prefix
				if testCase.wantSpecific {
					wantPattern += `ACM\+Virt--`
				}
				matcher := component.component.FindMatch(&v1.TestInfo{Suite: prefix + testCase.suffix})
				if matcher == nil || matcher.SuiteRegEx == nil {
					t.Fatalf("FindMatch() = %+v, want regex %q", matcher, wantPattern)
				}
				if matcher.SuiteRegEx.String() != wantPattern {
					t.Errorf("FindMatch() regex = %q, want %q", matcher.SuiteRegEx, wantPattern)
				}
			})
		}
		t.Run(component.name+"/anchored prefix", func(t *testing.T) {
			t.Parallel()
			if matcher := component.component.FindMatch(&v1.TestInfo{Suite: "unrelated--" + prefix + "ACM+Virt--GRC"}); matcher != nil {
				t.Errorf("FindMatch() = %+v, want no match", matcher)
			}
		})
		t.Run(component.name+"/exact flat suite", func(t *testing.T) {
			t.Parallel()
			suite := component.name + "-lp-interop"
			matcher := component.component.FindMatch(&v1.TestInfo{Suite: suite})
			if matcher == nil || matcher.Suite != suite || matcher.SuiteRegEx != nil {
				t.Errorf("FindMatch() = %+v, want exact suite %q", matcher, suite)
			}
			if matcher := component.component.FindMatch(&v1.TestInfo{Suite: suite + "--extra"}); matcher != nil {
				t.Errorf("FindMatch() = %+v, want no match for suffixed flat suite", matcher)
			}
		})
	}
}

func TestCanonicalOPPMatcherPrecedence(t *testing.T) {
	t.Parallel()
	components := []struct {
		name      string
		component *config.Component
	}{
		{"ACM", lpacm.LPacmComponent.Component},
		{"ACS", lpacs.LPacsComponent.Component},
		{"ODF", lpodf.LPodfComponent.Component},
		{"Quay", lpquay.LPquayComponent.Component},
	}
	for _, component := range components {
		t.Run(component.name+"/canonical OPP matcher wins", func(t *testing.T) {
			t.Parallel()
			matcher := component.component.FindMatch(&v1.TestInfo{Suite: "lp-interop--" + component.name + "--OPP--suite/all"})
			if matcher == nil || matcher.SuiteRegEx == nil {
				t.Fatalf("FindMatch() = %+v, want the OPP matcher", matcher)
			}
			if want := "^lp-interop--" + component.name + "--OPP--"; matcher.SuiteRegEx.String() != want {
				t.Errorf("FindMatch() regex = %q, want %q", matcher.SuiteRegEx, want)
			}
			if !reflect.DeepEqual(matcher.Capabilities, []string{"opp"}) {
				t.Errorf("FindMatch() capabilities = %v, want [opp]", matcher.Capabilities)
			}
		})
		t.Run(component.name+"/broad rule unchanged for non-OPP names", func(t *testing.T) {
			t.Parallel()
			matcher := component.component.FindMatch(&v1.TestInfo{Suite: "lp-interop--" + component.name + "--standalone"})
			if matcher == nil || matcher.SuiteRegEx == nil {
				t.Fatalf("FindMatch() = %+v, want the broad rule", matcher)
			}
			if want := "^lp-interop--" + component.name + "--"; matcher.SuiteRegEx.String() != want {
				t.Errorf("FindMatch() regex = %q, want %q (broad rule bytes must stay identical)", matcher.SuiteRegEx, want)
			}
			// The broad rule's capabilities must stay byte-identical too: ACM
			// carries ACM+Virt, the others carry none.
			var wantCapabilities []string
			if component.name == "ACM" {
				wantCapabilities = []string{"ACM+Virt"}
			}
			if !reflect.DeepEqual(matcher.Capabilities, wantCapabilities) {
				t.Errorf("FindMatch() broad-rule capabilities = %v, want %v", matcher.Capabilities, wantCapabilities)
			}
		})
	}
	t.Run("OCP platform-infra class", func(t *testing.T) {
		t.Parallel()
		matcher := lpinteropopp.LPinteropOPPComponent.FindMatch(&v1.TestInfo{Suite: "lp-interop--OCP--OPP--wait-mcp"})
		if matcher == nil || matcher.SuiteRegEx == nil || matcher.SuiteRegEx.String() != `^lp-interop--OCP--OPP--` {
			t.Fatalf("FindMatch() = %+v, want the OCP OPP matcher", matcher)
		}
		if !reflect.DeepEqual(matcher.Capabilities, []string{"opp"}) {
			t.Errorf("FindMatch() capabilities = %v, want [opp]", matcher.Capabilities)
		}
	})
}

func TestEqualPriorityMultiClaimErrors(t *testing.T) {
	t.Parallel()
	tied := &v1.TestOwnership{Name: "health", Suite: "lp-interop--X--OPP--all", Component: "LP--A", Priority: 0}
	other := &v1.TestOwnership{Name: "health", Suite: "lp-interop--X--OPP--all", Component: "LP--B", Priority: 0}
	_, err := getHighestPriority([]*v1.TestOwnership{tied, other})
	if err == nil {
		t.Fatal("getHighestPriority() = nil error, want an equal-priority multi-claim error")
	}
	if !strings.Contains(err.Error(), "unable to resolve conflict") {
		t.Errorf("getHighestPriority() error = %q, want the tie-conflict message", err)
	}
	high := &v1.TestOwnership{Name: "health", Suite: "lp-interop--X--OPP--all", Component: "LP--C", Priority: 1}
	owner, err := getHighestPriority([]*v1.TestOwnership{tied, high})
	if err != nil || owner.Component != "LP--C" {
		t.Fatalf("getHighestPriority() = %v, %v; want the priority-1 owner without error", owner, err)
	}
}

func TestLPAdditionalCapabilities(t *testing.T) {
	t.Parallel()
	tests := []struct {
		suite            string
		name             string
		wantCapabilities []string
		wantComponent    string
	}{
		{"CNV-lp-interop", "[Feature:LPFeature] checks workload", []string{"ACM+Virt", "LPFeature"}, "LP--CNV"},
		{"CNV-lp-interop", "[Feature:ACM+Virt] checks workload", []string{"ACM+Virt"}, "LP--CNV"},
		{"ACM-lp-interop", "[Feature:LPFeature] checks workload", []string{"LPFeature"}, "LP--ACM"},
		{"lp-chaos--ACM--GRC", "[Feature:Chaos] checks workload", []string{"ACM+Virt", "Chaos"}, "LP--ACM"},
		{"lp-ocp-compat--CNV--Testing VM workload.", "[Feature:Compatibility] checks workload", []string{"ACM+Virt", "Compatibility"}, "LP--CNV"},
		{"lp-interop--OPP--Observability", "[Feature:LPFeature] checks workload", []string{"LPFeature", "opp"}, "LP--ACM"},
		{"CNV-lp-interop", "[Feature:ResourceQuota] checks quota", []string{"ResourceQuota"}, "kube-controller-manager"},
		{"lp-interop--OPP--Observability", "[Feature:ResourceQuota] checks quota", []string{"ResourceQuota"}, "kube-controller-manager"},
	}
	identifier := NewTestIdentifier(registry.NewComponentRegistry(), nil)
	for _, testCase := range tests {
		t.Run(testCase.suite+"/"+testCase.name, func(t *testing.T) {
			t.Parallel()
			testInfo := &v1.TestInfo{Name: testCase.name, Suite: testCase.suite}
			ownership, err := identifier.Identify(testInfo)
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if !reflect.DeepEqual(ownership.Capabilities, testCase.wantCapabilities) {
				t.Errorf("Identify() capabilities = %v, want %v", ownership.Capabilities, testCase.wantCapabilities)
			}
			if ownership.Component != testCase.wantComponent {
				t.Errorf("Identify() component = %q, want %q", ownership.Component, testCase.wantComponent)
			}
			if ownership.Name != testInfo.Name || ownership.ID != util.StableID(testInfo, testInfo.Name) {
				t.Errorf("Identify() changed test identity: %+v", ownership)
			}
		})
	}
}
