package thirdparty

import (
	"encoding/json"
	"errors"

	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance models.RecipeImplementation, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) (models.TypeNormalisedInput, error) {
	sessionFeature := validateAndNormaliseSessionFeatureConfig(nil)
	if config != nil {
		sessionFeature = validateAndNormaliseSessionFeatureConfig(config.SessionFeature)
	}
	emailVerificationFeature := validateAndNormaliseEmailVerificationConfig(recipeInstance, config)
	signInAndUpFeature, err := validateAndNormaliseSignInAndUpConfig(config.SignInAndUpFeature)
	if err != nil {
		return models.TypeNormalisedInput{}, err
	}

	typeNormalisedInput := models.TypeNormalisedInput{
		SessionFeature:           sessionFeature,
		EmailVerificationFeature: emailVerificationFeature,
		SignInAndUpFeature:       signInAndUpFeature,
	}
	typeNormalisedInput.Override.Functions = func(originalImplementation models.RecipeImplementation) models.RecipeImplementation {
		return originalImplementation
	}
	typeNormalisedInput.Override.APIs = func(originalImplementation models.APIImplementation) models.APIImplementation {
		return originalImplementation
	}
	typeNormalisedInput.Override.EmailVerificationFeature = nil

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

func defaultSetJwtPayloadForSession(User models.User, thirdPartyAuthCodeResponse interface{}, action string) map[string]interface{} {
	return nil
}

func defaultSetSessionDataForSession(User models.User, thirdPartyAuthCodeResponse interface{}, action string) map[string]interface{} {
	return nil
}

func validateAndNormaliseSessionFeatureConfig(config *models.TypeNormalisedInputSessionFeature) models.TypeNormalisedInputSessionFeature {
	normalisedInputSessionFeature := models.TypeNormalisedInputSessionFeature{
		SetJwtPayload:  defaultSetJwtPayloadForSession,
		SetSessionData: defaultSetSessionDataForSession,
	}
	if config != nil {
		if config.SetJwtPayload != nil {
			normalisedInputSessionFeature.SetJwtPayload = config.SetJwtPayload
		}
		if config.SetSessionData != nil {
			normalisedInputSessionFeature.SetSessionData = config.SetSessionData
		}
	}
	return normalisedInputSessionFeature
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance models.RecipeImplementation, config *models.TypeInput) evm.TypeInput {
	var emailverificationTypeInput evm.TypeInput
	emailverificationTypeInput.GetEmailForUserID = getEmailForUserId

	emailverificationTypeInput.Override = nil
	if config != nil && config.Override != nil {
		override := config.Override
		if override.EmailVerificationFeature != nil {
			emailverificationTypeInput.Override = override.EmailVerificationFeature
			if config.EmailVerificationFeature.CreateAndSendCustomEmail == nil {
				emailverificationTypeInput.CreateAndSendCustomEmail = nil
			} else {
				emailverificationTypeInput.CreateAndSendCustomEmail = func(user evm.User, link string) error {
					userInfo := recipeInstance.GetUserByID(user.ID)
					if userInfo == nil {
						return errors.New("Unknown User ID provided")
					}
					return config.EmailVerificationFeature.CreateAndSendCustomEmail(*userInfo, link)
				}
			}

			if config.EmailVerificationFeature.GetEmailVerificationURL == nil {
				emailverificationTypeInput.GetEmailVerificationURL = nil
			} else {
				emailverificationTypeInput.GetEmailVerificationURL = func(user evm.User) (string, error) {
					userInfo := recipeInstance.GetUserByID(user.ID)
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

func validateAndNormaliseSignInAndUpConfig(config models.TypeInputSignInAndUp) (models.TypeNormalisedInputSignInAndUp, error) {
	providers := config.Providers
	if len(providers) == 0 {
		return models.TypeNormalisedInputSignInAndUp{}, supertokens.BadInputError{Msg: "thirdparty recipe requires atleast 1 provider to be passed in signInAndUpFeature.providers config"}
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
