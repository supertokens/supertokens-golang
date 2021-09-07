package models

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type UserInfo struct {
	ID    string
	Email *EmailStruct
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

type TypeProviderGetResponse struct {
	AccessTokenAPI        AccessTokenAPI
	AuthorisationRedirect AuthorisationRedirect
	GetProfileInfo        func(authCodeResponse interface{}) (UserInfo, error)
}

type AccessTokenAPI struct {
	URL    string
	Params map[string]string
}

type AuthorisationRedirect struct {
	URL    string
	Params map[string]interface{}
}

type TypeProvider struct {
	ID  string
	Get func(redirectURI *string, authCodeFromRequest *string) TypeProviderGetResponse
}

type User struct {
	ID         string
	TimeJoined uint64
	Email      string
	ThirdParty struct {
		ID     string
		UserID string
	}
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
}

type TypeInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeNormalisedInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeInput struct {
	SignInAndUpFeature       TypeInputSignInAndUp
	EmailVerificationFeature *TypeInputEmailVerificationFeature
	Override                 *struct {
		Functions                func(originalImplementation RecipeInterface) RecipeInterface
		APIs                     func(originalImplementation APIInterface) APIInterface
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeInterface) models.RecipeInterface
			APIs      func(originalImplementation models.APIInterface) models.APIInterface
		}
	}
}

type TypeNormalisedInput struct {
	SignInAndUpFeature       TypeNormalisedInputSignInAndUp
	EmailVerificationFeature models.TypeInput
	Override                 struct {
		Functions                func(originalImplementation RecipeInterface) RecipeInterface
		APIs                     func(originalImplementation APIInterface) APIInterface
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeInterface) models.RecipeInterface
			APIs      func(originalImplementation models.APIInterface) models.APIInterface
		}
	}
}
