package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type TypeFormField struct {
	ID    string
	Value interface{}
}

type TypeInputSetJwtPayloadForSession func(user User, formFields []TypeFormField, action string) map[string]interface{}

type TypeInputSetSessionDataForSession func(user User, formFields []TypeFormField, action string) map[string]interface{}

type TypeNormalisedInputSessionFeature struct {
	SetJwtPayload  TypeInputSetJwtPayloadForSession
	SetSessionData TypeInputSetSessionDataForSession
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
	GetEmailVerificationURL  func(user User) string
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
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
	Optional *bool
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
	TypeInputResetPasswordUsingTokenFeature
	FormFieldsForGenerateTokenForm []NormalisedFormField
	FormFieldsForPasswordResetForm []NormalisedFormField
}

type User struct {
	ID         string
	Email      string
	TimeJoined uint64
}

type TypeInput struct {
	SignUpFeature                  *TypeInputSignUp
	ResetPasswordUsingTokenFeature *TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
		}
	}
}

type RecipeImplementation struct {
	SignUp func(input struct {
		email    string
		password string
	}) SignInUpResponse
	SignIn func(input struct {
		email    string
		password string
	}) SignInUpResponse
	GetUserById              func(input struct{ userId string }) *User
	GetUserByEmail           func(input struct{ email string }) *User
	CreateResetPasswordToken func(input struct {
		userId string
	}) CreateResetPasswordTokenResponse
	ResetPasswordUsingToken func(input struct {
		token       string
		newPassword string
	}) ResetPasswordUsingTokenResponse
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
	EmailExistsGET                 func()
	GeneratePasswordResetTokenPOST func()
	PasswordResetPOST              func()
	SignInPOST                     func()
	SignUpPOST                     func()
}

type SignInUpResponse struct {
	Ok *struct {
		User User
	}
	Status string
}

type CreateResetPasswordTokenResponse struct {
	Ok *struct {
		Token string
	}
	Status string
}

type ResetPasswordUsingTokenResponse struct {
	Status string
}
