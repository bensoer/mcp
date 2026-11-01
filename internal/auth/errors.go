package auth

import "errors"

// ErrEmailExists indicates that registration was attempted for an email that is
// already registered.
var ErrEmailExists = errors.New("email already exists")
