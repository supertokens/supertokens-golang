package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() models.APIInterface {
	return models.APIInterface{
		VerifyEmailPOST: func(token string, options models.APIOptions) (models.VerifyEmailUsingTokenResponse, error) {
			return options.RecipeImplementation.VerifyEmailUsingToken((token))
		},

		IsEmailVerifiedGET: func(options models.APIOptions) (models.IsEmailVerifiedGETResponse, error) {
			session, err := session.GetSession(options.Req, options.Res, nil)
			if err != nil {
				return models.IsEmailVerifiedGETResponse{}, err
			}
			if session == nil {
				return models.IsEmailVerifiedGETResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
			}

			userID := session.GetUserID()

			email, err := options.Config.GetEmailForUserID(userID)
			if err != nil {
				return models.IsEmailVerifiedGETResponse{}, err
			}
			isVerified, err := options.RecipeImplementation.IsEmailVerified(userID, email)
			if err != nil {
				return models.IsEmailVerifiedGETResponse{}, err
			}
			return models.IsEmailVerifiedGETResponse{
				OK: &struct{ IsVerified bool }{
					IsVerified: isVerified,
				},
			}, nil
		},

		GenerateEmailVerifyTokenPOST: func(options models.APIOptions) (models.GenerateEmailVerifyTokenPOSTResponse, error) {
			session, err := session.GetSession(options.Req, options.Res, nil)
			if err != nil {
				return models.GenerateEmailVerifyTokenPOSTResponse{}, err
			}
			if session == nil {
				return models.GenerateEmailVerifyTokenPOSTResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
			}

			userID := session.GetUserID()
			email, err := options.Config.GetEmailForUserID(userID)
			if err != nil {
				return models.GenerateEmailVerifyTokenPOSTResponse{}, err
			}
			response, err := options.RecipeImplementation.CreateEmailVerificationToken(userID, email)
			if err != nil {
				return models.GenerateEmailVerifyTokenPOSTResponse{}, err
			}

			if response.EmailAlreadyVerifiedError != nil {
				return models.GenerateEmailVerifyTokenPOSTResponse{
					EmailAlreadyVerifiedError: &struct{}{},
				}, nil
			}

			user := models.User{
				ID:    userID,
				Email: email,
			}
			emailVerificationURL, err := options.Config.GetEmailVerificationURL(user)
			if err != nil {
				return models.GenerateEmailVerifyTokenPOSTResponse{}, err
			}
			emailVerifyLink := emailVerificationURL + "?token=" + response.OK.Token + "&rid=" + options.RecipeID

			options.Config.CreateAndSendCustomEmail(user, emailVerifyLink)

			return models.GenerateEmailVerifyTokenPOSTResponse{
				OK: &struct{}{},
			}, nil
		},
	}
}
