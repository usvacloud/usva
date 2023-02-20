package authenticator

import "context"

// Authenticator's job is to control access to specific resources.
// For example, one could:
// - implement Authenticator interface for accounts, where
//   - `AuthenticateSession` is used for session verification
//   - `NewSession` is used for login

type Authenticator[S, U, L, R any] interface {
	Authenticate(context.Context, S) (U, error)
	Register(context.Context, R) (S, error)
	NewSession(context.Context, L) (S, error)
}
