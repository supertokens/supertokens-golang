package api

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() sessmodels.APIInterface {
	return sessmodels.APIInterface{
		RefreshPOST: func(options sessmodels.APIOptions) error {
			_, err := options.RecipeImplementation.RefreshSession(options.Req, options.Res)
			return err
		},

		VerifySession: func(verifySessionOptions *sessmodels.VerifySessionOptions, options sessmodels.APIOptions) (*sessmodels.SessionContainer, error) {
			method := options.Req.Method
			if method == http.MethodOptions || method == http.MethodTrace {
				return nil, nil
			}

			incomingPath, err := supertokens.NewNormalisedURLPath(options.Req.RequestURI)
			if err != nil {
				return nil, err
			}

			refreshTokenPath := options.Config.RefreshTokenPath
			if incomingPath.Equals(refreshTokenPath) && method == http.MethodPost {
				session, err := options.RecipeImplementation.RefreshSession(options.Req, options.Res)
				return &session, err
			} else {
				return options.RecipeImplementation.GetSession(options.Req, options.Res, verifySessionOptions)
			}
		},

		SignOutPOST: func(options sessmodels.APIOptions) (sessmodels.SignOutPOSTResponse, error) {
			session, err := options.RecipeImplementation.GetSession(options.Req, options.Res, nil)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					return sessmodels.SignOutPOSTResponse{
						OK: &struct{}{},
					}, nil
				}
				return sessmodels.SignOutPOSTResponse{}, err
			}
			if session == nil {
				return sessmodels.SignOutPOSTResponse{}, defaultErrors.New("session is nil. Should not come here.")
			}

			err = session.RevokeSession()
			if err != nil {
				return sessmodels.SignOutPOSTResponse{}, err
			}

			return sessmodels.SignOutPOSTResponse{
				OK: &struct{}{},
			}, nil
		},
	}
}
