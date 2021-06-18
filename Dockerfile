# Builder image
FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 as builder
ENV GOPATH=$APP_ROOT
RUN env
COPY . .
RUN go build -o app github.com/aufi/must-gather-rest-wrapper/pkg


# Runner image
FROM registry.access.redhat.com/ubi8-minimal

#LABEL name="konveyor/must-gather-rest-wrapper" \
#      description="Konveyor Must Gather REST wrapper" \
#      help="For more information visit https://konveyor.io" \
#      license="Apache License 2.0" \
#      maintainer="maufart@redhat.com" \
#      summary="Konveyor Must Gather REST wrapper" \
#      url="https://quay.io/repository/konveyor/must-gather-rest-wrapper" \
#      usage="podman run konveyor/must-gather-rest-wrapper:latest" \
#      com.redhat.component="konveyor-must-gather-rest-wrapper-container" \
#      io.k8s.display-name="must-gather-rest-wrapper" \
#      io.k8s.description="Konveyor Must Gather REST wrapper" \
#      io.openshift.expose-services="" \
#      io.openshift.tags="operator,konveyor,forklift"

COPY --from=builder /opt/app-root/src/app /usr/bin/must-gather-rest-wrapper

# RUN microdnf -y install tar && microdnf clean all

ENTRYPOINT ["/usr/bin/must-gather-rest-wrapper"]
