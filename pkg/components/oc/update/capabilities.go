package ocupdate

import (
	"regexp"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/util"
	"k8s.io/apimachinery/pkg/util/sets"
)

// OCUpdate is about oc adm sub-command which belongs to OTA, such as oc adm upgrade or oc adm release extract
const OCUpdate = "OCUpdate"

var ocUpdateCapabilitiesIdentifiers = map[*regexp.Regexp]string{
	// The junit report of Prow upgrade ci is like pattern "upgrade should succeed: $UPGRADE_FAILURE_TYPE"
	regexp.MustCompile(`.*upgrade should succeed: (oc_update).*`): OCUpdate,
	regexp.MustCompile(`.*OTA oc should.*`):                       OCUpdate,
}

func identifyCapabilities(test *v1.TestInfo) []string {
	capabilities := sets.New[string](util.DefaultCapabilities(test)...)
	for matcher, capability := range ocUpdateCapabilitiesIdentifiers {
		if matcher.MatchString(test.Name) {
			capabilities.Insert(capability)
		}
	}
	return capabilities.UnsortedList()
}
