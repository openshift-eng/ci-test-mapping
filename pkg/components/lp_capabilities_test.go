package components

import (
	"os"
	"reflect"
	"strings"
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpacm"
	"github.com/openshift-eng/ci-test-mapping/pkg/components/lpcnv"
	"github.com/openshift-eng/ci-test-mapping/pkg/config"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
	"gopkg.in/yaml.v2"
)

// lpCapabilityTests is the LP capability fixture table. Every row resolves
// through the full identifier.Identify path, so an equal-priority
// multi-claim (tie) fails the row.
var lpCapabilityTests = []struct {
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
	{"lp-chaos--ACM--ACM+Virt--GRC", "LP--ACM", ""},
	{"lp-chaos--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", ""},
	{"lp-ocp-compat--ACM--ACM+Virt--GRC", "LP--ACM", ""},
	{"lp-ocp-compat--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", ""},
	{"lp-interop--OPP--Observability", "LP--ACM", "opp"},
	{"acm-opp-app", "LP--ACM", "opp"},
	{"lp-interop--OPP--acm", "LP--ACM", "opp"},
	{"lp-interop--ACM-OBS-ODF", "LP--ACM", "opp"},
	// Canonical component-first OPP suite names (agreed 2026-10-23) route
	// with the opp capability; for ACM this is also the ACM+Virt negative
	// (canonical OPP suites carry the opp capability, not ACM+Virt).
	// OPP matchers exist for ACM/ACS/CNV/ODF/Quay; canonical OCP-class
	// platform-infra suites (lp-interop--OCP--OPP--*) intentionally remain
	// Unknown — no OCP component is registered, and the lp-interop--OPP
	// component was removed per review (0 rows in both datasets). When an
	// OCP home is decided, add a matcher and rows for it.
	{"lp-interop--ACM--OPP--acm/policy/all", "LP--ACM", "opp"},
	{"lp-interop--ACM--OPP--acm/policy/smoke", "LP--ACM", "opp"},
	{"lp-interop--ACM--OPP--acm/lifecycle/all", "LP--ACM", "opp"},
	{"lp-interop--ACM--OPP--acm/lifecycle/smoke", "LP--ACM", "opp"},
	{"lp-interop--ACS--OPP--acs/all", "LP--ACS", "opp"},
	{"lp-interop--CNV--OPP--cnv/lifecycle/all", "LP--CNV", "opp"},
	{"lp-interop--CNV--OPP--cnv/lifecycle/smoke", "LP--CNV", "opp"},
	{"lp-interop--ODF--OPP--odf/all", "LP--ODF", "opp"},
	{"lp-interop--Quay--OPP--quay/smoke", "LP--Quay", "opp"},
	// Broad-rule outcome pins: the broad ^lp-interop--<COMP>-- rules are
	// unchanged by the OPP additions above (no capability for ACS/ODF/Quay).
	{"lp-interop--ACS--standalone", "LP--ACS", ""},
	{"lp-interop--ODF--standalone", "LP--ODF", ""},
	{"lp-interop--Quay--standalone", "LP--Quay", ""},
	{"ACM-lp-interop", "LP--ACM", ""},
	{"CNV-lp-interop", "LP--CNV", "ACM+Virt"},
}

func TestLPCapabilities(t *testing.T) {
	t.Parallel()
	tests := lpCapabilityTests
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

func TestPriorityResolution(t *testing.T) {
	t.Parallel()
	// The highest priority value wins (pkg/api/types/v1/component.go: "The
	// highest value wins"); Cluster Version Operator deliberately loses
	// upgrade tests with Priority -10. Ties at the highest value hard-error.
	zero := &v1.TestOwnership{Name: "health", Suite: "s", Component: "LP--A", Priority: 0}
	neg := &v1.TestOwnership{Name: "health", Suite: "s", Component: "LP--B", Priority: -10}
	one := &v1.TestOwnership{Name: "health", Suite: "s", Component: "LP--C", Priority: 1}
	for _, tc := range []struct {
		name   string
		owners []*v1.TestOwnership
		want   string
	}{
		{"priority 0 beats -10", []*v1.TestOwnership{zero, neg}, "LP--A"},
		{"priority 0 beats -10, reversed order", []*v1.TestOwnership{neg, zero}, "LP--A"},
		{"priority 1 beats 0", []*v1.TestOwnership{zero, one}, "LP--C"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			owner, err := getHighestPriority(tc.owners)
			if err != nil || owner == nil || owner.Component != tc.want {
				t.Fatalf("getHighestPriority() = %v, %v; want component %s without error", owner, err, tc.want)
			}
		})
	}
	if owner, err := getHighestPriority(nil); err != nil || owner != nil {
		t.Fatalf("getHighestPriority(nil) = %v, %v; want nil, nil", owner, err)
	}
}

