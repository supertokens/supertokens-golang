package schema

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type ReturnMap map[string]interface{}

var (
	CreateEmailVerificationTokenOk    ReturnMap = map[string]interface{}{"status": "OK", "token": ""}
	CreateEmailVerificationTokenError ReturnMap = map[string]interface{}{"status": "EMAIL_ALREADY_VERIFIED_ERROR"}
	VerifyEmailUsingTokenOk           ReturnMap = map[string]interface{}{"status": "OK", "user": User{}}
	VerifyEmailUsingTokenError        ReturnMap = map[string]interface{}{"status": "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR"}
)

type RecipeImplementation struct {
	Querier supertokens.Querier
}

type APIImplementation struct{}

type TypeInput struct {
	getEmailForUserId        func(userId string) string
	getEmailVerificationURL  func(userId string) string
	createAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	override                 struct {
		functions func(originalImplementation RecipeImplementation) RecipeInterface
		apis      func(originalImplementation APIImplementation) APIInterface
	}
}

type TypeNormalisedInput struct {
	getEmailForUserId        func(userId string) string
	getEmailVerificationURL  func(user User) string
	createAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	override                 struct {
		functions func(originalImplementation RecipeImplementation) RecipeInterface
		apis      func(originalImplementation APIImplementation) APIInterface
	}
}

type User struct {
	ID    string
	Email string
}

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
}

type APIInterface interface {
	VerifyEmailPOST(token string, options APIOptions) map[string]interface{}
	IsEmailVerifiedGET(options APIOptions) map[string]interface{}
	GenerateEmailVerifyTokenPOST(options APIOptions) map[string]interface{}
}

type RecipeInterface interface {
	CreateEmailVerificationToken(userId, email string) ReturnMap
	VerifyEmailUsingToken(token string) ReturnMap
	IsEmailVerified(userId, email string) bool
}
