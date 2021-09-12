package emailverification

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config evmodels.TypeInput) evmodels.TypeNormalisedInput {
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

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) evmodels.TypeNormalisedInput {
	return evmodels.TypeNormalisedInput{
		GetEmailForUserID:        func(userID string) (string, error) { return "", errors.New("not defined by user") },
		GetEmailVerificationURL:  DefaultGetEmailVerificationURL(appInfo),
		CreateAndSendCustomEmail: DefaultCreateAndSendCustomEmail(appInfo),
		Override: evmodels.OverrideStruct{
			Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation evmodels.APIInterface) evmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}
