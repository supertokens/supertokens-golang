package thirdpartyemailpassword

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (models.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	if config != nil && config.SignUpFeature != nil {
		typeNormalisedInput.SignUpFeature = validateAndNormaliseSignUpConfig(config.SignUpFeature)
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

func makeTypeNormalisedInput(recipeInstance *Recipe) models.TypeNormalisedInput {
	return models.TypeNormalisedInput{
		SignUpFeature:                  validateAndNormaliseSignUpConfig(nil),
		Providers:                      nil,
		ResetPasswordUsingTokenFeature: nil,
		EmailVerificationFeature:       validateAndNormaliseEmailVerificationConfig(recipeInstance, nil),
		Override: struct {
			Functions                func(originalImplementation models.RecipeInterface) models.RecipeInterface
			APIs                     func(originalImplementation models.APIInterface) models.APIInterface
			EmailVerificationFeature *struct {
				Functions func(originalImplementation evm.RecipeInterface) evm.RecipeInterface
				APIs      func(originalImplementation evm.APIInterface) evm.APIInterface
			}
		}{
			Functions: func(originalImplementation models.RecipeInterface) models.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation models.APIInterface) models.APIInterface {
				return originalImplementation
			},
			EmailVerificationFeature: nil,
		},
	}
}

func validateAndNormaliseSignUpConfig(config *models.TypeInputSignUp) models.TypeNormalisedInputSignUp {
	if config != nil {
		return models.TypeNormalisedInputSignUp{
			FormFields: emailpassword.NormaliseSignUpFormFields(config.FormFields),
		}
	}
	return models.TypeNormalisedInputSignUp{
		FormFields: emailpassword.NormaliseSignUpFormFields(nil),
	}
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance *Recipe, config *models.TypeInput) evm.TypeInput {
	emailverificationTypeInput := evm.TypeInput{
		GetEmailForUserID: recipeInstance.getEmailForUserId,
		Override:          nil,
	}

	if config != nil {
		if config.Override != nil {
			emailverificationTypeInput.Override = config.Override.EmailVerificationFeature
		}
		if config.EmailVerificationFeature != nil {
			if config.EmailVerificationFeature.CreateAndSendCustomEmail != nil {
				emailverificationTypeInput.CreateAndSendCustomEmail = func(user evm.User, link string) {
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
				emailverificationTypeInput.GetEmailVerificationURL = func(user evm.User) (string, error) {
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

// used to convert FormField of thirdpartyemailpassword to emailpassword's
// TODO: why can't we just use ep's form field type in this recipe too?
func normalisedToType(normalisedformFields []epm.NormalisedFormField) []epm.TypeInputFormField {
	var formFields []epm.TypeInputFormField
	for _, normalisedformField := range normalisedformFields {
		formFields = append(formFields, epm.TypeInputFormField{
			ID:       normalisedformField.ID,
			Validate: normalisedformField.Validate,
			Optional: &normalisedformField.Optional,
		})
	}
	return formFields
}
