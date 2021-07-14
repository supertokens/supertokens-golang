package emailpassword

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	emailverificationModels "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance models.RecipeImplementation, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) models.TypeNormalisedInput {
	sessionFeature := validateAndNormaliseSessionFeatureConfig(nil)
	signUpFeature := validateAndNormaliseSignupConfig(nil)
	resetPasswordUsingTokenFeature := validateAndNormaliseResetPasswordUsingTokenConfig(appInfo, signUpFeature, nil)

	if config != nil && config.SessionFeature != nil {
		sessionFeature = validateAndNormaliseSessionFeatureConfig(config.SessionFeature)
	}

	if config != nil && config.SignUpFeature != nil {
		signUpFeature = validateAndNormaliseSignupConfig(config.SignUpFeature)
	}

	signInFeature := validateAndNormaliseSignInConfig(signUpFeature)

	if config != nil && config.ResetPasswordUsingTokenFeature != nil {
		resetPasswordUsingTokenFeature = validateAndNormaliseResetPasswordUsingTokenConfig(appInfo, signUpFeature, config.ResetPasswordUsingTokenFeature)
	}

	emailVerificationFeature := validateAndNormaliseEmailVerificationConfig(recipeInstance, config)

	typeNormalisedInput := models.TypeNormalisedInput{
		SessionFeature:                 sessionFeature,
		SignUpFeature:                  signUpFeature,
		SignInFeature:                  signInFeature,
		ResetPasswordUsingTokenFeature: resetPasswordUsingTokenFeature,
		EmailVerificationFeature:       emailVerificationFeature,
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

	return typeNormalisedInput
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance models.RecipeImplementation, config *models.TypeInput) emailverificationModels.TypeInput {
	var emailverificationTypeInput emailverificationModels.TypeInput
	emailverificationTypeInput.GetEmailForUserID = getEmailForUserId

	emailverificationTypeInput.Override = nil
	if config != nil && config.Override != nil {
		override := config.Override
		if override.EmailVerificationFeature != nil {
			emailverificationTypeInput.Override = override.EmailVerificationFeature
		}
	}
	if config != nil && config.EmailVerificationFeature.CreateAndSendCustomEmail == nil {
		emailverificationTypeInput.CreateAndSendCustomEmail = nil
	} else {
		emailverificationTypeInput.CreateAndSendCustomEmail = func(user emailverificationModels.User, link string) error {
			userInfo := recipeInstance.GetUserById(user.ID)
			if userInfo == nil {
				return errors.New("Unknown User ID provided")
			}
			return config.EmailVerificationFeature.CreateAndSendCustomEmail(*userInfo, link)
		}
	}

	if config != nil && config.EmailVerificationFeature.GetEmailVerificationURL == nil {
		emailverificationTypeInput.GetEmailVerificationURL = nil
	} else {
		emailverificationTypeInput.GetEmailVerificationURL = func(user emailverificationModels.User) (string, error) {
			userInfo := recipeInstance.GetUserById(user.ID)
			if userInfo == nil {
				return "", errors.New("Unknown User ID provided")
			}
			return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo)
		}
	}

	return emailverificationTypeInput
}

func validateAndNormaliseResetPasswordUsingTokenConfig(appInfo supertokens.NormalisedAppinfo, signUpConfig models.TypeNormalisedInputSignUp, config *models.TypeInputResetPasswordUsingTokenFeature) models.TypeNormalisedInputResetPasswordUsingTokenFeature {
	normalisedInputResetPasswordUsingTokenFeature := models.TypeNormalisedInputResetPasswordUsingTokenFeature{
		FormFieldsForGenerateTokenForm: nil,
		FormFieldsForPasswordResetForm: nil,
		GetResetPasswordURL:            defaultGetResetPasswordURL(appInfo),
		CreateAndSendCustomEmail:       defaultCreateAndSendCustomPasswordResetEmail(appInfo),
	}

	if len(signUpConfig.FormFields) > 0 {
		var (
			formFieldsForPasswordResetForm []models.NormalisedFormField
			formFieldsForGenerateTokenForm []models.NormalisedFormField
		)
		for _, FormField := range signUpConfig.FormFields {
			if FormField.ID == constants.FormFieldPasswordID {
				formFieldsForPasswordResetForm = append(formFieldsForPasswordResetForm, FormField)
			}
			if FormField.ID == constants.FormFieldEmailID {
				formFieldsForGenerateTokenForm = append(formFieldsForGenerateTokenForm, FormField)
			}
		}
	}

	if config != nil && config.GetResetPasswordURL != nil {
		normalisedInputResetPasswordUsingTokenFeature.GetResetPasswordURL = config.GetResetPasswordURL
	}
	if config != nil && config.CreateAndSendCustomEmail != nil {
		normalisedInputResetPasswordUsingTokenFeature.CreateAndSendCustomEmail = config.CreateAndSendCustomEmail
	}

	return normalisedInputResetPasswordUsingTokenFeature
}

func defaultSetJwtPayloadForSession(_ models.User, _ []models.FormFieldValue, _ string) map[string]interface{} {
	return nil
}

