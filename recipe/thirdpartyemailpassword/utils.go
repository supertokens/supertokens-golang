package thirdpartyemailpassword

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *tpepmodels.TypeInput) (tpepmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	if config != nil && config.SignUpFeature != nil {
		typeNormalisedInput.SignUpFeature = config.SignUpFeature
	}

	if config != nil && config.Providers != nil {
		typeNormalisedInput.Providers = config.Providers
	}

	typeNormalisedInput.EmailVerificationFeature = validateAndNormaliseEmailVerificationConfig(recipeInstance, config)

	if config != nil && config.ResetPasswordUsingTokenFeature != nil {
		typeNormalisedInput.ResetPasswordUsingTokenFeature = config.ResetPasswordUsingTokenFeature
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
		if config.Override.EmailVerificationFeature != nil {
			typeNormalisedInput.Override.EmailVerificationFeature = config.Override.EmailVerificationFeature
		}
	}

	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(recipeInstance *Recipe) tpepmodels.TypeNormalisedInput {
	return tpepmodels.TypeNormalisedInput{
		SignUpFeature:                  nil,
		Providers:                      nil,
		ResetPasswordUsingTokenFeature: nil,
		EmailVerificationFeature:       validateAndNormaliseEmailVerificationConfig(recipeInstance, nil),
		Override: tpepmodels.OverrideStruct{
			Functions: func(originalImplementation tpepmodels.RecipeInterface) tpepmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
				return originalImplementation
			},
			EmailVerificationFeature: nil,
		},
	}
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance *Recipe, config *tpepmodels.TypeInput) evmodels.TypeInput {
	emailverificationTypeInput := evmodels.TypeInput{
		GetEmailForUserID: recipeInstance.getEmailForUserId,
		Override:          nil,
	}

	if config != nil {
		if config.Override != nil {
			emailverificationTypeInput.Override = config.Override.EmailVerificationFeature
		}
		if config.EmailVerificationFeature != nil {
			if config.EmailVerificationFeature.CreateAndSendCustomEmail != nil {
				emailverificationTypeInput.CreateAndSendCustomEmail = func(user evmodels.User, link string) {
					userInfo, err := recipeInstance.RecipeImpl.GetUserByID(user.ID)
					if err != nil {
						return
					}
					if userInfo == nil {
						return
					}
					config.EmailVerificationFeature.CreateAndSendCustomEmail(*userInfo, link)
				}
			}

			if config.EmailVerificationFeature.GetEmailVerificationURL != nil {
				emailverificationTypeInput.GetEmailVerificationURL = func(user evmodels.User) (string, error) {
					userInfo, err := recipeInstance.RecipeImpl.GetUserByID(user.ID)
					if err != nil {
						return "", err
					}
					if userInfo == nil {
						return "", errors.New("Unknown User ID provided")
					}
					return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo)
				}
			}
		}
	}

	return emailverificationTypeInput
}
