package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		VerifyEmailPOST: func(token string, options models.APIOptions) (*models.VerifyEmailUsingTokenResponse, error) {
			return options.RecipeImplementation.VerifyEmailUsingToken((token))
		},

		IsEmailVerifiedGET: func(options models.APIOptions) (bool, error) {
			// TODO: session management
			userId := "TODO"
			email := "TODO"
			return options.RecipeImplementation.IsEmailVerified(userId, email)
		},

		GenerateEmailVerifyTokenPOST: func(options models.APIOptions) (*models.CreateEmailVerificationTokenAPIResponse, error) {
			// TODO: session management
			userId := "TODO"
			email := "TODO"
			response, err := options.RecipeImplementation.CreateEmailVerificationToken(userId, email)

			if err != nil {
				return nil, err
			}

			if response.EmailAlreadyVerifiedError == true {
				return &models.CreateEmailVerificationTokenAPIResponse{
					EmailAlreadyVerifiedError: true,
				}, nil
			}

			emailVerifyLink := options.Config.GetEmailVerificationURL(models.User{ID: userId, Email: email}) +
				"?token=" + response.Ok.Token + "&rid=" + options.RecipeID

			options.Config.CreateAndSendCustomEmail(models.User{ID: userId, Email: email}, emailVerifyLink)

			return &models.CreateEmailVerificationTokenAPIResponse{
				OK: true,
			}, nil
		},
	}
}
