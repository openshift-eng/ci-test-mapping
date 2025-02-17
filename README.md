# Component Readiness Test Mapping

Component Readiness needs to know how to map each test to a particular
component and its capabilities. This tool:

1. Takes a set of metadata about all tests such as its name and suite
2. Maps the test to exactly one component that provides details about the capabilities that test is testing
3. Writes the result to a json file comitted to this repo in `data/`
4. Pushes the result to BigQuery

## Updating Mappings

# Overview

Teams own their component code under `pkg/components/<component_name>`
and can handle the mapping as they see fit. New components can copy from
`pkg/components/example` and modify it, or write their own
implementation of the interface. The sample Config interface, which you
can include in your component, provides rich filters on substrings,
sigs, operators, etc.

TRT has made a first pass at assigning ownership to teams, but it's
likely we haven't correctly assigned all tests. Please re-assign tests
to their correct components as needed.  Please also create an OWNERS
file in your directory so your team can manage components and
capabilities without TRT's intervention.

Component owners should return a `TestOwnership` struct from their
identification function. See the details in `pkg/api/v1` for details
about the `TestOwnership` struct.

They should return nil when the test is not theirs.  They should ONLY
return an error on a fatal error such as inability to read from a file.

A test must only map to one component, but may map to several
capabilities.  In the event that two components are vying for a test's
ownership, you may use the `Priority` field in the `TestOwnership`
struct.  The highest value wins.

# Process

These are the general steps for updating a test's ownership:

1. Open the `component.go` for a particular component in
   `pkg/components`. If a test is tied to a namespace, operator, or
   particular SIG, update that section. If you want to match on the
   suite or test name, add a new matcher using one of the existing
   types, such as `Suite`, or `IncludeAny`, which matches a test
   substring.
2. Update the `Priority` field if you need to force a higher or lower
   priority for this component.
3. Run `make mapping`, to apply the new rules to the test names. Review
   the results of the changes by using `git diff`.  Ensure that the
   changes you see are expected.
4. Commit the result and open a pull request on GitHub.

If you'd like to annotate a test as having additional capabilities,
update `capabilities.go`. A test may have multiple capabilities, but it
can only belong to a single component.

`make mapping` enforces this, as well as disallowing any test to move
back to `Unknown`. If a test is incorrectly attributed, you must find
the owner and move it to their component.

## Renaming tests

The unfortunate reality is tests may get renamed, so we need to have a
way to compare the test results across renames. To do that, each test
has a stable ID which is the current test name stored in the DB as an
md5sum.

The test's first stable ID is the one that remains. Component owners are
responsible for ensuring the `StableID` function in their component
returns the same ID for all names of a given test. This can be done with
a simple look-up map; see the monitoring component for an example.

## Removing tests

If a test is removed, or is refactored in such a way (i.e. one to many)
that it's not reasonable to mark them as renames, then it should be
tracked as an obsolete test. In `pkg/obsoletetests` one can manage their
obsolete tests by adding an entry to the set.  For OCP, only staff
engineers can approve a test's removal.

# Test Sources

Currently the tests used for mapping come from the corpus of tests
previously seen in job results. This list is filtered down to a
smaller quantity by selecting only those in certain suites, jobs, or
matching certain names.  This is configurable by specifying a
configuration file. An example is present in
`config/openshift-eng.yaml`.

At a mimimum though, for compatibility with component readiness (and all
other OpenShift tooling), a test must:

* always have a result when it runs, indicating success, flake, or failure (historically some tests only report failure)

* belong to a test suite

* have a stable name: do not use dynamic names such as full pod names in tests

* have a reasonable way to map to component/capabilities, such as `[sig-XYZ]` present in the test name, and using `[Feature:XYZ]` or `[Capability:XYZ]` to make mapping to capabilities easier

# Usage

See --help for more info.

## Test Mapping

To find unmapped tests, run `make unmapped`.

### Development

For production usage we fetch and push data to BigQuery, but for local
testing you can used locally comitted copies of that data by using
`--mode local`:

```
ci-test-mapping map --mode local
```

### Production

For production, use `--mode bigquery` and provide credentials:

```
ci-test-mapping map --mode bigquery \
  --google-service-account-credential-file ~/bq.json \
  --log-level debug \
  --mapping-file mapping.json \
  --push-to-bigquery
```

### Alternative data sources/destinations

The BigQuery project, dataset, JUnit table, and component mapping tables
are all configurable.

```
ci-test-mapping map \
    --mode bigquery \
    --bigquery-project openshift-gce-devel \
    --bigquery-dataset ci_analysis_us \
    --table-junit junit \
    --table-mapping component_mapping
```

### Using the BigQuery table for lookups

The BigQuery mapping table should be updated only in append mode (aside from
older entries being trimmed), so mappings should limit their results to the
most recent entry.

## Updating Jira Components

The command `./ci-test-mapping jira-verify` will inform you of any Jira
components that have been removed, and report on any new Jira components
that are available.

To create any missing components, run `./ci-test-mapping jira-create`.
If any components are being renamed, run `jira-create` first and then
move the configuration to the new component.

Run `./ci-test-mapping jira-verify` to ensure it reports back cleanly.

You'll need to set the env var `JIRA_TOKEN` to your personal API token
(which you can create from your Jira profile page). Then:

1. Move any configuration for renamed components
2. Delete any obsolete `pkg/components/<component>` directory
3. Remove references to removed components from `pkg/registry`.

# Component Readiness Variant Mapping

Variant mapping maps a particular job variant to a Jira component. This helps us map 
a column on Component Readiness dashboard to its proper owners. Variant mapping is not
enabled by default. One can use ```--map-variant=true``` (or simply ```--map-variant```) to enable it.

## Add Variant Mapping to components

To add variant mapping to one particular component, all you need to do is to initialize 
Variants slice in your component definition. For example, the following was added to
pkg/components/cloudcompute/vsphereprovider/component.go to map Platform:vsphere to 
component "Cloud Compute / vSphere Provider"

```
		Variants:             []string{"Platform:vsphere"},
```
