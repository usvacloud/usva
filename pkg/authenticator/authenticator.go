package authenticator

import "context"

// Authenticator's job is to control access to specific resources.
// For example, one could:
// - implement Authenticator interface for accounts using database
// - implement Authenticator interface for account using private keys
//
// S stands for session
// U stands for user (e.g. identifiable object)
// L stands for login (or authentication) credinteals
// R stands for Registration credinteals
type Authenticator[S, U, L, R any] interface {
	Authenticate(context.Context, S) (U, error)
	Register(context.Context, R) (S, error)
	NewSession(context.Context, L) (S, error)
}
