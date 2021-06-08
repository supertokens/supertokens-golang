package emailverification

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config schema.TypeInput) schema.TypeNormalisedInput {
	var typeNormalisedInput schema.TypeNormalisedInput
	typeNormalisedInput.GetEmailVerificationURL = *config.GetEmailVerificationURL
	if config.GetEmailVerificationURL == nil {
		typeNormalisedInput.GetEmailVerificationURL = GetEmailVerificationURL(appInfo)
	}
	typeNormalisedInput.CreateAndSendCustomEmail = *config.CreateAndSendCustomEmail
	if config.CreateAndSendCustomEmail == nil {
		typeNormalisedInput.CreateAndSendCustomEmail = CreateAndSendCustomEmail(appInfo)
	}

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = *config.Override.Functions
		} else {
			typeNormalisedInput.Override.Functions = func(originalImplementation schema.RecipeInterface) schema.RecipeInterface {
				return originalImplementation
			}
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = *config.Override.APIs
		} else {
			typeNormalisedInput.Override.APIs = func(originalImplementation schema.APIInterface) schema.APIInterface {
				return originalImplementation
			}
		}

	}

	typeNormalisedInput.GetEmailForUserID = config.GetEmailForUserID
	return typeNormalisedInput
}
