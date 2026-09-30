---
uri: standards://go/architecture
name: Go (Golang) Coding Architecture Practices
description: "Best practices for go code organisation and application architecture"
languages:
    - go
file_types:
    - "*.go"
priority: required
related_resources:
    - standards://go/syntax
    - standards://go/logging
    - standards://go/cli
    - standards://go/tooling
    - standards://go/testing
    - workflows://go-new-project
---

# Go (Golang) Coding Architecture Practices
Best practices for go code organisation and application architecture

- ALWAYS follow the repository when it already standardizes on different tools. Repository rules win.
- ALWAYS apply these rules when writing new Go. NEVER restyle an existing repository to match them unless asked.
- ALWAYS implement following idiomatic Go solutions.
- NEVER copy a six-layer CLI cake (`service → action → runtime → backend`) into every Go app. Add a layer ONLY when this repository already has that layer.

## Layout

- ALWAYS place application entrypoints under `cmd/<binary>/main.go`.
- ALWAYS keep `cmd/<binary>/main.go` a tiny composition and execute entrypoint. It builds the container, registers components, and calls `Execute`.
- NEVER put business logic or library types in `main.go`.
- ALWAYS place feature packages under `internal/<feature>/`.
- ALWAYS default to `internal/`. Create a top-level public package ONLY for an intentional reusable API.
- ALWAYS place compiled binaries in `bin/`. Build commands MUST output to `bin/`.
- ALWAYS place 3rd-party library source under `lib/` when the repository vendors it.
- ALWAYS place 3rd-party compiled binaries under `dist/` when the repository vendors them.
- See `standards://go/syntax` for file and package naming.

```
cmd/<binary>/main.go          # wire container, execute
internal/<feature>/
  user_repository.go          # one behavior-owning struct
  protocols.go                # consumer interfaces
  types.go                    # request/response and other DTOs
  errors.go                   # sentinel and typed errors
```

## Files in a feature package

- ALWAYS put one behavior-owning struct per file. Keep every method for that struct in the same file.
- NEVER split one receiver across sibling files such as `foo_create.go` and `foo_update.go`.
- IF a file is too large, ALWAYS split responsibilities into composed types. NEVER hide an oversized type behind sibling method files.
- ALWAYS group supporting kinds in the same package:
    - `protocols.go` — consumer interfaces
    - `errors.go` — error sentinels and error types
    - `types.go` — DTOs, requests, responses, and options structs
- ALWAYS name the file after the primary type, in snake_case. Rename the file when you rename the type.
- NEVER let `port_forward_service.go` hold `TunnelManager`.
- Files with only functions, constants, or vars and no primary type (`errors.go`, `helpers.go`) are allowed.

## Interfaces

- ALWAYS return structs and accept interfaces. Constructors return `*Foo` or `(*Foo, error)`.
- NEVER return an interface from a constructor.
- ALWAYS put consumer interfaces beside the consumer, in that package's `protocols.go`. The consumer owns the interface.
- PREFER fewer interfaces, even if each one is slightly larger. A second interface means a separate role.
- NEVER pass one concrete object through many tiny consumer interfaces.

```go
// ✅ GOOD — constructor returns the struct
func NewProcessor() *Processor {
    return &Processor{}
}

// ❌ BAD — constructor returns the interface
func NewProcessor() ProcessorInterface {
    return &Processor{}
}

// ✅ GOOD — consumer owns a focused contract in protocols.go
type Forwarder interface {
    Forward(ctx context.Context, namespace, pod string, localPort, remotePort int) error
}

// ❌ BAD — constructor returns the interface
func NewPortForwardService(...) Forwarder { ... }
```

- AVOID interface assertions. IF one is used, ALWAYS include a comment explaining why. These assertions have no application-code value and are confusing without that comment.

```go
// ✅ GOOD — comment explains the compile-time check
// Compile-time assertion that *Client satisfies flyerfinder.FlyerFinder.
// If *Client stops implementing FlyerFinder, this line fails here.
var _ flyerfinder.FlyerFinder = (*Client)(nil)

// ❌ BAD — assertion with no explanation
var _ flyerfinder.FlyerFinder = (*Client)(nil)
```

## Constructors

- ALWAYS prefix constructor functions with `New`. Example: `NewRuntime()`, `NewStore()`, `NewFlyerFinder()`.
- ALWAYS give a behavior-owning struct that has dependencies a constructor.
- Simple data and value types MAY use a usable zero value (`var buf bytes.Buffer`).
- ALWAYS pass direct parameters until there are more than 4. THEN pass one explicit typed config struct.
- NEVER use functional options or config methods on a constructor. No `NewFoo(WithTimeout(d))`, no `...Option`, no `func(*Foo)` passed into `New`.
- ALWAYS put optional fields on the config struct, or as extra explicit parameters while the count is still 4 or fewer.
- Constructors are the ONLY dependency-injection point for a struct.
- NEVER read globals or package vars inside methods to obtain clients, configs, or loggers.
- A constructor assigns fields and MAY check local invariants (nil deps, empty required strings).
- NEVER do network or file I/O in a constructor. Validate the selected component when it is actually used.
- NEVER add singleton globals or `sync.Once` "one instance" types. The container holds the one instance for the process.

