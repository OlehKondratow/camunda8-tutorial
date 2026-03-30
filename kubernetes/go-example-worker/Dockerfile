# Kontekst buildu — katalog główny repozytorium camunda8-tutorial (zob. ci/cloudbuild-workers.yaml).
FROM golang:1.21-bookworm AS build
WORKDIR /src
COPY zeebe-tutorial/ ./
RUN go mod download && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/example-worker ./cmd/example-worker

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/example-worker /example-worker
USER nonroot:nonroot
ENTRYPOINT ["/example-worker"]
