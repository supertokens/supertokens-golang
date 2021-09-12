package thirdparty

import (
	"encoding/json"
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *tpmodels.TypeInput) (tpmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	signInAndUpFeature, err := validateAndNormaliseSignInAndUpConfig(config.SignInAndUpFeature)
	if err != nil {
		return tpmodels.TypeNormalisedInput{}, err
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

func makeTypeNormalisedInput(recipeInstance *Recipe) tpmodels.TypeNormalisedInput {
	return tpmodels.TypeNormalisedInput{
		SignInAndUpFeature:       tpmodels.TypeNormalisedInputSignInAndUp{},
		EmailVerificationFeature: validateAndNormaliseEmailVerificationConfig(recipeInstance, nil),
		Override: tpmodels.OverrideStruct{
			Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tpmodels.APIInterface) tpmodels.APIInterface {
				return originalImplementation
			},
			EmailVerificationFeature: nil,
		},
	}
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance *Recipe, config *tpmodels.TypeInput) evmodels.TypeInput {
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
						return "", errors.New("unknown User ID provided")
					}
					return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo)
				}
			}
		}
	}

	return emailverificationTypeInput
}

func validateAndNormaliseSignInAndUpConfig(config tpmodels.TypeInputSignInAndUp) (tpmodels.TypeNormalisedInputSignInAndUp, error) {
	providers := config.Providers
	if len(providers) == 0 {
		return tpmodels.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "thirdparty recipe requires at least 1 provider to be passed in signInAndUpFeature.providers config"}
	}
	return tpmodels.TypeNormalisedInputSignInAndUp{
		Providers: providers,
	}, nil
}

func parseUser(value interface{}) (*tpmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user tpmodels.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func parseUsers(value interface{}) ([]tpmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var user []tpmodels.User
	err = json.Unmarshal(respJSON, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
