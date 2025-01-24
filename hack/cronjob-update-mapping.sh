#!/bin/bash
set -x
export HOME=/tmp
# Prune older records first, we can't prune after we write
# to the table with records in the streaming buffer
ci-test-mapping prune

set -o errexit
set -o pipefail

# Clone the repo
cd /tmp
git clone https://github.com/openshift-eng/ci-test-mapping.git
cd ci-test-mapping
git checkout -b update

# Generate and push mapping
## OCP Engineering Mappings
ci-test-mapping map --mode=bigquery --push-to-bigquery --map-variant
## QE Mappings
ci-test-mapping map --mode=bigquery --push-to-bigquery \
	            --bigquery-dataset ci_analysis_qe --bigquery-project openshift-gce-devel \
	            --table-junit junit --table-mapping component_mapping --config "" --map-variant

if git diff --quiet
then
  echo "No changes."
  exit 0
fi

# get token with write ability (after querying DB; do not give the token a chance to expire)
keyfile="${GHAPP_KEYFILE:-/secrets/ghapp/private.key}"
set +x
trt_token=`gh-token generate --app-id 1046118 --key "$keyfile" --installation-id 57361690 --token-only`  # 57361690 = openshift-trt
git remote add openshift-trt "https://oauth2:${trt_token}@github.com/openshift-trt/ci-test-mapping.git"
set -x
git commit -a -m "Update test mapping"
git push openshift-trt update --force

pr-creator -github-app-private-key-path "${keyfile}" -github-app-id 1046118 \
           -org openshift-eng -repo ci-test-mapping \
	   -source openshift-trt:update -branch main \
	   -body "Automatically generated component mapping update" \
	   -title "Update component readiness test mapping" \
	   -confirm