func defaultSetSessionDataForSession(_ models.User, _ []models.FormFieldValue, _ string) map[string]interface{} {
	return nil
}

func validateAndNormaliseSessionFeatureConfig(config *models.TypeNormalisedInputSessionFeature) models.TypeNormalisedInputSessionFeature {
	normalisedInputSessionFeature := models.TypeNormalisedInputSessionFeature{
		SetJwtPayload:  defaultSetJwtPayloadForSession,
		SetSessionData: defaultSetSessionDataForSession,
	}

	if config != nil && config.SetJwtPayload != nil {
		normalisedInputSessionFeature.SetJwtPayload = config.SetJwtPayload
	}

	if config != nil && config.SetSessionData != nil {
		normalisedInputSessionFeature.SetSessionData = config.SetSessionData
	}

	return normalisedInputSessionFeature
}

func validateAndNormaliseSignInConfig(signUpConfig models.TypeNormalisedInputSignUp) models.TypeNormalisedInputSignIn {
	return models.TypeNormalisedInputSignIn{
		FormFields: normaliseSignInFormFields(signUpConfig.FormFields),
	}
}

func normaliseSignInFormFields(formFields []models.NormalisedFormField) []models.NormalisedFormField {
	normalisedFormFields := make([]models.NormalisedFormField, 0)
	if len(formFields) > 0 {
		for _, formField := range formFields {
			var (
				validate func(value interface{}) *string
				optional bool = false
			)
			if formField.ID == constants.FormFieldPasswordID {
				validate = formField.Validate
			} else if formField.ID == constants.FormFieldEmailID {
				validate = defaultEmailValidator
			}
			normalisedFormFields = append(normalisedFormFields, models.NormalisedFormField{
				ID:       formField.ID,
				Validate: validate,
				Optional: optional,
			})
		}
	}
	return normalisedFormFields
}

func validateAndNormaliseSignupConfig(config *models.TypeInputSignUp) models.TypeNormalisedInputSignUp {
	if config == nil {
		return models.TypeNormalisedInputSignUp{
			FormFields: normaliseSignUpFormFields(nil),
		}
	}
	return models.TypeNormalisedInputSignUp{
		FormFields: normaliseSignUpFormFields(config.FormFields),
	}
}

func normaliseSignUpFormFields(formFields []models.TypeInputFormField) []models.NormalisedFormField {
	var (
		normalisedFormFields     []models.NormalisedFormField
		formFieldPasswordIDCount = 0
		formFieldEmailIDCount    = 0
	)

	if len(formFields) > 0 {
		for _, formField := range formFields {
			var (
				validate func(value interface{}) *string
				optional bool = false
			)
			if formField.ID == constants.FormFieldPasswordID {
				formFieldPasswordIDCount++
				validate = defaultPasswordValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
			} else if formField.ID == constants.FormFieldEmailID {
				formFieldEmailIDCount++
				validate = defaultEmailValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
			} else {
				validate = defaultValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
				if formField.Optional != nil {
					optional = *formField.Optional
				}
			}
			normalisedFormFields = append(normalisedFormFields, models.NormalisedFormField{
				ID:       formField.ID,
				Validate: validate,
				Optional: optional,
			})
		}
	}
	if formFieldPasswordIDCount == 0 {
		normalisedFormFields = append(normalisedFormFields, models.NormalisedFormField{
			ID:       constants.FormFieldPasswordID,
			Validate: defaultPasswordValidator,
			Optional: false,
		})
	}
	if formFieldEmailIDCount == 0 {
		normalisedFormFields = append(normalisedFormFields, models.NormalisedFormField{
			ID:       constants.FormFieldEmailID,
			Validate: defaultEmailValidator,
			Optional: false,
		})
	}
	return normalisedFormFields
}

func defaultValidator(_ interface{}) *string {
	return nil
}

func defaultPasswordValidator(value interface{}) *string {
	// length >= 8 && < 100
    // must have a number and a character
	
	if reflect.TypeOf(value).Kind() != reflect.String {
		msg := "Development bug: Please make sure the password field yields a string"
		return &msg
	}
	if len(value.(string)) < 8 {
		msg := "Password must contain at least 8 characters, including a number"
		return &msg
	}
	if len(value.(string)) >= 100 {
		msg := "Password's length must be lesser than 100 characters"
		return &msg
	}
	alphaCheck, err := regexp.Match(`^.*[A-Za-z]+.*$`, []byte(value.(string)))
	if err != nil || !alphaCheck {
		msg := "Password must contain at least one alphabet"
		return &msg
	}
	numCheck, err := regexp.Match(`^.*[0-9]+.*$`, []byte(value.(string)))
	if err != nil || !numCheck {
		msg := "Password must contain at least one number"
		return &msg
	}
	return nil
}

func defaultEmailValidator(value interface{}) *string {
	if reflect.TypeOf(value).Kind() != reflect.String {
		msg := "Development bug: Please make sure the email field yields a string"
		return &msg
	}
	emailCheck, err := regexp.Match(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, []byte(value.(string)))
	if err != nil || !emailCheck {
		msg := "Email is invalid"
		return &msg
	}
	return nil
}
