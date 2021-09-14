package interceptor

import (
	"context"
	"errors"
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/session"
	sessionError "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/twitchtv/twirp"
)

func NewSuperTokensErrorHandlerInterceptor() twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			resp, err := next(ctx, req)
			if err != nil {
				if errors.As(err, &sessionError.TryRefreshTokenError{}) {
					return resp, twirp.NewError(twirp.Unauthenticated, "try refresh token")
				} else if errors.As(err, &sessionError.UnauthorizedError{}) {
					fmt.Println("in unauth!!")
					return resp, twirp.NewError(twirp.Unauthenticated, "unauthorized")
				} else if errors.As(err, &sessionError.TokenTheftDetectedError{}) {
					_, err = session.RevokeSession(err.(sessionError.TokenTheftDetectedError).Payload.SessionHandle)
					if err != nil {
						return resp, err
					}
					return resp, twirp.NewError(twirp.Unauthenticated, "token theft detected")
				}
			}
			return resp, err
		}
	}
}
