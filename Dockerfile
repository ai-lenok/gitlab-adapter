FROM golang AS builder

RUN mkdir /app && mkdir /app/data
ADD . /app/
WORKDIR /app
RUN make build

FROM registry.access.redhat.com/ubi9-minimal

COPY --from=builder /app/out/bin/gitlab-adapter /app/gitlab-adapter

WORKDIR /app
RUN mkdir /app/config
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/app/gitlab-adapter", "start-server"]