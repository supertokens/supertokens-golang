package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type APIInterface struct {
	AuthorisationUrlGET func(provider TypeProvider, options APIOptions) (AuthorisationUrlGETResponse, error)
	SignInUpPOST        func(provider TypeProvider, code string, redirectURI string, options APIOptions) (SignInUpPOSTResponse, error)
}

type AuthorisationUrlGETResponse struct {
	OK *struct{ Url string }
}

type SignInUpPOSTResponse struct {
	OK *struct {
		CreatedNewUser   bool
		User             User
		AuthCodeResponse interface{}
	}
	NoEmailGivenByProviderError *struct{}
	FieldError                  *struct{ Error string }
}

type APIOptions struct {
	RecipeImplementation                  RecipeInterface
	EmailVerificationRecipeImplementation models.RecipeInterface
	Config                                TypeNormalisedInput
	RecipeID                              string
	Providers                             []TypeProvider
	Req                                   *http.Request
	Res                                   http.ResponseWriter
	OtherHandler                          http.HandlerFunc
}