// knownOPPMisroutes pins today's WRONG routing for runtime-first OPP infra
// suites, observed in data/openshift-gce-devel/ci_analysis_us only (0 rows in
// ci_analysis_qe): smoke/backup/wait-mcp have no suite matcher and resolve to
// Unknown; odf-upgrade/quay-upgrade are stolen by Cluster Version Operator,
// which matches the TEST NAME substring "upgrade" (Priority -10, capability
// ClusterUpgrade). TODO(INTEROP-9629): update these rows when the producer
// rename (openshift/release#86217) emits canonical lp-interop--OCP--OPP--*
// names.
var knownOPPMisroutes = []struct {
	testName         string
	suite            string
	wantComponent    string
	wantCapabilities []string
	wantPriority     int
}{
	{"acm-connectivity", "lp-interop--OPP--smoke", "Unknown", []string{DefaultCapability}, 0},
	{"backup", "lp-interop--OPP--backup", "Unknown", []string{DefaultCapability}, 0},
	{"wait-mcp", "lp-interop--OPP--wait-mcp", "Unknown", []string{DefaultCapability}, 0},
	// Canonical OCP-class platform-infra suites stay Unknown too: no OCP
	// component is registered and the lp-interop--OPP component was removed
	// per review (0 rows in both datasets). TODO: give these a home when the
	// OCP component question is decided.
	{"wait-mcp", "lp-interop--OCP--OPP--wait-mcp", "Unknown", []string{DefaultCapability}, 0},
	{"odf-upgrade", "lp-interop--OPP--odf-upgrade", "Cluster Version Operator", []string{"ClusterUpgrade"}, -10},
	{"quay-upgrade", "lp-interop--OPP--quay-upgrade", "Cluster Version Operator", []string{"ClusterUpgrade"}, -10},
}

// oppYieldPins proves product suite matchers at Priority 0 beat Cluster
// Version Operator's Priority -10 name-based claim on upgrade-flavoured test
// names — the no-tie/no-steal property the canonical matchers rely on.
var oppYieldPins = []struct {
	testName         string
	suite            string
	wantComponent    string
	wantCapabilities []string
	wantPriority     int
}{
	{"acm-upgrade", "lp-interop--OPP--acm-upgrade", "LP--ACM", []string{"opp"}, 0},
	{"acs-upgrade", "lp-interop--OPP--acs-upgrade", "LP--ACS", []string{"opp"}, 0},
}

func TestKnownOPPMisroutes(t *testing.T) {
	t.Parallel()
	identifier := NewTestIdentifier(registry.NewComponentRegistry(), nil)
	for _, row := range append(append([]struct {
		testName         string
		suite            string
		wantComponent    string
		wantCapabilities []string
		wantPriority     int
	}{}, knownOPPMisroutes...), oppYieldPins...) {
		t.Run(row.testName+"/"+row.suite, func(t *testing.T) {
			t.Parallel()
			ownership, err := identifier.Identify(&v1.TestInfo{Name: row.testName, Suite: row.suite})
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if ownership.Component != row.wantComponent || ownership.Priority != row.wantPriority {
				t.Errorf("Identify() = %s/%d, want %s/%d", ownership.Component, ownership.Priority, row.wantComponent, row.wantPriority)
			}
			if !reflect.DeepEqual(ownership.Capabilities, row.wantCapabilities) {
				t.Errorf("Identify() capabilities = %v, want %v", ownership.Capabilities, row.wantCapabilities)
			}
		})
	}
}

// TestOPPFixtureSuitesAdmittedByUSConfig guards against fixture rows (and the
// matchers they witness) that the production US mapping pipeline can never
// reach: config/openshift-eng.yaml admits suites only via includeSuites and
// includeSuitePatterns, compiled to SQL in pkg/bigquery/test_table.go
// (buildSuitesFilter/buildExcludeSuitesFilter) before any matcher runs. A
// bare exact-Suite matcher such as {Suite: "preflight"} is unreachable there,
// and reachable only in the unfiltered QE run (--config ""), where it would
// over-claim. US-only guard; not QE evidence.
func TestOPPFixtureSuitesAdmittedByUSConfig(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../config/openshift-eng.yaml")
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var cfg v1.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse config: %v", err)
	}
	// Fail closed on shapes this guard does not model: exclusion lists and
	// non-trailing-% LIKE patterns.
	if len(cfg.ExcludeSuites) > 0 || len(cfg.ExcludeSuitePatterns) > 0 {
		t.Fatalf("config now has exclusion lists; extend this guard to model them")
	}
	admitted := func(suite string) bool {
		for _, s := range cfg.IncludeSuites {
			if s == suite {
				return true
			}
		}
		for _, p := range cfg.IncludeSuitePatterns {
			prefix, ok := strings.CutSuffix(p, "%")
			if !ok {
				t.Fatalf("unexpected include pattern shape %q; extend this guard, do not assume LIKE", p)
			}
			if strings.HasPrefix(suite, prefix) {
				return true
			}
		}
		return false
	}
	if admitted("preflight") {
		t.Fatal("guard is vacuous: bare suite 'preflight' must not be admitted")
	}
	for _, tc := range lpCapabilityTests {
		if !strings.Contains(tc.suite, "--OPP--") {
			continue
		}
		if !admitted(tc.suite) {
			t.Errorf("fixture suite %q is not admitted by config/openshift-eng.yaml; its matcher can never fire in the US pipeline", tc.suite)
		}
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
		{"lp-chaos--ACM--GRC", "[Feature:Chaos] checks workload", []string{"Chaos"}, "LP--ACM"},
		{"lp-ocp-compat--CNV--Testing VM workload.", "[Feature:Compatibility] checks workload", []string{"Compatibility"}, "LP--CNV"},
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
