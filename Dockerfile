FROM registry.access.redhat.com/ubi9/ubi:latest AS builder
WORKDIR /go/src/openshift-eng/ci-test-mapping
COPY . .
ENV PATH="/go/bin:${PATH}"
ENV GOPATH="/go"
RUN dnf install -y \
        git \
        go \
        make && make build

FROM gcr.io/k8s-prow/pr-creator:latest AS prcreator

FROM registry.access.redhat.com/ubi9/ubi:latest AS base
COPY --from=builder /go/src/openshift-eng/ci-test-mapping/ci-test-mapping /bin/ci-test-mapping
COPY --from=prcreator /ko-app/pr-creator /bin/pr-creator
RUN dnf install -y git
ENTRYPOINT ["/bin/ci-test-mapping"]
