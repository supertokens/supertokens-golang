package api

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() schema.APIImplementation {
	return schema.APIImplementation{
		RefreshPOST: func(options schema.APIOptions) error {
			_, err := options.RecipeImplementation.RefreshSession(options.Req, options.Res)
			return err
		},
		VerifySession: func(verifySessionOptions *schema.VerifySessionOptions, options schema.APIOptions) {
			method := options.Req.Method
			if method == http.MethodOptions || method == http.MethodTrace {
				options.OtherHandler(options.Res, options.Req)
				return
			}

			incomingPath, err := supertokens.NewNormalisedURLPath(options.Req.RequestURI)
			if err != nil {
				options.OtherHandler(options.Res, options.Req)
				return
			}
			refreshTokenPath := options.Config.RefreshTokenPath

			if incomingPath.Equals(refreshTokenPath) && method == http.MethodPost {

			} else {

			}
			options.OtherHandler(options.Res, options.Req)
			return
		},
		SignOutPOST: func(options schema.APIOptions) (map[string]string, error) {
			session, err := options.RecipeImplementation.GetSession(options.Req, options.Res, nil)
			if err != nil {
				// TODO: error check handle
				return nil, err
			}
			if session == nil {
				return nil, errors.New("Session is undefined. Should not come here.")
			}
			session.RevokeSession()
			return map[string]string{
				"status": "OK",
			}, nil
		},
	}
}
