package thirdparty

import (
	"encoding/json"
	"errors"

	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance Recipe, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (models.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	signInAndUpFeature, err := validateAndNormaliseSignInAndUpConfig(config.SignInAndUpFeature)
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}
	typeNormalisedInput.SignInAndUpFeature = signInAndUpFeature

	typeNormalisedInput.EmailVerificationFeature = validateAndNormaliseEmailVerificationConfig(recipeInstance, config)

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

func makeTypeNormalisedInput(recipeInstance Recipe) models.TypeNormalisedInput {
	return models.TypeNormalisedInput{
		SignInAndUpFeature:       models.TypeNormalisedInputSignInAndUp{},
		EmailVerificationFeature: validateAndNormaliseEmailVerificationConfig(recipeInstance, nil),
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

func validateAndNormaliseEmailVerificationConfig(recipeInstance Recipe, config *models.TypeInput) evm.TypeInput {
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
						return "", errors.New("unknown User ID provided")
					}
					return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo)
				}
			}
		}
	}

	return emailverificationTypeInput
}

func validateAndNormaliseSignInAndUpConfig(config models.TypeInputSignInAndUp) (models.TypeNormalisedInputSignInAndUp, error) {
	providers := config.Providers
	if len(providers) == 0 {
		return models.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "thirdparty recipe requires at least 1 provider to be passed in signInAndUpFeature.providers config"}
	}
	return models.TypeNormalisedInputSignInAndUp{
		Providers: providers,
	}, nil
}

func parseUser(value interface{}) (*models.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func parseUsers(value interface{}) ([]models.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user []models.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
