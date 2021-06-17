package api

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/errors"
	sessionErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		RefreshPOST: func(options models.APIOptions) error {
			_, err := options.RecipeImplementation.RefreshSession(options.Req, options.Res)
			return err
		},
		VerifySession: func(verifySessionOptions *models.VerifySessionOptions, options models.APIOptions) error {
			method := options.Req.Method
			if method == http.MethodOptions || method == http.MethodTrace {
				options.OtherHandler(options.Res, options.Req)
				return nil
			}

			incomingPath, err := supertokens.NewNormalisedURLPath(options.Req.RequestURI)
			if err != nil {
				options.OtherHandler(options.Res, options.Req)
				return err
			}
			refreshTokenPath := options.Config.RefreshTokenPath

			if incomingPath.Equals(refreshTokenPath) && method == http.MethodPost {
				_, err := options.RecipeImplementation.RefreshSession(options.Req, options.Res)
				if err != nil {
					options.OtherHandler(options.Res, options.Req)
					return err
				}
			} else {
				_, err := options.RecipeImplementation.GetSession(options.Req, options.Res, verifySessionOptions)
				if err != nil {
					options.OtherHandler(options.Res, options.Req)
					return err
				}
			}
			options.OtherHandler(options.Res, options.Req)
			return nil
		},
		SignOutPOST: func(options models.APIOptions) (map[string]string, error) {
			session, err := options.RecipeImplementation.GetSession(options.Req, options.Res, nil)
			if err != nil {
				if sessionErrors.IsUnauthorizedError(err) {
					return map[string]string{
						"status": "OK",
					}, nil
				}
				return nil, err
			}
			if session == nil {
				return nil, errors.BadInputError{Msg: "Session is undefined. Should not come here."}
			}
			session.RevokeSession()
			return map[string]string{
				"status": "OK",
			}, nil
		},
	}
}
