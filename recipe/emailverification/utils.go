package emailverification

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config schema.TypeInput) schema.TypeNormalisedInput {
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

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) schema.TypeNormalisedInput {
	return schema.TypeNormalisedInput{
		GetEmailForUserID:        func(userID string) (string, error) { return "", errors.New("Not defined by user") },
		GetEmailVerificationURL:  DefaultGetEmailVerificationURL(appInfo),
		CreateAndSendCustomEmail: DefaultCreateAndSendCustomEmail(appInfo),
		Override: struct {
			Functions func(originalImplementation schema.RecipeImplementation) schema.RecipeImplementation
			APIs      func(originalImplementation schema.APIImplementation) schema.APIImplementation
		}{
			Functions: func(originalImplementation schema.RecipeImplementation) schema.RecipeImplementation {
				return originalImplementation
			},
			APIs: func(originalImplementation schema.APIImplementation) schema.APIImplementation {
				return originalImplementation
			},
		},
	}
}
