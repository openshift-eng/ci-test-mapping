package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"time"

	"github.com/pkg/errors"

	"cloud.google.com/go/civil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/openshift-eng/ci-test-mapping/cmd/ci-test-mapping/flags"
	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/bigquery"
	"github.com/openshift-eng/ci-test-mapping/pkg/components"
	"github.com/openshift-eng/ci-test-mapping/pkg/jira"
	"github.com/openshift-eng/ci-test-mapping/pkg/obsoletetests"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
)

const ModeBigQuery = "bigquery"
const ModeLocal = "local"

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Map tests and job variants to component ownership",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := verifyParams(); err != nil {
			_ = cmd.Usage()
			return err
		}

		var tests []v1.TestInfo
		var testTableManager *bigquery.MappingTableManager[v1.TestOwnership]
		var variantTableManager *bigquery.MappingTableManager[v1.VariantMapping]

		testsFile := path.Join("data", f.bigqueryFlags.Project, f.bigqueryFlags.Dataset, fmt.Sprintf("%s.json", f.junitTable))
		testMappingFile := path.Join("data", f.bigqueryFlags.Project, f.bigqueryFlags.Dataset, fmt.Sprintf("%s.json", f.testMappingTable))
		variantMappingFile := path.Join("data", f.bigqueryFlags.Project, f.bigqueryFlags.Dataset, fmt.Sprintf("%s.json", f.variantMappingTable))

		config, err := f.configFlags.GetConfig()
		if err != nil {
			return err
		}

		if f.mode == ModeBigQuery {
			// Get a bigquery client
			bigqueryClient, err := bigquery.NewClient(context.Background(),
				f.bigqueryFlags.ServiceAccountCredentialFile,
				f.bigqueryFlags.OAuthClientCredentialFile, f.bigqueryFlags.Project, f.bigqueryFlags.Dataset)
			if err != nil {
				return errors.WithMessage(err, "could not obtain bigquery client")
			}

			// Create or update schema for test mapping table
			testTableManager = bigquery.NewMappingTableManager[v1.TestOwnership](context.Background(), bigqueryClient, f.testMappingTable, v1.TestMappingTableSchema)
			if err := testTableManager.Migrate(); err != nil {
				return errors.WithMessage(err, "could not migrate test mapping table")
			}
			// Create or update schema for variant mapping table
			variantTableManager = bigquery.NewMappingTableManager[v1.VariantMapping](context.Background(), bigqueryClient, f.variantMappingTable, v1.VariantMappingTableSchema)
			if err := variantTableManager.Migrate(); err != nil {
				return errors.WithMessage(err, "could not migrate variant mapping table")
			}

			// Get a list of all tests from bigquery - this could be swapped out with other
			// mechanisms to get test details later on.
			testLister := bigquery.NewTestTableManager(context.Background(), bigqueryClient, config, f.junitTable, f.jobVariantsTable)
			tests, err = testLister.ListTests()
			if err != nil {
				return errors.WithMessage(err, "could not list tests")
			}
			if err := writeRecords(tests, testsFile); err != nil {
				return errors.WithMessage(err, "couldn't write records")
			}
		} else {
			data, err := os.ReadFile(testsFile)
			if err != nil {
				return errors.WithMessage(err, "could not fetch tests from file")
			}
			if err := json.Unmarshal(data, &tests); err != nil {
				return errors.WithMessage(err, "could not marshal tests from file")
			}
		}

		// Create a registry of components
		componentRegistry := registry.NewComponentRegistry()

		// Query each component for each test
		now := time.Now()
		createdAt := civil.DateTimeOf(now)
		log.Infof("mapping tests to ownership")

		jiraComponentIDs, err := jira.GetJiraComponents()
		if err != nil {
			return errors.WithMessage(err, "could not get jira component mapping")
		}
		testObsoleter := &obsoletetests.OCPObsoleteTestManager{}
		testIdentifier := components.NewTestIdentifier(componentRegistry, jiraComponentIDs)
		var newTestMappings []v1.TestOwnership
		var matched, unmatched int
		success := true
		for i := range tests {
			ownership, err := testIdentifier.Identify(tests[i])
			if err != nil {
				log.WithError(err).Warningf("encountered error in component identification")
				success = false
				continue
			}
			if ownership != nil {
				if ownership.Component == components.DefaultComponent {
					unmatched++
				} else {
					matched++
				}
				ownership.CreatedAt = createdAt

				ownership.StaffApprovedObsolete = testObsoleter.IsObsolete(&tests[i])
				newTestMappings = append(newTestMappings, *ownership)
			}
		}
		if !success {
			return fmt.Errorf("encountered errors while trying to identify tests")
		}

		// Ensure slice is sorted
		sort.Slice(newTestMappings, func(i, j int) bool {
			return newTestMappings[i].Name < newTestMappings[j].Name && newTestMappings[i].Suite < newTestMappings[j].Suite
		})

		log.WithFields(log.Fields{
			"matched":   matched,
			"unmatched": unmatched,
		}).Infof("mapping tests to ownership complete in %v", time.Since(now))

		variantMappings := []v1.VariantMapping{}
		if f.mapVariants {
			now = time.Now()
			log.Infof("mapping variants to ownership")
			variantIdentifier := components.NewVariantIdentifier(componentRegistry, jiraComponentIDs)
			variantMappings, err = variantIdentifier.Identify()
			if err != nil {
				log.WithError(err).Warningf("encountered error in variant ownership identification")
			}
			log.Infof("mapping variants to ownership complete in %v", time.Since(now))
		}

		if f.mode == ModeBigQuery && f.pushToBQ {
			now = time.Now()
			log.Infof("pushing test mappings to bigquery...")
			if err := testTableManager.PushMappings(newTestMappings); err != nil {
				return errors.WithMessage(err, "could not push test mappings to bigquery")
			}
			log.Infof("done pushing test mappings to bigquery...")
			log.Infof("pushing variant mappings to bigquery...")
			if err := variantTableManager.PushMappings(variantMappings); err != nil {
				return errors.WithMessage(err, "could not push variant mappings to bigquery")
			}
			log.Infof("done pushing variant mappings to bigquery...")
			log.Infof("push finished in %+v", time.Since(now))
		}

		if err := writeRecords(newTestMappings, testMappingFile); err != nil {
			return errors.WithMessage(err, "could not write records to test mapping file")
		}
		if f.mapVariants {
			if err := writeRecords(variantMappings, variantMappingFile); err != nil {
				return errors.WithMessage(err, "could not write records to variant mapping file")
			}
		}

		return nil
	},
}

