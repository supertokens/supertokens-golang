package emailpassword

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance Recipe, appInfo supertokens.NormalisedAppinfo, config *models.TypeInput) models.TypeNormalisedInput {

	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	if config != nil && config.SignUpFeature != nil {
		typeNormalisedInput.SignUpFeature = validateAndNormaliseSignupConfig(config.SignUpFeature)
	}

	typeNormalisedInput.SignInFeature = validateAndNormaliseSignInConfig(typeNormalisedInput.SignUpFeature)

	if config != nil && config.ResetPasswordUsingTokenFeature != nil {
		typeNormalisedInput.ResetPasswordUsingTokenFeature = validateAndNormaliseResetPasswordUsingTokenConfig(appInfo, typeNormalisedInput.SignUpFeature, config.ResetPasswordUsingTokenFeature)
	}

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

	return typeNormalisedInput
}

func makeTypeNormalisedInput(recipeInstance Recipe) models.TypeNormalisedInput {
	signUpConfig := validateAndNormaliseSignupConfig(nil)
	return models.TypeNormalisedInput{
		SignUpFeature:                  signUpConfig,
		SignInFeature:                  validateAndNormaliseSignInConfig(signUpConfig),
		ResetPasswordUsingTokenFeature: validateAndNormaliseResetPasswordUsingTokenConfig(recipeInstance.RecipeModule.GetAppInfo(), signUpConfig, nil),
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
						return "", errors.New("Unknown User ID provided")
					}
					return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo)
				}
			}
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
			if FormField.ID == "password" {
				formFieldsForPasswordResetForm = append(formFieldsForPasswordResetForm, FormField)
			}
			if FormField.ID == "email" {
				formFieldsForGenerateTokenForm = append(formFieldsForGenerateTokenForm, FormField)
			}
		}
		normalisedInputResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm = formFieldsForGenerateTokenForm
		normalisedInputResetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm = formFieldsForPasswordResetForm
	}

	if config != nil && config.GetResetPasswordURL != nil {
		normalisedInputResetPasswordUsingTokenFeature.GetResetPasswordURL = config.GetResetPasswordURL
	}
	if config != nil && config.CreateAndSendCustomEmail != nil {
		normalisedInputResetPasswordUsingTokenFeature.CreateAndSendCustomEmail = config.CreateAndSendCustomEmail
	}

	return normalisedInputResetPasswordUsingTokenFeature
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
			if formField.ID == "password" {
				validate = formField.Validate
			} else if formField.ID == "email" {
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
			FormFields: NormaliseSignUpFormFields(nil),
		}
	}
	return models.TypeNormalisedInputSignUp{
		FormFields: NormaliseSignUpFormFields(config.FormFields),
	}
}

func NormaliseSignUpFormFields(formFields []models.TypeInputFormField) []models.NormalisedFormField {
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
			if formField.ID == "password" {
				formFieldPasswordIDCount++
				validate = defaultPasswordValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
			} else if formField.ID == "email" {
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
			ID:       "password",
			Validate: defaultPasswordValidator,
			Optional: false,
		})
	}
	if formFieldEmailIDCount == 0 {
		normalisedFormFields = append(normalisedFormFields, models.NormalisedFormField{
			ID:       "email",
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
