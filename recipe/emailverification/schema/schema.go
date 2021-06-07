package schema

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/recipeimplementation"
)

type ReturnMap map[string]interface{}

var (
	CreateEmailVerificationTokenOk    ReturnMap = map[string]interface{}{"status": "OK", "token": ""}
	CreateEmailVerificationTokenError ReturnMap = map[string]interface{}{"status": "EMAIL_ALREADY_VERIFIED_ERROR"}
	VerifyEmailUsingTokenOk           ReturnMap = map[string]interface{}{"status": "OK", "user": User{}}
	VerifyEmailUsingTokenError        ReturnMap = map[string]interface{}{"status": "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR"}
)

// type RecipeImplementation struct {
// 	Querier supertokens.Querier
// }

// type APIImplementation struct{}

type TypeInput struct {
	GetEmailForUserID        func(userID string) string
	GetEmailVerificationURL  *func(userID User) string
	CreateAndSendCustomEmail *func(user User, emailVerificationURLWithToken string)
	Override                 Override
}

type Override struct {
	Functions func(originalImplementation recipeimplementation.RecipeImplementation) RecipeInterface
	APIs      func(originalImplementation api.APIImplementation) APIInterface
}

type TypeNormalisedInput struct {
	GetEmailForUserID        func(userID string) string
	GetEmailVerificationURL  func(userID User) string
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	Override                 Override
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
