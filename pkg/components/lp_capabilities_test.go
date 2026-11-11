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
	{"lp-chaos--ACM--ACM+Virt--GRC", "LP--ACM", ""},
	{"lp-chaos--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", ""},
	{"lp-ocp-compat--ACM--ACM+Virt--GRC", "LP--ACM", ""},
	{"lp-ocp-compat--CNV--ACM+Virt--Testing VM workload.", "LP--CNV", ""},
	// Canonical component-first OPP suite names (agreed 2026-10-23) route
	// with the OPP capability; for ACM this is also the ACM+Virt negative
	// (canonical OPP suites carry the OPP capability, not ACM+Virt).
	// OPP matchers exist for ACM/ACS/CNV/ODF/Quay; canonical OCP-class
	// platform-infra suites (lp-interop--OCP--OPP--*) intentionally remain
	// Unknown — no OCP component is registered, and the lp-interop--OPP
	// component was removed per review (0 rows in both datasets). When an
	// OCP home is decided, add a matcher and rows for it. Suites whose
	// matchers were removed per the exact structure agreed in review
	// (broad component rules, ACM-Virt, and the runtime-first OPP test
	// patterns) are pinned as Unknown in TestKnownUnroutedLPSuites, except
	// the two runtime-first names whose producers are still live
	// (odf-health, acm-obs-odf), which are routed below because their rows
	// stay ingested.
	{"lp-interop--ACM--OPP--acm/policy/all", "LP--ACM", "OPP"},
	{"lp-interop--ACM--OPP--acm/policy/smoke", "LP--ACM", "OPP"},
	{"lp-interop--ACM--OPP--acm/lifecycle/all", "LP--ACM", "OPP"},
	{"lp-interop--ACM--OPP--acm/lifecycle/smoke", "LP--ACM", "OPP"},
	{"lp-interop--ACS--OPP--acs/all", "LP--ACS", "OPP"},
	{"lp-interop--CNV--OPP--cnv/lifecycle/all", "LP--CNV", "OPP"},
	{"lp-interop--CNV--OPP--cnv/lifecycle/smoke", "LP--CNV", "OPP"},
	{"lp-interop--ODF--OPP--odf/all", "LP--ODF", "OPP"},
	{"lp-interop--Quay--OPP--quay/smoke", "LP--Quay", "OPP"},
	// Runtime-first OPP emissions that are still live in openshift/release
	// (interop-opp-odf-health, interop-opp-observability-odf,
	// interop-opp-product-upgrade-acs). Their exact suite names stay admitted
	// by config/openshift-eng.yaml and routed here; the exclusion now covers
	// only the OPP names with no producer left.
	{"lp-interop--OPP--odf-health", "LP--ODF", "OPP"},
	{"lp-interop--OPP--acm-obs-odf", "LP--ACM", "OPP"},
	{"lp-interop--OPP--acs-upgrade", "LP--ACS", "OPP"},
	{"lp-interop--OPP--acs-deploy", "LP--ACS", "OPP"},
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
			suffix      string
			wantPattern string // empty means no matcher may claim the suite
		}{
			{"ACM+Virt--GRC", `ACM\+Virt--`},
			{"ACM+Virt--Testing VM workload.", `ACM\+Virt--`},
			{"OPP--GRC", `OPP--`},
			{"OPP--acm/policy/all", `OPP--`},
			{"ACMVirt--GRC", ""},
			{"ACMMVirt--GRC", ""},
			{"ACM+Virtual--GRC", ""},
			{"ACM+Virt", ""},
			{"ACM+Virt-GRC", ""},
			{"standalone", ""},
		}
		for _, testCase := range tests {
			t.Run(component.name+"/"+testCase.suffix, func(t *testing.T) {
				t.Parallel()
				matcher := component.component.FindMatch(&v1.TestInfo{Suite: prefix + testCase.suffix})
				if testCase.wantPattern == "" {
					// The broad component rules are gone: only the exact
					// runtime arms may claim; everything else stays Unknown.
					if matcher != nil {
						t.Errorf("FindMatch() = %+v, want no match", matcher)
					}
					return
				}
				wantPattern := "^" + prefix + testCase.wantPattern
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
		// The legacy exact-Suite matchers ("<X>-lp-interop") were removed
		// per maintainer review (2026-10-29), with one exception directed by
		// etirta's 2026-10-30 22:34Z correction: only {Suite:
		// "ACM-lp-interop"} is restored. The ACM claim is exact — suffixed
		// variants must not match — and the other flat names stay unclaimed
		// (ACS/ODF/Quay are pinned in knownUnroutedSuites).
		t.Run(component.name+"/flat suite claim", func(t *testing.T) {
			t.Parallel()
			suite := component.name + "-lp-interop"
			matcher := component.component.FindMatch(&v1.TestInfo{Suite: suite})
			if component.name == "ACM" {
				if matcher == nil || matcher.Suite != suite {
					t.Errorf("FindMatch() = %+v, want exact-Suite match for flat suite %q", matcher, suite)
				}
			} else if matcher != nil {
				t.Errorf("FindMatch() = %+v, want no match for flat suite %q", matcher, suite)
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

// usStatus is a pinned name's intended reachability in the production US
// mapping pipeline (config/openshift-eng.yaml). Suites are filtered by
// includeSuites/includeSuitePatterns minus excludeSuites/excludeSuitePatterns
// (compiled to SQL in pkg/bigquery/test_table.go) before any matcher runs, so
// a filtered name can never produce a row however the matchers resolve.
type usStatus string

const (
	// usAdmitted: the name survives the US suite filter, so its pin
	// describes routing the pipeline can still produce.
	usAdmitted usStatus = "admitted"
	// usUnreachable: the name is filtered out before matching, so its pin
	// fixes matcher behaviour for a name US ingestion never sees (retired
	// producers, superseded flat suites). The unfiltered QE run
	// (--config "") is out of scope here.
	usUnreachable usStatus = "unreachable"
)

// oppRoutingPin is one pinned routing row: the test identity to resolve, the
// ownership the identifier must yield today, and the row's intended US
// reachability.
type oppRoutingPin struct {
	testName         string
	suite            string
	wantComponent    string
	wantCapabilities []string
	wantPriority     int
	us               usStatus
}

// knownOPPMisroutes pins today's WRONG routing for runtime-first OPP infra
// suites, observed in data/openshift-gce-devel/ci_analysis_us only (0 rows in
// ci_analysis_qe): smoke/backup/wait-mcp have no suite matcher and resolve to
// Unknown. The upgrade-flavoured names were all stolen by Cluster Version
// Operator, which matches the TEST NAME substring "upgrade" (Priority -10,
// capability ClusterUpgrade); acs-upgrade, odf-upgrade and quay-upgrade each
// now keep their component through an exact Suite matcher instead (see
// oppYieldPins).
// TODO(INTEROP-9629): update these rows when the producer rename
// (openshift/release#86217) emits canonical lp-interop--OCP--OPP--* names.
//
// The us field declares each row's intended reachability, so a pin on a name
// US ingestion no longer produces reads as matcher history rather than a live
// misroute; TestOPPFixtureSuitesAdmittedByUSConfig enforces it.
var knownOPPMisroutes = []oppRoutingPin{
	{"acm-connectivity", "lp-interop--OPP--smoke", "Unknown", []string{DefaultCapability}, 0, usAdmitted},
	{"backup", "lp-interop--OPP--backup", "Unknown", []string{DefaultCapability}, 0, usAdmitted},
	{"wait-mcp", "lp-interop--OPP--wait-mcp", "Unknown", []string{DefaultCapability}, 0, usAdmitted},
	// Canonical OCP-class platform-infra suites stay Unknown too: no OCP
	// component is registered and the lp-interop--OPP component was removed
	// per review (0 rows in both datasets). TODO: give these a home when the
	// OCP component question is decided.
	{"wait-mcp", "lp-interop--OCP--OPP--wait-mcp", "Unknown", []string{DefaultCapability}, 0, usAdmitted},
	// The runtime-first OPP test patterns were removed per the exact
	// structure agreed in review (etirta's 2026-10-30 22:34Z correction
	// restored only the ACM flat Suite), so their upgrade-flavoured names
	// lose the LP component claim and are stolen by Cluster Version
	// Operator. acm-upgrade is also excluded from US ingestion (no producer
	// left in openshift/release), so only the routing is pinned here, and the
	// row is declared unreachable rather than a live misroute.
	{"acm-upgrade", "lp-interop--OPP--acm-upgrade", "Cluster Version Operator", []string{"ClusterUpgrade"}, -10, usUnreachable},
}

// oppYieldPins proves product suite matchers at Priority 0 beat Cluster
// Version Operator's Priority -10 name-based claim on upgrade-flavoured test
// names: the no-tie/no-steal property the canonical matchers rely on. Every
// row is usAdmitted, because a yield pin only means something where US
// ingestion still produces the name.
var oppYieldPins = []oppRoutingPin{
	{"acm-upgrade", "lp-interop--ACM--OPP--acm/lifecycle/all", "LP--ACM", []string{"OPP"}, 0, usAdmitted},
	{"acs-upgrade", "lp-interop--ACS--OPP--acs/all", "LP--ACS", []string{"OPP"}, 0, usAdmitted},
	// Runtime-first acs-upgrade emission: its exact Suite matcher must also
	// outrank Cluster Version Operator's -10 name-based claim.
	{"acs-upgrade", "lp-interop--OPP--acs-upgrade", "LP--ACS", []string{"OPP"}, 0, usAdmitted},
	// Same yield property for the runtime-first odf-upgrade and quay-upgrade
	// emissions, whose test names also carry the "upgrade" substring CVO
	// claims by name.
	{"odf-upgrade", "lp-interop--OPP--odf-upgrade", "LP--ODF", []string{"OPP"}, 0, usAdmitted},
	{"quay-upgrade", "lp-interop--OPP--quay-upgrade", "LP--Quay", []string{"OPP"}, 0, usAdmitted},
}

// unroutedPin is one pinned Unknown row: the suite that must resolve to
// Unknown, plus its intended US reachability.
type unroutedPin struct {
	suite string
	us    usStatus
}

// knownUnroutedSuites pins the LP suites whose matchers were removed per the
// exact structure agreed in review: the broad ^lp-interop--<COMP>-- rules,
// ^lp-interop--ACM-Virt--, the runtime-first OPP test patterns, and the
// legacy exact-Suite matchers ("<X>-lp-interop", removed per maintainer
// review 2026-10-29 — producers migrate to the canonical prefix). Etirta's
// 2026-10-30 22:34Z correction restored only {Suite: "ACM-lp-interop"};
// the other flat names stay unrouted. The runtime-first OPP patterns stay
// unrouted too, except the two live emissions (odf-health, acm-obs-odf)
// that the ingestion narrowing keeps admitted and routes to LP--ODF /
// LP--ACM. These resolve to Unknown; TODO(INTEROP-9629): producers emit
// canonical lp-interop--<COMP>--OPP--* names (openshift/release#86217) and
// these rows disappear as the raw emissions retire.
//
// The us field separates live Unknown rows (still admitted, still visible in
// US data) from retired names that only matcher hygiene keeps here, so no row
// reads as a live unrouted suite when the pipeline cannot produce it.
var knownUnroutedSuites = []unroutedPin{
	// Flat "<X>-lp-interop" names: dropped from includeSuites, and nothing
	// matches them by pattern, so US ingestion never sees them.
	{"ACS-lp-interop", usUnreachable},
	{"CNV-lp-interop", usUnreachable},
	{"ODF-lp-interop", usUnreachable},
	{"Quay-lp-interop", usUnreachable},
	// Admitted by "lp-interop--%": live rows that must stay Unknown.
	{"lp-interop--ACM-Virt--GRC", usAdmitted},
	{"lp-interop--ACM-Virt--observability", usAdmitted},
	{"lp-interop--ACM--ACMVirt--GRC", usAdmitted},
	{"lp-interop--CNV--ACMVirt--Testing VM workload.", usAdmitted},
	{"lp-interop--ACM--ACM+Virtual--GRC", usAdmitted},
	{"lp-interop--CNV--ACM+Virtual--Testing VM workload.", usAdmitted},
	{"lp-interop--ACM--ACM+Virt", usAdmitted},
	{"lp-interop--CNV--ACM+Virt", usAdmitted},
	{"lp-interop--ACM--ACM+Virt-GRC", usAdmitted},
	{"lp-interop--CNV--ACM+Virt-Testing VM workload.", usAdmitted},
	{"lp-interop--ACM--standalone", usAdmitted},
	{"lp-interop--CNV--standalone", usAdmitted},
	{"lp-interop--ACS--standalone", usAdmitted},
	{"lp-interop--ODF--standalone", usAdmitted},
	{"lp-interop--Quay--standalone", usAdmitted},
	{"lp-interop--OPP--Observability", usAdmitted},
	// Bare name with no include entry and no pattern match.
	{"acm-opp-app", usUnreachable},
	{"lp-interop--OPP--acm", usAdmitted},
	// Excluded by exact name (retired producer / superseded flat suite).
	{"lp-interop--ACM-OBS-ODF", usUnreachable},
	{"lp-interop--OPP--acs", usAdmitted},
	{"lp-interop--OPP--pytest", usUnreachable},
	{"lp-interop--ODF", usUnreachable},
	// The bare runtime-first name stays admitted; only the per-case
	// "lp-interop--OPP--QUAY OCP-*" suites are excluded.
	{"lp-interop--OPP--QUAY", usAdmitted},
}

func TestKnownUnroutedLPSuites(t *testing.T) {
	t.Parallel()
	identifier := NewTestIdentifier(registry.NewComponentRegistry(), nil)
	for _, pin := range knownUnroutedSuites {
		suite := pin.suite
		t.Run(suite, func(t *testing.T) {
			t.Parallel()
			testInfo := &v1.TestInfo{Name: "should preserve VM workload availability", Suite: suite}
			ownership, err := identifier.Identify(testInfo)
			if err != nil || ownership == nil {
				t.Fatalf("Identify() = %v, %v", ownership, err)
			}
			if ownership.Component != "Unknown" || ownership.Priority != 0 {
				t.Errorf("Identify() = %s/%d, want Unknown/0", ownership.Component, ownership.Priority)
			}
			if !reflect.DeepEqual(ownership.Capabilities, []string{DefaultCapability}) {
				t.Errorf("Identify() capabilities = %v, want [%s]", ownership.Capabilities, DefaultCapability)
			}
		})
	}
}

func TestKnownOPPMisroutes(t *testing.T) {
	t.Parallel()
	identifier := NewTestIdentifier(registry.NewComponentRegistry(), nil)
	for _, row := range append(append([]oppRoutingPin{}, knownOPPMisroutes...), oppYieldPins...) {
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
// includeSuitePatterns, minus excludeSuites and excludeSuitePatterns, all
// compiled to SQL in pkg/bigquery/test_table.go
// (buildSuitesFilter/buildExcludeSuitesFilter) before any matcher runs. A
// bare exact-Suite matcher such as {Suite: "preflight"} is unreachable there,
// and reachable only in the unfiltered QE run (--config ""), where it would
// over-claim. US-only guard; not QE evidence. It covers every pinned name
// (capability fixtures and the three pin lists), each asserted against the
// reachability its us field declares, so a pin cannot describe live routing
// for a name the filter drops (or dead routing for a name it still admits).
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
	// Fail closed on pattern shapes this guard does not model:
	// non-trailing-% LIKE patterns. The *testing.T is the caller's so a
	// subtest's failure is reported against the subtest.
	admitted := func(t *testing.T, suite string) bool {
		t.Helper()
		included := false
		for _, s := range cfg.IncludeSuites {
			if s == suite {
				included = true
				break
			}
		}
		if !included {
			for _, p := range cfg.IncludeSuitePatterns {
				prefix, ok := strings.CutSuffix(p, "%")
				if !ok {
					t.Fatalf("unexpected include pattern shape %q; extend this guard, do not assume LIKE", p)
				}
				if strings.HasPrefix(suite, prefix) {
					included = true
					break
				}
			}
		}
		if !included {
			return false
		}
		// Inclusion is necessary but not sufficient: exclusions win over
		// inclusions (buildExcludeSuitesFilter ANDs NOT IN / NOT LIKE onto
		// the OR'd include filter).
		for _, s := range cfg.ExcludeSuites {
			if s == suite {
				return false
			}
		}
		for _, p := range cfg.ExcludeSuitePatterns {
			prefix, ok := strings.CutSuffix(p, "%")
			if !ok {
				t.Fatalf("unexpected exclude pattern shape %q; extend this guard, do not assume LIKE", p)
			}
			if strings.HasPrefix(suite, prefix) {
				return false
			}
		}
		return true
	}
	if admitted(t, "preflight") {
		t.Fatal("guard is vacuous: bare suite 'preflight' must not be admitted")
	}
	// The exclusion path must actually bite: the legacy runtime-first OPP
	// prefix is include-pattern-admitted (lp-interop--%) yet excluded, so a
	// matcher on such a suite would be unreachable in the US pipeline.
	if admitted(t, "lp-interop--OPP--pytest") {
		t.Fatal("guard is vacuous: excluded legacy suite 'lp-interop--OPP--pytest' must not be admitted")
	}
	if !admitted(t, "lp-interop--ACM--OPP--acm") {
		t.Fatal("guard over-excludes: canonical component-first suite 'lp-interop--ACM--OPP--acm' must stay admitted")
	}
	for _, tc := range lpCapabilityTests {
		if !strings.Contains(tc.suite, "--OPP--") {
			continue
		}
		if !admitted(t, tc.suite) {
			t.Errorf("fixture suite %q is not admitted by config/openshift-eng.yaml; its matcher can never fire in the US pipeline", tc.suite)
		}
	}
	// Every pinned name carries its intended reachability, and the guard must
	// check all of them: the pin lists exist to fix behaviour for names that
	// US ingestion may stop (or never start) producing, so an unchecked pin
	// can read as a live routing fact after the filter moves out from under it.
	assertUSStatus := func(t *testing.T, list, suite string, want usStatus) {
		t.Helper()
		got := admitted(t, suite)
		switch want {
		case usAdmitted:
			if !got {
				t.Errorf("%s pin %q is not admitted by config/openshift-eng.yaml but the row declares %q; the pin claims live US routing the pipeline can never produce", list, suite, want)
			}
		case usUnreachable:
			if got {
				t.Errorf("%s pin %q is admitted by config/openshift-eng.yaml but the row declares %q; the pin claims dead US routing the pipeline still produces", list, suite, want)
			}
		default:
			t.Errorf("%s pin %q declares no intended US reachability (%q); use %q or %q", list, suite, want, usAdmitted, usUnreachable)
		}
	}
	for _, row := range knownOPPMisroutes {
		t.Run("knownOPPMisroutes/"+row.suite, func(t *testing.T) {
			t.Parallel()
			assertUSStatus(t, "knownOPPMisroutes", row.suite, row.us)
		})
	}
	for _, row := range oppYieldPins {
		t.Run("oppYieldPins/"+row.suite, func(t *testing.T) {
			t.Parallel()
			assertUSStatus(t, "oppYieldPins", row.suite, row.us)
		})
	}
	for _, pin := range knownUnroutedSuites {
		t.Run("knownUnroutedSuites/"+pin.suite, func(t *testing.T) {
			t.Parallel()
			assertUSStatus(t, "knownUnroutedSuites", pin.suite, pin.us)
		})
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
		// The flat "<X>-lp-interop" suites lost their legacy claimant per
		// review; name-derived capabilities still apply on the Unknown path.
		// Only ACM keeps its flat claim (etirta's 2026-10-30 22:34Z
		// correction).
		{"CNV-lp-interop", "[Feature:LPFeature] checks workload", []string{"LPFeature"}, "Unknown"},
		{"CNV-lp-interop", "[Feature:ACM+Virt] checks workload", []string{"ACM+Virt"}, "Unknown"},
		{"ACM-lp-interop", "[Feature:LPFeature] checks workload", []string{"LPFeature"}, "LP--ACM"},
		{"lp-chaos--ACM--GRC", "[Feature:Chaos] checks workload", []string{"Chaos"}, "LP--ACM"},
		{"lp-ocp-compat--CNV--Testing VM workload.", "[Feature:Compatibility] checks workload", []string{"Compatibility"}, "LP--CNV"},
		{"lp-interop--ACM--OPP--acm", "[Feature:LPFeature] checks workload", []string{"LPFeature", "OPP"}, "LP--ACM"},
		{"CNV-lp-interop", "[Feature:ResourceQuota] checks quota", []string{"ResourceQuota"}, "kube-controller-manager"},
		{"lp-interop--ACM--OPP--acm", "[Feature:ResourceQuota] checks quota", []string{"ResourceQuota"}, "kube-controller-manager"},
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
