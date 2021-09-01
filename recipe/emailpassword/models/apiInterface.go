package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type APIOptions struct {
	RecipeImplementation                  RecipeInterface
	EmailVerificationRecipeImplementation models.RecipeInterface
	Config                                TypeNormalisedInput
	RecipeID                              string
	Req                                   *http.Request
	Res                                   http.ResponseWriter
	OtherHandler                          http.HandlerFunc
}

type APIInterface struct {
	EmailExistsGET                 func(email string, options APIOptions) (EmailExistsGETResponse, error)
	GeneratePasswordResetTokenPOST func(formFields []TypeFormField, options APIOptions) (GeneratePasswordResetTokenPOSTResponse, error)
	PasswordResetPOST              func(formFields []TypeFormField, token string, options APIOptions) (ResetPasswordUsingTokenResponse, error)
	SignInPOST                     func(formFields []TypeFormField, options APIOptions) (SignInResponse, error)
	SignUpPOST                     func(formFields []TypeFormField, options APIOptions) (SignUpResponse, error)
}

type EmailExistsGETResponse struct {
	OK *struct{ Exists bool }
}

type GeneratePasswordResetTokenPOSTResponse struct {
	OK *struct{}
}
