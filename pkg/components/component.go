package components

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
)

const (
	DefaultComponent  = "Unknown"
	DefaultCapability = "Other"
	DefaultProduct    = "OpenShift"
)

type TestIdentifier struct {
	reg          *registry.Registry
	componentIDs map[string]uint
}

func New(reg *registry.Registry, componentIDs map[string]uint) *TestIdentifier {
	if componentIDs == nil {
		componentIDs = make(map[string]uint)
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

	testOwnership.Kind = v1.Kind
	testOwnership.APIVersion = v1.APIVersion

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
		testOwnership.JIRAComponentID = &id
	}

	if len(testOwnership.Capabilities) == 0 {
		testOwnership.Capabilities = []string{DefaultCapability}
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
	for _, ownership := range ownerships {
		if highest != nil && ownership.Priority == highest.Priority {
			return nil, fmt.Errorf("test %q is claimed by %s, %s - unable to resolve conflict "+
				"-- please use priority field", highest.Name, highest.Component, ownership.Component)
		}

		if highest == nil || ownership.Priority > highest.Priority {
			highest = ownership
		}
	}

	return highest, nil
}
