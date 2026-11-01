# Helm chart and in-cluster e2e

Implement this file. It is the task. `AGENTS.md` says `docs/plans/` is not the source of truth. That does not apply here.

Do not redesign the server. Do not add auth, TLS, Cobra, or a second listen address. If a step is blocked, stop and report. Do not invent a workaround.

## Read first

| Path | Why |
|------|-----|
| `Dockerfile` | Image already listens on `0.0.0.0:8080`, `MCP_ASSET_ROOT=/app/assets`, distroless `nonroot` (UID/GID 65532), `ENTRYPOINT ["/mcp"]`, `CMD ["-http", "0.0.0.0:8080"]`. Do not change it. The release workflow still builds this file as-is, then packages the chart beside that image. |
| `cmd/main.go` | `-http` is streamable HTTP. Routes are only `GET /healthz` (`200`, body `ok\n`) and `/mcp`. No `/mcp/` redirect. Logs on stderr. Stateful SDK handler, no `EventStore`. |
| `cmd/http_handler_test.go` | Client pattern to copy: `mcp.NewClient` + `mcp.StreamableClientTransport`. Loopback DNS-rebinding notes do not apply in-cluster: the process binds `0.0.0.0`, which is not loopback. |
| `internal/server.go` | One server, resources from the image assets, then tools. Implementation name is the string in `mcp.NewServer`. |
| `internal/tools/registry.go` | Tool set. Names: `discover_standards`, `discover_workflows`, `search_standards`, `search_workflows`. |
| `internal/tools/tools_test.go` | Overlap these calls against the live service. Do not port the in-memory mock. |
| `internal/server_test.go` | Overlap `TestBootstrapServer_Integration` and `TestBootstrapServer_Integration_UnknownResource` (read includes frontmatter; unknown URI errors). |
| `.agents/skills/test-resource-document/scripts/test_resource.py` | Overlap list + `--read standards://git/commit-messages`. That harness stays stdio. Do not retarget it at Kubernetes. |
| `Makefile` | New targets follow this style (`##` help comments, `.PHONY`). `make check` is fmt-check, typecheck, vet, lint, vuln, `go test -race`. |
| `.agents/rules/validation.md` | Local gate. E2e is extra, not a replacement. |
| `.github/workflows/ci.yml` | PR and main gate. Add one minikube e2e job here. Do not matrix it across Go versions. |
| `.github/workflows/release.yml` | Same `release: published` trigger. Also package the chart and attach the `.tgz` to that GitHub Release. Do not add a registry. |
| `.pre-commit-config.yaml` | `check-yaml` will reject Helm `{{ }}` templates unless excluded. |
| `assets/standards/git/branches.md` | Branch from updated `main`. No ticket was given; this plan authorizes `helm-http-chart`. Do not stop to ask. |
| `assets/standards/git/commit-messages.md` | Conventional Commits. Subject says what and why. |
| `assets/standards/git/commit-staging.md` | Separate source, tests, and docs. ≤500 lines changed per commit. |
| `assets/standards/git/operations.md` | No amend, no force-push, no rebase. |
| `assets/standards/git/pull-requests.md` | Title lists every commit type. Body has Added / Removed / Updated / Deleted. |
| `assets/standards/go/testing.md` | Standard `testing` package. Do not reshape production code for tests. |
| `assets/standards/go/tooling.md` | Route build and test through `make`. Do not put minikube on a pre-commit hook. |

