# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22.0
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./cmd/app

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        postgresql-client \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser


RUN mkdir -p /migrations && chown -R appuser:appuser /migrations

RUN mkdir -p /seeds && chown -R appuser:appuser /seeds

USER appuser

COPY --from=build /bin/server /bin/
COPY ./migrations /migrations
COPY ./seeds /seeds

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
