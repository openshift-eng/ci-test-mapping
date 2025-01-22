package components

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
)

const (
	DefaultProject    = "OCPBUGS"
	DefaultComponent  = "Unknown"
	DefaultCapability = "Other"
	DefaultProduct    = "OpenShift"
)

type TestIdentifier struct {
	reg          *registry.Registry
	componentIDs map[string]int64
}

func NewTestIdentifier(reg *registry.Registry, componentIDs map[string]int64) *TestIdentifier {
	if componentIDs == nil {
		componentIDs = make(map[string]int64)
	}

	return &TestIdentifier{
		reg:          reg,
		componentIDs: componentIDs,
	}
}

func (t *TestIdentifier) Identify(test *v1.TestInfo) (*v1.TestOwnership, error) {
	var ownerships []*v1.TestOwnership

	log.WithFields(testInfoLogFields(test)).Debugf("attempting to identify test using %d components", len(t.reg.Components))
	for name, component := range t.reg.Components {
		log.WithFields(testInfoLogFields(test)).Tracef("checking component %q", name)
		ownership, err := component.IdentifyTest(test)
		if err != nil {
			log.WithError(err).Errorf("component %q returned an error", name)
			return nil, err
		}
		if ownership != nil {
			log.WithFields(testInfoLogFields(test)).Tracef("component %q claimed this test", name)
			ownerships = append(ownerships, t.setDefaults(test, ownership, component))
		}
	}

	if len(ownerships) == 0 {
		ownerships = append(ownerships, t.setDefaults(test, &v1.TestOwnership{
			ID:   util.StableID(test, test.Name),
			Name: test.Name,
		}, nil))
	}

	highestPriority, err := getHighestPriority(ownerships)
	if err != nil {
		return nil, err
	}

	uniqueCapabilities := sets.New[string](highestPriority.Capabilities...)
	highestPriority.Capabilities = uniqueCapabilities.UnsortedList()
	sort.Strings(highestPriority.Capabilities)
	return highestPriority, nil
}

func (t *TestIdentifier) setDefaults(testInfo *v1.TestInfo, testOwnership *v1.TestOwnership, c v1.Component) *v1.TestOwnership {
	if testOwnership.ID == "" && c != nil {
		testOwnership.ID = util.StableID(testInfo, c.StableID(testInfo))
	}

	testOwnership.Kind = v1.TestOwnershipKind
	testOwnership.APIVersion = v1.TestOwnershipAPIVersion

	if testOwnership.Product == "" {
		testOwnership.Product = DefaultProduct
	}

	if testOwnership.Component == "" {
		testOwnership.Component = DefaultComponent
	}

	if testOwnership.JIRAComponent == "" {
		testOwnership.JIRAComponent = DefaultComponent
	}

	if id, ok := t.componentIDs[testOwnership.JIRAComponent]; ok {
		bID := big.Rat{}
		testOwnership.JIRAComponentID = bID.SetInt64(id)
	}

	if len(testOwnership.Capabilities) == 0 {
		capabilities := util.DefaultCapabilities(testInfo)
		if len(capabilities) == 0 {
			capabilities = []string{DefaultCapability}
		}
		testOwnership.Capabilities = capabilities
	}

	if testOwnership.Suite == "" {
		testOwnership.Suite = testInfo.Suite
	}

	return testOwnership
}

func testInfoLogFields(testInfo *v1.TestInfo) log.Fields {
	return log.Fields{
		"name":  testInfo.Name,
		"suite": testInfo.Suite,
	}
}

func getHighestPriority(ownerships []*v1.TestOwnership) (*v1.TestOwnership, error) {
	var highest *v1.TestOwnership
	// First find the highest priority
	for _, ownership := range ownerships {
		if highest == nil || ownership.Priority > highest.Priority {
			highest = ownership
		}
	}
	// Now find duplicates with the same highest priority
	if highest != nil {
		for _, ownership := range ownerships {
			if ownership.Priority == highest.Priority && ownership != highest {
				return nil, fmt.Errorf("suite=%q test=%q is claimed by %s, %s - unable to resolve conflict "+
					"-- please use priority field", highest.Suite, highest.Name, highest.Component, ownership.Component)
			}
		}
	}

	return highest, nil
}

type VariantIdentifier struct {
	reg          *registry.Registry
	componentIDs map[string]int64
}

func (vi *VariantIdentifier) Identify() ([]v1.VariantMapping, error) {
	log.Debugf("attempting to map variants to jira using %d components", len(vi.reg.Components))
	now := time.Now()
	createdAt := civil.DateTimeOf(now)
	variantToMapping := map[string]v1.VariantMapping{}
	for name, component := range vi.reg.Components {
		log.Tracef("checking component %q", name)
		variants, err := component.IdentifyVariants()
		if err != nil {
			log.WithError(err).Errorf("component %q returned an error", name)
			return nil, err
		}
		for _, v := range variants {
			if vm, ok := variantToMapping[v]; ok {
				log.Errorf("component %s is trying to claim variant %s, which is already mapped to project %s component %s", name, v, vm.JiraProject, vm.JiraComponent)
				return nil, fmt.Errorf("duplicate variant mapping")
			}
			parts := strings.Split(v, ":")
			if len(parts) != 2 {
				log.Errorf("Incorrect format for variant %s", v)
				continue
			}
			jiraComponents := component.JiraComponents()
			if len(jiraComponents) > 0 {
				mapping := v1.VariantMapping{
					VariantName:   parts[0],
					VariantValue:  parts[1],
					JiraProject:   component.JiraProject(),
					JiraComponent: jiraComponents[0],
					CreatedAt:     createdAt,
				}
				variantToMapping[v] = mapping
			}
		}
	}
	mappings := []v1.VariantMapping{}
	for _, mapping := range variantToMapping {
		mappings = append(mappings, vi.setDefaults(mapping))
	}
	// sort by variants
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].VariantName < mappings[j].VariantName && mappings[i].VariantValue < mappings[j].VariantValue
	})
	return mappings, nil
}

func (vi *VariantIdentifier) setDefaults(variantMapping v1.VariantMapping) v1.VariantMapping {
	variantMapping.Kind = v1.VariantMappingKind
	variantMapping.APIVersion = v1.VariantMappingAPIVersion

	if variantMapping.Product == "" {
		variantMapping.Product = DefaultProduct
	}

	if variantMapping.JiraComponent == "" {
		variantMapping.JiraComponent = DefaultComponent
	}

	if variantMapping.JiraProject == "" {
		variantMapping.JiraProject = DefaultProject
	}

	return variantMapping
}

func NewVariantIdentifier(reg *registry.Registry, componentIDs map[string]int64) *VariantIdentifier {
	if componentIDs == nil {
		componentIDs = make(map[string]int64)
	}

	return &VariantIdentifier{
		reg:          reg,
		componentIDs: componentIDs,
	}
}
