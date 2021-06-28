package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config models.TypeInput) models.TypeNormalisedInput {
	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	if config.GetEmailVerificationURL != nil {
		typeNormalisedInput.GetEmailVerificationURL = config.GetEmailVerificationURL
	}

	if config.CreateAndSendCustomEmail != nil {
		typeNormalisedInput.CreateAndSendCustomEmail = config.CreateAndSendCustomEmail
	}

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}

	}

	if config.GetEmailForUserID != nil {
		typeNormalisedInput.GetEmailForUserID = config.GetEmailForUserID
	}
	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) models.TypeNormalisedInput {
	return models.TypeNormalisedInput{
		GetEmailForUserID:        func(userID string) (string, error) { return "", supertokens.BadInputError{Msg: "Not defined by user"} },
		GetEmailVerificationURL:  DefaultGetEmailVerificationURL(appInfo),
		CreateAndSendCustomEmail: DefaultCreateAndSendCustomEmail(appInfo),
		Override: struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
		}{
			Functions: func(originalImplementation models.RecipeImplementation) models.RecipeImplementation {
				return originalImplementation
			},
			APIs: func(originalImplementation models.APIImplementation) models.APIImplementation {
				return originalImplementation
			},
		},
	}
}
