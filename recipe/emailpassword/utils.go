package emailpassword

import (
	"reflect"
	"regexp"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	emailverificationModels "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance models.RecipeImplementation, appInfo supertokens.NormalisedAppinfo, config models.TypeInput) models.TypeNormalisedInput {
	// sessionFeature := validateAndNormaliseSessionFeatureConfig(config.SessionFeature)
	// signUpFeature := validateAndNormaliseSignupConfig(config.SignUpFeature)
	// signInFeature := validateAndNormaliseSignInConfig(signUpFeature)
	// resetPasswordUsingTokenFeature := validateAndNormaliseResetPasswordUsingTokenConfig(appInfo, signUpFeature, config.ResetPasswordUsingTokenFeature)
	// emailVerificationFeature := validateAndNormaliseEmailVerificationConfig(recipeInstance, config)
	return models.TypeNormalisedInput{}
}
func validateAndNormaliseEmailVerificationConfig(recipeInstance models.RecipeImplementation, config models.TypeInput) emailverificationModels.TypeInput {
	return emailverificationModels.TypeInput{
		// GetEmailForUserID: recipeInstance.GetEmailForUserID,
	}
}

func validateAndNormaliseResetPasswordUsingTokenConfig(appInfo supertokens.NormalisedAppinfo, signUpConfig models.TypeNormalisedInputSignUp, config *models.TypeInputResetPasswordUsingTokenFeature) models.TypeNormalisedInputResetPasswordUsingTokenFeature {
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
	getResetPasswordURL := defaultGetResetPasswordURL(appInfo)
	if config != nil && config.GetResetPasswordURL != nil {
		getResetPasswordURL = config.GetResetPasswordURL
	}
	createAndSendCustomEmail := defaultCreateAndSendCustomPasswordResetEmail(appInfo)
	if config != nil && config.CreateAndSendCustomEmail != nil {
		createAndSendCustomEmail = config.CreateAndSendCustomEmail
	}
	return models.TypeNormalisedInputResetPasswordUsingTokenFeature{
		FormFieldsForGenerateTokenForm: formFieldsForGenerateTokenForm,
		FormFieldsForPasswordResetForm: formFieldsForPasswordResetForm,
		GetResetPasswordURL:            getResetPasswordURL,
		CreateAndSendCustomEmail:       createAndSendCustomEmail,
	}
}

func defaultSetJwtPayloadForSession(user models.User, formFields []models.TypeFormField, action string) map[string]interface{} {
	return map[string]interface{}{}
}

func defaultSetSessionDataForSession(user models.User, formFields []models.TypeFormField, action string) map[string]interface{} {
	return map[string]interface{}{}
}

func validateAndNormaliseSessionFeatureConfig(config *models.TypeNormalisedInputSessionFeature) models.TypeNormalisedInputSessionFeature {
	setJwtPayload := defaultSetJwtPayloadForSession
	if config != nil || config.SetJwtPayload != nil {
		setJwtPayload = config.SetJwtPayload
	}

	setSessionData := defaultSetSessionDataForSession
	if config != nil || config.SetSessionData != nil {
		setJwtPayload = config.SetSessionData
	}

	return models.TypeNormalisedInputSessionFeature{
		SetJwtPayload:  setJwtPayload,
		SetSessionData: setSessionData,
	}
}

func validateAndNormaliseSignInConfig(signUpConfig models.TypeNormalisedInputSignUp) models.TypeNormalisedInputSignIn {
	return models.TypeNormalisedInputSignIn{
		FormFields: normaliseSignInFormFields(signUpConfig.FormFields),
	}
}

func normaliseSignInFormFields(formFields []models.NormalisedFormField) []models.NormalisedFormField {
	var normalisedFormFields []models.NormalisedFormField
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
		return models.TypeNormalisedInputSignUp{}
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
	if err != nil || alphaCheck {
		msg := "Password must contain at least one alphabet"
		return &msg
	}
	numCheck, err := regexp.Match(`^.*[0-9]+.*$`, []byte(value.(string)))
	if err != nil || numCheck {
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
	if err != nil || emailCheck {
		msg := "Email is invalid"
		return &msg
	}
	return nil
}
