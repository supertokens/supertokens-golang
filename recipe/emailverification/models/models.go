package models

import (
	"net/http"
)

type TypeInput struct {
	GetEmailForUserID        func(userID string) (string, error)
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
	Override                 *struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type TypeNormalisedInput struct {
	GetEmailForUserID        func(userID string) (string, error)
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
	Override                 struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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
	Status string `json:"status"`
	Token  string `json:"token"`
}

type VerifyEmailUsingTokenResponse struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}

type GenerateEmailVerifyTokenPOSTResponse struct {
	Status string `json:"status"`
}

type APIImplementation struct {
	VerifyEmailPOST              func(token string, options APIOptions) (*VerifyEmailUsingTokenResponse, error)
	IsEmailVerifiedGET           func(options APIOptions) (bool, error)
	GenerateEmailVerifyTokenPOST func(options APIOptions) (*GenerateEmailVerifyTokenPOSTResponse, error)
}

type RecipeImplementation struct {
	CreateEmailVerificationToken func(userID, email string) (*CreateEmailVerificationTokenResponse, error)
	VerifyEmailUsingToken        func(token string) (*VerifyEmailUsingTokenResponse, error)
	IsEmailVerified              func(userID, email string) (bool, error)
}