Helm docs (do not restate them): [helm create](https://helm.sh/docs/helm/helm_create/), [chart tests](https://helm.sh/docs/topics/chart_tests/), [chart hooks](https://helm.sh/docs/topics/charts_hooks/), [helm test](https://helm.sh/docs/helm/helm_test/), [helm package](https://helm.sh/docs/helm/helm_package/). Minikube image load: [pushing images](https://minikube.sigs.k8s.io/docs/handbook/pushing/). Minikube on GitHub Actions is supported: [setup minikube in GitHub Actions](https://minikube.sigs.k8s.io/docs/tutorials/setup_minikube_in_github_actions/) (`medyagh/setup-minikube`, Linux docker driver, no nested virt). Restricted pod baseline: [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/). HPA: [Horizontal Pod Autoscaling](https://kubernetes.io/docs/concepts/workloads/autoscaling/horizontal-pod-autoscale/).

No personal standard covers Helm. Follow this plan.

## Locked decisions

- Chart path is `charts/mcp`. Create it only with `helm create` (from `charts/`, so the chart name is `mcp`). Then edit the generated files. Do not hand-write the skeleton.
- HPA stays. `helm create` already emits `templates/hpa.yaml`. Enable it by default: `minReplicas: 1`, `maxReplicas: 3`. Keep the generated CPU target. Set CPU and memory requests; a CPU utilization target does nothing without them. When autoscaling is enabled, omit Deployment `replicas` so the HPA owns the count. `replicaCount` is only the fallback when autoscaling is turned off. Do not `fail` the template for a replica count other than 1. Do not add a PodDisruptionBudget. Do not add `EventStore` or any other server change to make sessions cluster-wide.
- Service `sessionAffinity: ClientIP` (timeout 10800). Sessions are in-memory. Stickiness is the chart mitigation. It is not a guarantee if the client IP changes or the pod is replaced. Say that in NOTES and the README. Do not refuse the HPA because of it.
- Do not set container `command` or `args`. The image entrypoint is the contract. A Helm-created nginx command will break the pod.
- Service is ClusterIP. Container port, service port, and targetPort are `8080`, name `http`. Probes are `httpGet` `/healthz`. Distroless has no shell; do not use `exec` probes.
- Ingress stays disabled. No TLS, no cert-manager, no auth. Same threat model as `docs/DECISIONS.md`: anyone who can open the port can read every resource.
- Pod security matches the image: `runAsNonRoot`, `runAsUser`/`runAsGroup` `65532`, `readOnlyRootFilesystem`, `allowPrivilegeEscalation: false`, drop all capabilities, `seccompProfile: RuntimeDefault`. `automountServiceAccountToken: false` on the app pod. The app has no egress need.
- Do not set `MCP_ASSET_ROOT` in the chart. The image default is the baked `assets/` tree. A wrong mount would hide those files.
- Default image repository `mcp`, tag `dev`, `pullPolicy` `IfNotPresent`. Local Make and CI override tag `dev` and `pullPolicy` `Never`. Do not hardcode a GCP project in git. The release job rewrites a temp copy of the values before `helm package` so the attached chart defaults to `gcr.io/${GCP_PROJECT_ID}/mcp` and the release tag. Git `values.yaml` stays `mcp` / `dev`.
- Tests that need a cluster are `//go:build e2e` under `test/e2e`. `make test` / `make check` / pre-commit must not start a cluster. Add an untagged `doc.go` so the package still typechecks without the tag.
- No new Go dependencies. Reuse `github.com/modelcontextprotocol/go-sdk/mcp` the way `cmd/http_handler_test.go` does. No client-go, no kubectl image, no Helm unittest plugin.
- Protocol e2e and endpoint acceptance run inside the cluster via `helm test`. Workload acceptance (securityContext, ready replicas, endpoints, restart count) runs on the host in the same Make target, after `helm test`. Do not weaken that split to avoid a second image.
- Prefer overlap with unit tests. A live check that repeats a unit assertion is required, not redundant.
- Namespace and Helm release name are both `mcp`. Local minikube profile is `mcp` (do not touch the default profile). CI uses whatever profile `medyagh/setup-minikube` already started. `make e2e` must accept that and must not start a second cluster.

## Chart

After `helm create`, keep `_helpers.tpl`, `NOTES.txt`, ServiceAccount, Service, Deployment, HPA, Ingress (disabled), and `templates/tests/`. Update NOTES with the Service DNS, `/healthz`, `/mcp`, HPA bounds, and the session-affinity limit. Remove nginx-specific values (config, probes paths, command).

Values to keep configurable: image repository/tag/pullPolicy, `autoscaling.enabled` / `minReplicas` / `maxReplicas` / CPU target, service port, `sessionAffinity`, resources, probe timings, `tests.enabled`, test image repository/tag/pullPolicy. Defaults: autoscaling on, min 1, max 3. Ship conservative resources (request 50m/64Mi, limit 500m/256Mi).

Add a NetworkPolicy: ingress to the app on 8080 from the same namespace; app egress denied. Allow the test pod DNS plus the app port, or the test pod will not resolve the Service. If the policy blocks `helm test`, fix the policy. Do not delete it to go green.

Replace the generated wget test with one Pod (or Job) annotated `helm.sh/hook: test`. Hook delete policy: `before-hook-creation,hook-succeeded` only, so failures remain for debug. `helm test --logs` must show the Go test log. The test container uses the e2e image, not the app image. Inject the base URL from the chart (`http://<fullname>:8080`), not a hardcoded host. `restartPolicy: Never`. Test ServiceAccount may automount a token only if a later test needs the API; the planned suite does not, so leave it false.

`helm lint --strict` and `helm template` must succeed with defaults. Rendered output must include an HPA with min 1 and max 3, and must not set Deployment `replicas` while autoscaling is enabled. A render with `autoscaling.enabled=false` and `replicaCount=2` must succeed and set replicas to 2. That path is not the production default.

## E2e image and suite

`test/e2e/Dockerfile`, context repo root. Multi-stage: `golang:1.26.6-bookworm` (same as `Dockerfile`), `CGO_ENABLED=0`, `go test -c -tags e2e -o /out/e2e.test ./test/e2e`, then `gcr.io/distroless/static-debian12:nonroot`. Entrypoint runs that binary with `-test.v`. `.dockerignore` already drops `.git` and `bin`; do not exclude `test/`.

The suite reads `MCP_E2E_BASE_URL`. It must fail fast if unset. Use the SDK client against `$MCP_E2E_BASE_URL/mcp`. Plain `net/http` for non-MCP checks.

Required cases, mapped to existing tests:

| Live check | Copy the assertion from |
|------------|-------------------------|
| `GET /healthz` is 200 and body `ok\n` | `cmd/http_handler_test.go` |
| Connect, list resources, resource count > 0, `standards://git/commit-messages` present, `text/markdown` | `TestHTTPHandlerListsResourcesAndHealthz`, `test_resource.py` |
| Read that URI; body contains its frontmatter `uri:` line | `TestBootstrapServer_Integration` |
| Read `standards://test/does-not-exist` returns an error | `TestBootstrapServer_Integration_UnknownResource` |
| List tools is exactly the four names in `internal/tools/registry.go` | `TestRegisterAll` plus the name constants |
| `discover_standards` returns only `standards://` URIs, including `standards://git/commit-messages` | `TestBootstrapServer_CallTool_DiscoverStandards` |
| `discover_workflows` returns only `workflows://` URIs and is non-empty | same shape as discover standards |
| `search_standards` query `commit` matches at least one `standards://` URI | `TestBootstrapServer_CallTool_SearchStandards` |
| `search_standards` query `nonexistent-token-xyz` returns an empty resource list, not a transport error | same test, case 4 |
| `search_workflows` with `query: ""` does not error and returns at least one workflow | `TestBootstrapServer_CallTool_SearchWorkflows_EmptyQuery` |
| A second client connects while the first session is open | stateful handler in `cmd/main.go` |
| `POST /mcp` with a garbage body does not kill the process: `/healthz` still returns `ok\n` afterward | acceptance of the running server |
| `GET /` is not 200 with body `ok\n` | only `/healthz` is the probe |
| Initialize reports the server name from `internal/server.go` | bootstrap contract |

Do not re-test frontmatter parser tables, stopword tables, or bootstrap error paths. Those need a broken assets tree and must not be aimed at the production image.

## Makefile

Add targets, each with a `##` help line:

| Target | Behavior |
|--------|----------|
| `helm-lint` | `helm lint --strict charts/mcp` and `helm template` into a temp file. No cluster. |
| `e2e` | The one command the operator asked for. |
| `e2e-down` | `helm uninstall` the release, delete namespace `mcp`, then `minikube delete -p <profile>` only for the profile this target started. Do not delete a cluster CI or the operator already had running. |

`e2e` must:

1. Fail with a clear message if `docker`, `minikube`, `helm`, or `kubectl` is missing.
2. Profile: `MINIKUBE_PROFILE` if set, otherwise the running profile, otherwise `mcp`. Start only when that profile is not already `Running`. Do not pass a driver if the profile exists. Do not start a second cluster when `medyagh/setup-minikube` already started one.
3. `minikube -p "$PROFILE" addons enable metrics-server`. Do not fail the run if HPA metrics are still `<unknown>` before the suite ends. The HPA object must still be present.
4. `docker build -t mcp:dev .` and `docker build -f test/e2e/Dockerfile -t mcp-e2e:dev .`
5. Load both tags into that profile. Prefer `minikube image load`. If the GitHub docker driver cannot see host images, build in a subshell with `minikube docker-env` so the env does not leak into later recipes. `pullPolicy` stays `Never`.
6. `kubectl` context is that profile. Create namespace `mcp` if needed.
7. `helm upgrade --install mcp charts/mcp --namespace mcp --wait --timeout 5m` with image tags `dev`, `pullPolicy=Never`, and the test image `mcp-e2e:dev`. Do not disable the HPA.
8. Wait until the app rollout is complete. Idle ready replicas must be 1 (the HPA minimum), not 3.
9. `helm test mcp --namespace mcp --logs --timeout 10m`
10. Assert with kubectl: HPA `minReplicas` is 1 and `maxReplicas` is 3, Service `sessionAffinity` is `ClientIP`, container user is 65532, `readOnlyRootFilesystem` is true, Service has ready endpoints, app container `restartCount` is 0. Do not generate load to force a scale-up. That is flaky on a laptop and on a runner.
11. Leave the cluster up. Say how to tear down only the profile this command started.

Do not add `e2e` or `helm-lint` to `make check` or to pre-commit. `make check` must stay runnable without Docker.

Extend `fmt` and `fmt-check` so `//go:build e2e` files are included (`gofmt` / `go fmt -tags e2e ./test/e2e`). Extend `vet` and `lint` the same way for `./test/e2e` only, so default `./...` lint does not start requiring the tag. `make test` stays `go test -v ./...` with no tag.

Exclude `charts/mcp/templates/` from the pre-commit `check-yaml` hook. Keep `Chart.yaml` and `values.yaml` checked. Helm templates are not valid YAML.

## Pull request CI

Minikube is supported on GitHub-hosted Ubuntu. Do not skip this job. Official path: [setup minikube in GitHub Actions](https://minikube.sigs.k8s.io/docs/tutorials/setup_minikube_in_github_actions/). Docker driver. No nested virtualization. No kind, k3d, or a hand-rolled cluster.

Add one job to `.github/workflows/ci.yml`, same `pull_request` and `push` triggers as the other jobs. `ubuntu-latest` only. Not in the Go version matrix. `fail-fast` of the other jobs must not change.

Steps:

1. `actions/checkout@v4`.
2. `medyagh/setup-minikube` pinned to a release tag, not `@latest` or `@master`. That action starts the cluster. Record its profile and pass it as `MINIKUBE_PROFILE`.
3. Install Helm with `azure/setup-helm` pinned to a major version. The minikube action does not install Helm. `kubectl` comes from the minikube action.
4. `make e2e`.

Do not duplicate the helm install or the test list in the workflow. If `make e2e` and the action disagree about the profile, fix the Makefile. Do not add a second workflow file.

The existing Docker smoke job stays. It does not install the chart.

## Release bundle

Edit `.github/workflows/release.yml` only. Same trigger: GitHub Release `published`. After the image push succeeds, package `charts/mcp` and attach the `.tgz` to that same Release.

- `contents: write` is required for the upload. The workflow is `contents: read` today. Grant write on the job that uploads, not on unrelated jobs.
- Install Helm the same way as CI.
- Version: strip a leading `v` from `${{ github.event.release.tag_name }}`. `helm package --version` and `--app-version` use that value.
- Before packaging, rewrite a temp copy of the chart values. Do not commit it. Set `image.repository` to `gcr.io/${{ secrets.GCP_PROJECT_ID }}/mcp`, `image.tag` to that version, `image.pullPolicy` to `IfNotPresent`. Leave HPA defaults (min 1, max 3) and `tests.enabled` as in git. Do not push the e2e image. `helm install` does not pull the test image; `helm test` does.
- `gh release upload` that `.tgz` onto the tag that triggered the workflow, using `GH_TOKEN: ${{ github.token }}`. Fail the job if the upload fails. Do not overwrite an existing asset silently; if the asset name already exists, stop and report.
- Do not push the chart to `gcr.io`. This repo's registry is Container Registry, and Helm OCI charts belong in Artifact Registry, which this workflow does not use. Do not add a `gh-pages` Helm repo, Chart Releaser, or a new GCP API. The GitHub Release asset is the bundle.
- The image job's tags (semver, major.minor, `latest`, `sha-`) stay as they are. The chart's default tag is the full release version only.

## Docs

Separate doc commit:

- `README.md`: prerequisites (`docker`, `minikube`, `helm`, `kubectl`), `make e2e`, `make e2e-down`, `make helm-lint`, HPA defaults, session-affinity limit, and that a published GitHub Release carries both the image and the chart `.tgz`.
- `.agents/rules/validation.md`: new targets. `make e2e` is not part of `make check` or pre-commit. It does run in CI.
- `docs/DECISIONS.md`: one bullet — Kubernetes deploy is this chart; HPA min 1 max 3; ClientIP affinity because sessions are in-memory; probes use `/healthz`; no in-process TLS or auth; chart bundles attach to the GitHub Release, not gcr.io.
- `docs/ARCHITECTURE.md`: one line that HTTP on a cluster is `charts/mcp`, not a second binary.

Do not edit `.agents/skills/server-development`.

## Git and PR

Work only in a new worktree. From this repo, after `git fetch` and an up-to-date `main`:

`git worktree add -b helm-http-chart ../mcp-helm-http-chart main`

Commit there. Suggested split, still subject to the 500-line cap (split further if a commit would exceed it):

1. `feat: add helm chart for the HTTP server` — `charts/mcp` only.
2. `test: add in-cluster HTTP e2e suite` — `test/e2e` including its Dockerfile.
3. `build: add minikube helm e2e make target` — `Makefile`, `.pre-commit-config.yaml`.
4. `ci: run minikube e2e and attach the chart to the release` — `.github/workflows/ci.yml`, `.github/workflows/release.yml`.
5. `docs: document helm chart and e2e target` — README, validation rule, architecture, decisions.

Do not commit `bin/`, `coverage.out`, or rendered manifests.

Push `helm-http-chart` with a normal `git push -u origin helm-http-chart`. Open a PR with `gh`. Title must list every type present, in the priority order in `assets/standards/git/pull-requests.md` (expected: `feat, test, build, ci, docs: serve the MCP server on Kubernetes`). Description uses Added / Removed / Updated / Deleted, including `- None` where empty. Link this plan.

## Verify before the PR

Run from the worktree. Stop on the first failure. This matches `.agents/rules/validation.md` and the CI `checks` / `precommit` jobs. Then run the cluster command. CI runs the same `make e2e` on Ubuntu; a green local run is still required before the PR.

```bash
make fmt
make fmt-check
make typecheck
make vet
make test
make lint
make vuln
make check
pre-commit run --all-files
make helm-lint
make build
python3 .agents/skills/test-resource-document/scripts/test_resource.py
python3 .agents/skills/test-resource-document/scripts/test_resource.py --read standards://git/commit-messages
make e2e
helm package charts/mcp --version 0.0.0-test --app-version 0.0.0-test
```

`make test` must still pass with zero e2e cases executed. `make e2e` must show the helm test pod exit 0, idle replicas at 1, and an HPA with min 1 and max 3. If Docker or minikube cannot start, stop and report. Do not open the PR. `helm package` must produce a `.tgz`. Do not commit it. Do not publish a GitHub Release to test the upload.

Also build the app image the way CI does (`docker build -t mcp:ci .`) and `curl` `/healthz` on a published port, matching the smoke step in `.github/workflows/ci.yml`. The chart must not be required for that smoke test to keep passing.
