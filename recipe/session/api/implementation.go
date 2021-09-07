package api

import (
	defaultErrors "errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() models.APIInterface {
	return models.APIInterface{
		RefreshPOST: func(options models.APIOptions) error {
			_, err := options.RecipeImplementation.RefreshSession(options.Req, options.Res)
			return err
		},

		VerifySession: func(verifySessionOptions *models.VerifySessionOptions, options models.APIOptions) (*models.SessionContainer, error) {
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

		SignOutPOST: func(options models.APIOptions) (models.SignOutPOSTResponse, error) {
			session, err := options.RecipeImplementation.GetSession(options.Req, options.Res, nil)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					return models.SignOutPOSTResponse{
						OK: &struct{}{},
					}, nil
				}
				return models.SignOutPOSTResponse{}, err
			}
			if session == nil {
				return models.SignOutPOSTResponse{}, defaultErrors.New("session is nil. Should not come here.")
			}

			err = session.RevokeSession()
			if err != nil {
				return models.SignOutPOSTResponse{}, err
			}

			return models.SignOutPOSTResponse{
				OK: &struct{}{},
			}, nil
		},
	}
}
