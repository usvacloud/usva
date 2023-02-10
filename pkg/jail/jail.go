package jail

import (
	"context"
)

type Ban interface {
	Ban(context.Context, string) error
}

type IsAuthorized interface {
	IsAuthorized(context.Context, string) (bool, error)
}

type Jail interface {
	Ban
	IsAuthorized
}
