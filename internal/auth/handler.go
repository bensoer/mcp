package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	maxRegisterBodySize = 1 << 10
	defaultIdentityRPS  = 10
	defaultRemoteRPS    = 1
)

// Handler wraps an HTTP handler with API key authentication and rate limits.
type Handler struct {
	next          http.Handler
	store         *Store
	identityLimit *Limiter
	remoteLimit   *Limiter
	logger        *zap.SugaredLogger
}

// NewHandler creates an authentication and rate-limiting HTTP handler.
func NewHandler(next http.Handler, store *Store, rateLimit int, logger *zap.SugaredLogger) *Handler {
	identityRPS := defaultIdentityRPS
	remoteRPS := defaultRemoteRPS
	if rateLimit > 0 {
		identityRPS = rateLimit
		remoteRPS = rateLimit
	}

	return &Handler{
		next:          next,
		store:         store,
		identityLimit: NewLimiter(identityRPS),
		remoteLimit:   NewLimiter(remoteRPS),
		logger:        logger,
	}
}

// ServeHTTP applies authentication and rate limits before delegating to the
// wrapped handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/register":
		h.serveRegister(w, r)
	case "/mcp":
		h.serveMCP(w, r)
	default:
		h.next.ServeHTTP(w, r)
	}
}

func (h *Handler) serveRegister(w http.ResponseWriter, r *http.Request) {
	if !h.allowRemote(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		h.writeForbidden(w)
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		h.writeForbidden(w)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRegisterBodySize)
	var request registerRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		h.writeForbidden(w)
		return
	}
	if err := ensureJSONEnd(decoder); err != nil {
		h.writeForbidden(w)
		return
	}

	email := normalizeEmail(request.Email)
	if email == "" {
		h.writeForbidden(w)
		return
	}
	apiKey, err := h.store.Register(email)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			h.logger.Infof("API key registration rejected for duplicate email %s", email)
		} else {
			h.logger.Errorf("API key registration failed: %v", err)
		}
		h.writeForbidden(w)
		return
	}

	h.logger.Infof("API key registered for email %s", email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(registerResponse{APIKey: apiKey})
}

func (h *Handler) serveMCP(w http.ResponseWriter, r *http.Request) {
	apiKey, ok := bearerToken(r.Header.Get("Authorization"))
	if !ok {
		h.serveUnauthorized(w, r)
		return
	}

	email, ok, err := h.store.Authenticate(apiKey)
	if err != nil {
		h.logger.Errorf("authenticate MCP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		h.serveUnauthorized(w, r)
		return
	}

	if !h.allowIdentity(w, email) {
		return
	}
	h.next.ServeHTTP(w, r)
}

func (h *Handler) serveUnauthorized(w http.ResponseWriter, r *http.Request) {
	if !h.allowRemote(w, r) {
		return
	}
	h.logger.Infof("MCP request rejected from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusUnauthorized)
}

func (h *Handler) allowRemote(w http.ResponseWriter, r *http.Request) bool {
	allowed, delay := h.remoteLimit.Allow(r.RemoteAddr)
	if allowed {
		return true
	}
	h.logger.Infof("request rate limited from %s", r.RemoteAddr)
	h.writeRateLimited(w, delay)
	return false
}

func (h *Handler) allowIdentity(w http.ResponseWriter, email string) bool {
	allowed, delay := h.identityLimit.Allow(email)
	if allowed {
		return true
	}
	h.logger.Infof("MCP request rate limited for email %s", email)
	h.writeRateLimited(w, delay)
	return false
}

func (h *Handler) writeRateLimited(w http.ResponseWriter, delay time.Duration) {
	retryAfter := max(1, int(math.Ceil(delay.Seconds())))
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	w.WriteHeader(http.StatusTooManyRequests)
}

func (h *Handler) writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"forbidden"}`))
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

func ensureJSONEnd(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("request body contains multiple JSON values")
		}
		return fmt.Errorf("decode trailing JSON: %w", err)
	}
	return nil
}
