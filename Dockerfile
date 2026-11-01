FROM golang:1.27.3-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mcp ./cmd
COPY test/e2e ./test/e2e
RUN CGO_ENABLED=0 go test -c -tags e2e -o /out/e2e.test ./test/e2e

FROM debian:bookworm-slim AS e2e
COPY --from=build --chown=65532:65532 --chmod=755 /out/e2e.test /e2e.test
USER 65532:65532
ENTRYPOINT ["/e2e.test", "-test.v"]

FROM gcr.io/distroless/static-debian12:nonroot AS app
COPY --from=build --chown=nonroot:nonroot --chmod=755 /out/mcp /mcp
COPY --chown=nonroot:nonroot assets /app/assets
ENV MCP_ASSET_ROOT=/app/assets
EXPOSE 8080
ENTRYPOINT ["/mcp"]
CMD ["-http", "0.0.0.0:8080"]
