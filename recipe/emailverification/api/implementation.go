package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		VerifyEmailPOST: func(token string, options models.APIOptions) (*models.VerifyEmailUsingTokenResponse, error) {
			return options.RecipeImplementation.VerifyEmailUsingToken((token))
		},

		IsEmailVerifiedGET: func(options models.APIOptions) (bool, error) {
			session, err := session.GetSession(options.Req, options.Res, nil)
			if err != nil {
				return false, err
			}
			if session == nil {
				return false, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
			}
			userID := session.GetUserID()
			email, err := options.Config.GetEmailForUserID(userID)
			if err != nil {
				return false, err
			}
			return options.RecipeImplementation.IsEmailVerified(userID, email)
		},

		GenerateEmailVerifyTokenPOST: func(options models.APIOptions) (*models.CreateEmailVerificationTokenAPIResponse, error) {
			session, err := session.GetSession(options.Req, options.Res, nil)
			if err != nil {
				return nil, err
			}
			if session == nil {
				return nil, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
			}

			userID := session.GetUserID()
			email, err := options.Config.GetEmailForUserID(userID)
			if err != nil {
				return nil, err
			}
			response, err := options.RecipeImplementation.CreateEmailVerificationToken(userID, email)
			if err != nil {
				return nil, err
			}

			if response.EmailAlreadyVerifiedError {
				return &models.CreateEmailVerificationTokenAPIResponse{
					EmailAlreadyVerifiedError: true,
				}, nil
			}
			emailVerificationURL, err := options.Config.GetEmailVerificationURL(models.User{ID: userID, Email: email})
			if err != nil {
				return nil, err
			}
			emailVerifyLink := emailVerificationURL + "?token=" + response.Ok.Token + "&rid=" + options.RecipeID

			options.Config.CreateAndSendCustomEmail(models.User{ID: userID, Email: email}, emailVerifyLink)
			return &models.CreateEmailVerificationTokenAPIResponse{
				OK: true,
			}, nil
		},
	}
}
