package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
)

func MakeAPIImplementation() schema.APIImplementation {
	return schema.APIImplementation{
		VerifyEmailPOST: func(token string, options schema.APIOptions) (*schema.VerifyEmailUsingTokenResponse, error) {
			return options.RecipeImplementation.VerifyEmailUsingToken((token))
		},

		IsEmailVerifiedGET: func(options schema.APIOptions) (bool, error) {
			// TODO: session management
			userId := "TODO"
			email := "TODO"
			return options.RecipeImplementation.IsEmailVerified(userId, email)
		},

		GenerateEmailVerifyTokenPOST: func(options schema.APIOptions) (*schema.CreateEmailVerificationTokenAPIResponse, error) {
			// TODO: session management
			userId := "TODO"
			email := "TODO"
			response, err := options.RecipeImplementation.CreateEmailVerificationToken(userId, email)

			if err != nil {
				return nil, err
			}

			if response.EmailAlreadyVerifiedError == true {
				return &schema.CreateEmailVerificationTokenAPIResponse{
					EmailAlreadyVerifiedError: true,
				}, nil
			}

			emailVerifyLink := options.Config.GetEmailVerificationURL(schema.User{ID: userId, Email: email}) +
				"?token=" + response.OK.Token + "&rid=" + options.RecipeID

			options.Config.CreateAndSendCustomEmail(schema.User{ID: userId, Email: email}, emailVerifyLink)

			return &schema.CreateEmailVerificationTokenAPIResponse{
				OK: true,
			}, nil
		},
	}
}