type MapFlags struct {
	mode                string
	pushToBQ            bool
	bigqueryFlags       *flags.BigQueryFlags
	configFlags         *flags.ConfigFlags
	junitTable          string
	jobVariantsTable    string
	testMappingTable    string
	variantMappingTable string
	mapVariants         bool
}

var f = NewMapFlags()

func NewMapFlags() *MapFlags {
	return &MapFlags{
		bigqueryFlags: flags.NewBigQueryFlags(),
		configFlags:   flags.NewConfigFlags(),
	}
}

func (f *MapFlags) BindFlags(fs *pflag.FlagSet) {
	f.bigqueryFlags.BindFlags(fs)
	f.configFlags.BindFlags(fs)
}

func init() {
	mapCmd.PersistentFlags().StringVar(&f.junitTable, "table-junit", "junit", "BigQuery table name storing JUnit test results")
	mapCmd.PersistentFlags().StringVar(&f.jobVariantsTable, "table-job-variants", "job_variants", "BigQuery table name storing job variant mappings")
	mapCmd.PersistentFlags().StringVar(&f.testMappingTable, "table-mapping", "component_mapping", "BigQuery table name storing component mappings")
	mapCmd.PersistentFlags().StringVar(&f.variantMappingTable, "table-variant-mapping", "variant_mapping", "BigQuery table name storing variant mappings")
	mapCmd.PersistentFlags().StringVar(&f.mode, "mode", "local", "Mode (one of: local, bigquery). Local mode doesn't require access to BigQuery and is suitable for local development.")
	mapCmd.PersistentFlags().BoolVar(&f.pushToBQ, "push-to-bigquery", false, "whether or not to push the updated records to bigquery")
	mapCmd.PersistentFlags().BoolVar(&f.mapVariants, "map-variant", false, "whether or not to map variants to jira projects and components")
	f.BindFlags(mapCmd.Flags())
	rootCmd.AddCommand(mapCmd)
}

func verifyParams() error {
	switch f.mode {
	case ModeBigQuery:
		if f.bigqueryFlags.ServiceAccountCredentialFile == "" && f.bigqueryFlags.OAuthClientCredentialFile == "" {
			return fmt.Errorf("please supply bigquery credentials, or use --mode=local") //nolint
		}
	case ModeLocal:
		if f.pushToBQ {
			return fmt.Errorf("cannot push to bigquery in --mode=local") //nolint
		}

		if f.bigqueryFlags.ServiceAccountCredentialFile != "" || f.bigqueryFlags.OAuthClientCredentialFile != "" {
			return fmt.Errorf("bigquery credentials not required for local mode, maybe you meant to specify --mode=bigquery") //nolint
		}
	default:
		return fmt.Errorf("invalid mode, must be one of: bigquery, local. got: %q", f.mode) //nolint
	}

	return nil
}

func writeRecords(records interface{}, filename string) error {
	now := time.Now()
	log.Infof("writing results to file")
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.WithError(err).Errorf("could not open file for writing")
		return err
	}
	jsonEncoder := json.NewEncoder(f)
	jsonEncoder.SetIndent("", "  ")

	err = jsonEncoder.Encode(records)
	log.Infof("write complete in %+v", time.Since(now))
	return err
}
