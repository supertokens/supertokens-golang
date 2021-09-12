package evmodels

import "net/http"

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type APIInterface struct {
	VerifyEmailPOST              func(token string, options APIOptions) (VerifyEmailUsingTokenResponse, error)
	IsEmailVerifiedGET           func(options APIOptions) (IsEmailVerifiedGETResponse, error)
	GenerateEmailVerifyTokenPOST func(options APIOptions) (GenerateEmailVerifyTokenPOSTResponse, error)
}

type IsEmailVerifiedGETResponse struct {
	OK *struct {
		IsVerified bool
	}
}

type GenerateEmailVerifyTokenPOSTResponse struct {
	OK                        *struct{}
	EmailAlreadyVerifiedError *struct{}
}
