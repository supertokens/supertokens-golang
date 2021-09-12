package interceptor

import (
	"context"

	"github.com/twitchtv/twirp"
)

// NewInterceptorMakeSmallHats builds an interceptor that modifies
// calls to MakeHat ignoring the request, and instead always making small hats.
func NewSuperTokensErrorHandlerInterceptor() twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return next(ctx, req)
		}
	}
}
