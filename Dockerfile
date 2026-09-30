FROM golang:1.26.6-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mcp ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build --chown=nonroot:nonroot --chmod=755 /out/mcp /mcp
COPY --chown=nonroot:nonroot assets /app/assets
ENV MCP_ASSET_ROOT=/app/assets
EXPOSE 8080
ENTRYPOINT ["/mcp"]
CMD ["-http", "0.0.0.0:8080"]
