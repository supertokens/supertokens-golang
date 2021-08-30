package models

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type TypeNormalisedInput struct {
	SignUpFeature                  TypeNormalisedInputSignUp
	SignInFeature                  TypeNormalisedInputSignIn
	ResetPasswordUsingTokenFeature TypeNormalisedInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       models.TypeInput
	Override                       struct {
		Functions                func(originalImplementation RecipeInterface) RecipeInterface
		APIs                     func(originalImplementation APIInterface) APIInterface
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeInterface) models.RecipeInterface
			APIs      func(originalImplementation models.APIInterface) models.APIInterface
		}
	}
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
}

type TypeInputFormField struct {
	ID       string
	Validate func(value interface{}) *string
	Optional *bool
}

type TypeInputSignUp struct {
	FormFields []TypeInputFormField
}

type NormalisedFormField struct {
	ID       string
	Validate func(value interface{}) *string
	Optional bool
}

type TypeNormalisedInputSignUp struct {
	FormFields []NormalisedFormField
}

type TypeNormalisedInputSignIn struct {
	FormFields []NormalisedFormField
}

type TypeInputResetPasswordUsingTokenFeature struct {
	GetResetPasswordURL      func(user User) string
	CreateAndSendCustomEmail func(user User, passwordResetURLWithToken string)
}

type TypeNormalisedInputResetPasswordUsingTokenFeature struct {
	GetResetPasswordURL            func(user User) string
	CreateAndSendCustomEmail       func(user User, passwordResetURLWithToken string)
	FormFieldsForGenerateTokenForm []NormalisedFormField
	FormFieldsForPasswordResetForm []NormalisedFormField
}

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	TimeJoined uint64 `json:"timejoined"`
}

type TypeInput struct {
	SignUpFeature                  *TypeInputSignUp
	ResetPasswordUsingTokenFeature *TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *struct {
		Functions                func(originalImplementation RecipeInterface) RecipeInterface
		APIs                     func(originalImplementation APIInterface) APIInterface
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeInterface) models.RecipeInterface
			APIs      func(originalImplementation models.APIInterface) models.APIInterface
		}
	}
}

type TypeFormField struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
