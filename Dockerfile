FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.24-openshift-4.21 AS builder
WORKDIR /go/src/openshift-eng/ci-test-mapping
ENV PATH="/go/bin:${PATH}"
ENV GOPATH="/go"
RUN go install k8s.io/test-infra/robots/pr-creator@latest
# install gh-token before it required go 1.23 which ubi9 doesn't have yet;
# unfortunately `go install` with the tag is broken, so just clone and install.
RUN git clone https://github.com/Link-/gh-token --branch v2.0.2 \
 && cd gh-token && go install .
COPY . .
RUN make build

FROM registry.access.redhat.com/ubi9/ubi:latest AS base
RUN dnf install -y git jq
COPY --from=builder /go/src/openshift-eng/ci-test-mapping/ci-test-mapping /bin/ci-test-mapping
COPY --from=builder /go/bin/gh-token /bin/gh-token
COPY --from=builder /go/bin/pr-creator /bin/pr-creator
COPY hack /hack
ENTRYPOINT ["/bin/ci-test-mapping"]
