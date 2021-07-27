package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type TypeNormalisedInputSessionFeature struct {
	SetJwtPayload  func(user User, formFields []TypeFormField, action string) map[string]interface{}
	SetSessionData func(user User, formFields []TypeFormField, action string) map[string]interface{}
}

type TypeNormalisedInput struct {
	SessionFeature                 TypeNormalisedInputSessionFeature
	SignUpFeature                  TypeNormalisedInputSignUp
	SignInFeature                  TypeNormalisedInputSignIn
	ResetPasswordUsingTokenFeature TypeNormalisedInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       models.TypeInput
	Override                       struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
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
	CreateAndSendCustomEmail func(user User, passwordResetURLWithToken string) error
}

type TypeNormalisedInputResetPasswordUsingTokenFeature struct {
	GetResetPasswordURL            func(user User) string
	CreateAndSendCustomEmail       func(user User, passwordResetURLWithToken string) error
	FormFieldsForGenerateTokenForm []NormalisedFormField
	FormFieldsForPasswordResetForm []NormalisedFormField
}

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	TimeJoined uint64 `json:"timejoined"`
}

type TypeInput struct {
	SessionFeature                 *TypeNormalisedInputSessionFeature
	SignUpFeature                  *TypeInputSignUp
	ResetPasswordUsingTokenFeature *TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
		}
	}
}

type RecipeImplementation struct {
	SignUp                   func(email string, password string) SignInUpResponse
	SignIn                   func(email string, password string) SignInUpResponse
	GetUserByID              func(userID string) *User
	GetUserByEmail           func(email string) *User
	CreateResetPasswordToken func(userID string) CreateResetPasswordTokenResponse
	ResetPasswordUsingToken  func(token string, newPassword string) ResetPasswordUsingTokenResponse
}

type APIOptions struct {
	RecipeImplementation RecipeImplementation
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type APIImplementation struct {
	EmailExistsGET                 func(email string, options APIOptions) EmailExistsGETResponse
	GeneratePasswordResetTokenPOST func(formFields []TypeFormField, options APIOptions) GeneratePasswordResetTokenPOSTResponse
	PasswordResetPOST              func(formFields []TypeFormField, token string, options APIOptions) PasswordResetPOSTResponse
	SignInPOST                     func(formFields []TypeFormField, options APIOptions) SignInUpResponse
	SignUpPOST                     func(formFields []TypeFormField, options APIOptions) SignInUpResponse
}

type TypeFormField struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type EmailExistsGETResponse struct {
	Status string `json:"status"`
	Exist  bool   `json:"exist"`
}

type GeneratePasswordResetTokenPOSTResponse struct {
	Status string `json:"status"`
}

type PasswordResetPOSTResponse struct {
	Status string `json:"status"`
}
type SignInUpResponse struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}

type CreateResetPasswordTokenResponse struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}

type ResetPasswordUsingTokenResponse struct {
	Status string `json:"status"`
}
