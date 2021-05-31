package emailverification

import (
	"net/http"
)

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
	id    string
	email string
}

type APIOptions struct {
	recipeImplementation RecipeInterface
	config               TypeNormalisedInput
	recipeId             string
	req                  http.Request
	res                  http.ResponseWriter
}

type APIInterface struct {
	verifyEmailPOST              func(token string, options APIOptions) map[string]interface{}
	isEmailVerifiedGET           func(options APIOptions) map[string]interface{}
	generateEmailVerifyTokenPOST func(options APIOptions) map[string]interface{}
}

type RecipeInterface struct {
	createEmailVerificationToken func(userId, email string) map[string]interface{}
	verifyEmailUsingToken        func(token string) map[string]interface{}
	isEmailVerified              func(userId, email string) bool
}