```go
// ✅ GOOD — four or fewer explicit params
func NewPortForwardService(config *rest.Config, logger *zap.SugaredLogger) *PortForwardService {
    return &PortForwardService{config: config, logger: logger}
}

// ✅ GOOD — more than four params: one typed struct, not option funcs
type UserRepositoryConfig struct {
    DB      DB
    Cache   Cache
    Clock   Clock
    Timeout time.Duration
    Limit   int
}

func NewUserRepository(cfg UserRepositoryConfig) (*UserRepository, error) {
    if cfg.DB == nil {
        return nil, fmt.Errorf("user repository: db is nil")
    }
    return &UserRepository{
        db: cfg.DB, cache: cfg.Cache, clock: cfg.Clock,
        timeout: cfg.Timeout, limit: cfg.Limit,
    }, nil
}

// ❌ BAD — functional options
func NewUserRepository(opts ...Option) *UserRepository { ... }

// ❌ BAD — missing New prefix
func CreateProcessor() *Processor { return &Processor{} }
```

## Container and wiring

- ALWAYS give anything beyond a tiny script a centralized container. The container, or composition root, owns registration.
- NEVER let a component self-register.
- NEVER write factory functions or factory switches that grow an `if` or `switch` for every new implementation.
- ALWAYS register implementations on the container. Callers look them up by key.
- ALWAYS verify new code is wired into the rest of the application.
- NEVER open a pull request with unwired code. Unused constructors, unregistered tools, and structs nothing reads are suspicious.
- IF the plan asked for a stub, show that. Otherwise confirm with the user.

## init()

- `init()` is allowed ONLY where the framework or library already uses it. Cobra command registration in a Cobra project is the allowed case.
- NEVER add a new `init()` in application or library code.
- IF setup seems to require `init()`, stop and review the design. Wire it from the container or an explicit constructor instead.
- See `standards://go/cli` and `workflows://go-cobra-subcommand` for the Cobra registration exception.

## Methods and helpers

- NEVER create a method for its own sake.
- ALWAYS write inline first.
- Extract a private helper ONLY when the logic is used more than once, OR when leaving it inline would make the call site much harder to read.
- IF every use is on one struct, the helper is an unexported method on that struct.
- Cross-type helpers in one feature go in the nearest `utils.go` or `utils` package.
- IF helpers are used across features, ALWAYS move them to a specifically named package.
- NEVER grow a global dumping-ground `utils`.
- NEVER import a helper from a sibling file just to reuse it. Move shared logic to a proper home first.
- Interface methods and exported API methods are contracts, not helpers. They MAY exist even when called once.
- Package-level functions are for constructors (`NewFoo`) and genuine utilities with no owner.
- NEVER leave floating methods with no owner.
- NEVER reshape production code for tests. NEVER extract or export a helper only so a test can call it.

## Embedding

- Anonymous embedding is allowed WHEN the parent name states what the combination is.
- NEVER embed a collaborator from another layer. Call it through a named field.

```go
// ✅ GOOD — parent name is the combination
type ReadWriter interface {
    io.Reader
    io.Writer
}

// ❌ BAD — collaborator from another layer promoted onto the parent
type UserService struct {
    *APIClient
}

// ✅ GOOD — named field, parent exposes its own methods
type UserService struct {
    client *APIClient
}
```

## Variable passing

- NEVER pass pointers into method parameters with the intention of filling them for the caller to use afterwards.
- ALWAYS treat variables, whether passed by value or by reference, as pass-by-value. Return the new value.

```go
// ❌ BAD — result is filled by solveMath and then printed
var result *string
err := solveMath(result)
fmt.Print(result)

// ✅ GOOD — treated as pass-by-value
var resultA *string
err, resultB := solveMath(resultA)
fmt.Print(resultB)
```

## HTTP client

- PREFER [gentleman](https://pkg.go.dev/github.com/h2non/gentleman) for new HTTP clients.
- Use `net/http` WHEN gentleman would only wrap a transport you still have to build, such as a custom `DialContext`.
- ALWAYS follow the repository if it already uses another HTTP client.

```go
res, err := gentleman.New().
    BaseURL("https://api.example.com").
    Request().Get().Path("/users").Send()
```
