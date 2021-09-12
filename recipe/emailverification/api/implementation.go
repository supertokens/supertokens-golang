package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() evmodels.APIInterface {
	return evmodels.APIInterface{
		VerifyEmailPOST: func(token string, options evmodels.APIOptions) (evmodels.VerifyEmailUsingTokenResponse, error) {
			return options.RecipeImplementation.VerifyEmailUsingToken((token))
		},

		IsEmailVerifiedGET: func(options evmodels.APIOptions) (evmodels.IsEmailVerifiedGETResponse, error) {
			session, err := session.GetSession(options.Req, options.Res, nil)
			if err != nil {
				return evmodels.IsEmailVerifiedGETResponse{}, err
			}
			if session == nil {
				return evmodels.IsEmailVerifiedGETResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
			}

			userID := session.GetUserID()

			email, err := options.Config.GetEmailForUserID(userID)
			if err != nil {
				return evmodels.IsEmailVerifiedGETResponse{}, err
			}
			isVerified, err := options.RecipeImplementation.IsEmailVerified(userID, email)
			if err != nil {
				return evmodels.IsEmailVerifiedGETResponse{}, err
			}
			return evmodels.IsEmailVerifiedGETResponse{
				OK: &struct{ IsVerified bool }{
					IsVerified: isVerified,
				},
			}, nil
		},

		GenerateEmailVerifyTokenPOST: func(options evmodels.APIOptions) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {
			session, err := session.GetSession(options.Req, options.Res, nil)
			if err != nil {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
			}
			if session == nil {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
			}

			userID := session.GetUserID()
			email, err := options.Config.GetEmailForUserID(userID)
			if err != nil {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
			}
			response, err := options.RecipeImplementation.CreateEmailVerificationToken(userID, email)
			if err != nil {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
			}

			if response.EmailAlreadyVerifiedError != nil {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{
					EmailAlreadyVerifiedError: &struct{}{},
				}, nil
			}

			user := evmodels.User{
				ID:    userID,
				Email: email,
			}
			emailVerificationURL, err := options.Config.GetEmailVerificationURL(user)
			if err != nil {
				return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
			}
			emailVerifyLink := emailVerificationURL + "?token=" + response.OK.Token + "&rid=" + options.RecipeID

			options.Config.CreateAndSendCustomEmail(user, emailVerifyLink)

			return evmodels.GenerateEmailVerifyTokenPOSTResponse{
				OK: &struct{}{},
			}, nil
		},
	}
}
