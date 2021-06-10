package schema

import (
	"net/http"
)

type TypeInput struct {
	GetEmailForUserID        func(userID string) (string, error)
	GetEmailVerificationURL  func(user User) string
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	Override                 *struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type TypeNormalisedInput struct {
	GetEmailForUserID        func(userID string) (string, error)
	GetEmailVerificationURL  func(user User) string
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	Override                 struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type User struct {
	ID    string
	Email string
}

type APIOptions struct {
	RecipeImplementation RecipeImplementation
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type CreateEmailVerificationTokenResponse struct {
	OK *struct {
		Token string
	}
	EmailAlreadyVerifiedError bool // Zero value will be false
}

type CreateEmailVerificationTokenAPIResponse struct {
	OK                        bool
	EmailAlreadyVerifiedError bool // Zero value will be false
}

type VerifyEmailUsingTokenResponse struct {
	OK *struct {
		User User
	}
	InvalidTokenError bool // Zero value will be false
}

type APIImplementation struct {
	VerifyEmailPOST              func(token string, options APIOptions) (*VerifyEmailUsingTokenResponse, error)
	IsEmailVerifiedGET           func(options APIOptions) (bool, error)
	GenerateEmailVerifyTokenPOST func(options APIOptions) (*CreateEmailVerificationTokenAPIResponse, error)
}

type RecipeImplementation struct {
	CreateEmailVerificationToken func(userID, email string) (*CreateEmailVerificationTokenResponse, error)
	VerifyEmailUsingToken        func(token string) (*VerifyEmailUsingTokenResponse, error)
	IsEmailVerified              func(userID, email string) (bool, error)
}
