FROM golang:1.26.6-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY test/e2e ./test/e2e
RUN CGO_ENABLED=0 go test -c -tags e2e -o /out/e2e.test ./test/e2e

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build --chown=nonroot:nonroot --chmod=755 /out/e2e.test /e2e.test
ENTRYPOINT ["/e2e.test", "-test.v"]
