package models

import "github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"

type TypeNormalisedInput struct {
	SignUpFeature                  TypeNormalisedInputSignUp
	SignInFeature                  TypeNormalisedInputSignIn
	ResetPasswordUsingTokenFeature TypeNormalisedInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       evmodels.TypeInput
	Override                       OverrideStruct
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evmodels.OverrideStruct
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
	Override                       *OverrideStruct
}

type TypeFormField struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
